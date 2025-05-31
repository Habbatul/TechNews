package service

import (
	"TechNews/config"
	"TechNews/data"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"google.golang.org/genai"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func resumeNews(text string, maxOutputToken int32) string {
	if len(text) < 10 {
		return ""
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GENAI_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Printf("Error creating genai client: %v", err)
		return ""
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemma-3n-e4b-it",
		genai.Text("Tanggapi hanya dengan satu paragraf bahasa indonesia tanpa ada tambahan baris. Jangan sertakan catatan. Pertahankan istilah teknis dalam Bahasa Inggris. Resume text berikut : "+text),
		&genai.GenerateContentConfig{
			MaxOutputTokens: maxOutputToken,
		},
	)
	if err != nil {
		log.Printf("Error generating content: %v", err)
		return ""
	}

	return result.Text()
}

// Todo: if ada yang sama on redis maka dont handle
var firstLinkNews string
var firstLinkVideo string

func getResumeData() data.ResumeResponse {
	result, err := config.RedisClient.Get(config.Ctx, "resume").Result()
	if err != nil {
		return data.ResumeResponse{Resume: data.Resumes{}}
	}
	var resumeResp data.ResumeResponse
	json.Unmarshal([]byte(result), &resumeResp)

	return data.ResumeResponse{
		Resume: resumeResp.Resume,
	}
}

func resumeFromTLDRTech() data.Resume {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://tldr.tech/api/rss/tech")
	if err != nil {
		log.Printf("Error parsing RSS feed: %v", err)
		return data.Resume{}
	}

	if len(feed.Items) == 0 {
		log.Println("No items found in RSS feed")
		return data.Resume{}
	}

	firstLinkNews = feed.Items[0].Link
	dataOld := getResumeData()
	if dataOld.Resume.Resume1.Result != "" && strings.Contains(dataOld.Resume.Resume1.Result, firstLinkNews) {
		fmt.Println("News Link : The link is already in old data, skip processing:")
		return dataOld.Resume.Resume1
	}

	//proses
	resp, err := http.Get(firstLinkNews)
	if err != nil {
		log.Printf("Failed to fetch article page: %v", err)
		return data.Resume{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Bad status code: %d", resp.StatusCode)
		return data.Resume{}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		return data.Resume{}
	}

	var result []string
	doc.Find(".newsletter-html").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i >= 2 {
			return false
		}
		text := strings.TrimSpace(s.Text())
		if text != "" {
			result = append(result, text)
		}
		return true
	})

	if len(result) == 0 {
		log.Println("Tidak menemukan elemen dengan class 'newsletter-html'")
		return data.Resume{}
	}

	output := ""
	if len(result) > 0 {
		output += "news-1:" + limitWord(result[0], 125)
	}
	if len(result) > 1 {
		output += ". news-2:" + limitWord(result[1], 125)
	}

	return data.Resume{Result: resumeNews(output, 225), Source: firstLinkNews}
}

func resumeFromFireshipVideo() data.Resume {
	//TODO : pastikan link youtube dapat baru eksekusi
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://www.youtube.com/feeds/videos.xml?channel_id=UCsBjURrPoezykLs9EqgamOA")
	if err != nil {
		log.Printf("Error parsing RSS feed: %v", err)
		return data.Resume{}
	}

	if len(feed.Items) == 0 {
		log.Println("No items found in RSS feed")
		return data.Resume{}
	}

	firstLinkVideo = feed.Items[0].Link
	//kalo link sebelumnya sama gausah dirangkum biar ga boros token
	dataOld := getResumeData()
	if dataOld.Resume.Resume2.Source != "" && strings.Contains(dataOld.Resume.Resume2.Source, firstLinkVideo) {
		fmt.Println("The link is already in old data, skip processing")
		return dataOld.Resume.Resume2
	}

	//proses
	reVideoID := regexp.MustCompile(`v=([a-zA-Z0-9_-]{11})`)
	matches := reVideoID.FindStringSubmatch(firstLinkVideo)
	if len(matches) < 2 {
		log.Println("Video ID not found in link")
		return data.Resume{}
	}
	videoID := matches[1]

	subtitle := downloadAndCleanSubtitle(videoID)
	return data.Resume{Result: resumeNews(limitWord(subtitle, 230), 230), Source: firstLinkVideo}
}

func downloadAndCleanSubtitle(videoID string) string {
	url := "https://www.youtube.com/watch?v=" + videoID

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching video page: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
	}
	body := string(bodyBytes)

	re := regexp.MustCompile(`ytInitialPlayerResponse\s*=\s*({.*});`)
	matches := re.FindStringSubmatch(body)
	if len(matches) < 2 {
		log.Println("ytInitialPlayerResponse JSON not found")
	}
	jsonStr := matches[1]

	var playerResponse data.InitialPlayerResponse
	err = json.Unmarshal([]byte(jsonStr), &playerResponse)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
	}

	if len(playerResponse.Captions.PlayerCaptionsTracklistRenderer.CaptionTracks) == 0 {
		log.Println("No captions available for this video")
	}

	subtitleUrl := ""
	for _, track := range playerResponse.Captions.PlayerCaptionsTracklistRenderer.CaptionTracks {
		if track.Lang == "en" {
			subtitleUrl = track.BaseUrl + "&fmt=vtt"
			break
		}
	}

	if subtitleUrl == "" {
		log.Println("No English subtitle found")
	}

	log.Println(subtitleUrl)

	subResp, err := http.Get(subtitleUrl)
	if err != nil {
		log.Printf("Error fetching subtitle: %v", err)
	}
	defer subResp.Body.Close()

	if subResp.StatusCode != 200 {
		log.Printf("Failed to get subtitle, status: %s", subResp.Status)
	}

	subtitleData, err := io.ReadAll(subResp.Body)
	if err != nil {
		log.Printf("Error reading subtitle data: %v", err)
	}

	//bersihkan subtitle WebVTT: header, timestamps, empty line, dll
	text := cleanVtt(string(subtitleData))
	return text
}

func cleanVtt(vtt string) string {
	lines := strings.Split(vtt, "\n")
	var result []string
	seen := make(map[string]bool)

	//regex untuk deteksi baris timestamp
	timestampRe := regexp.MustCompile(`^\d{2}:\d{2}:\d{2}\.\d{3}\s+-->\s+\d{2}:\d{2}:\d{2}\.\d{3}`)

	//regex untuk menangkap isi dalam tag: <...>isi</...> atau <...>isi
	tagContentRe := regexp.MustCompile(`<[^>]+>([^<]+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || timestampRe.MatchString(line) {
			continue
		}

		//ambil semua isi tag, gabungkan jadi satu baris
		matches := tagContentRe.FindAllStringSubmatch(line, -1)
		var parts []string
		for _, m := range matches {
			word := strings.TrimSpace(m[1])
			if word != "" {
				parts = append(parts, word)
			}
		}
		cleaned := strings.Join(parts, " ")
		if cleaned == "" {
			continue
		}

		//klo belum pernah muncul, masukkan ke hasil
		if !seen[cleaned] {
			seen[cleaned] = true
			result = append(result, cleaned)
		}
	}

	return strings.Join(result, " ")
}

func limitWord(text string, lenText int) string {
	word := strings.Fields(text)
	if len(word) > lenText {
		word = word[:lenText]
	}
	return strings.Join(word, " ")
}

func GetResume() data.ResumeResponse {
	resumeText := resumeFromTLDRTech()
	resumeVideo := resumeFromFireshipVideo()
	return data.ResumeResponse{
		Resume: data.Resumes{
			Resume1: resumeText,
			Resume2: resumeVideo,
		},
	}
}
