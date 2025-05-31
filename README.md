# 📰 TechNews API

**TechNews** is a RESTful API that collects technology news from multiple sources by aggregating RSS feeds. It also summarizes content from Fireship YouTube videos and TLDR.tech articles using Google’s Gemini LLM SDK.

🔗 **[Open Frontend Web (Not Done Yet)](belum)**

---

## ✨ Features

* **Aggregates Tech News from Multiple RSS Feeds**
  Collects and unifies tech-related news from various RSS sources using [`gofeed`](https://github.com/mmcdole/gofeed). Some sites, like Reddit, require proxy access due to regional ISP restrictions (Indonesian ISPs).
* **Proxy Support for Blocked Sites**
  Uses a proxy server to access sources like Reddit that may be blocked by Indonesian ISPs.
* **Summarizes Content Using the Google Gemini LLM SDK**

    * Summarizes YouTube videos from **Fireship** (latest video)
    * Summarizes articles from **TLDR.tech** (2 latest news)
    * All summarization is done after pre-processing and applying text limits, in order to optimize token usage.

* **Scraping and Text Preprocessing**
  Content scraping is done manually (not by the LLM), with the help of [`gofeed`](https://github.com/mmcdole/gofeed) & [`goquery`](https://github.com/PuerkitoBio/goquery). Text is processed before being sent to the LLM to ensure efficiency.
---

## 🚀 How It Works

1. Fetches data from various RSS sources.
2. Uses a proxy if required to bypass geo-restrictions (e.g., Reddit access).
3. Scrapes and preprocesses the text content.
4. Sends cleaned and truncated text to the Gemini LLM for summarization.
5.  API endpoints are available to update and store summaries and news lists directly into Redis.
6. You can trigger the POST /api/resume and POST /api/feed endpoints periodically using a cron job, a scheduler, or any approach you prefer.
7. Returns the summarized content and news lists via RESTful endpoints (retrieved from Redis storage).
---

## ⚙️ Environment Variables

| Variable         | Description                                                                                               | Required? |
| ---------------- |-----------------------------------------------------------------------------------------------------------|-----------|
| `PROXY_URL`      | Proxy server URL and port (e.g., `http://your-proxy.com:8080`) used to access restricted sources. | DEPEND    |
| `REDIS_ADDR`     | Redis server address used for storage.                                                                    | YES       |
| `REDIS_PASSWORD` | Redis password (if authentication is required).                                                           | DEPEND    |
| `GENAI_KEY` | Your Google AI API Key.                                                                                   | YES       |

---

## 🛠️ Tech Stack

* **RSS Parsing:** [`gofeed`](https://github.com/mmcdole/gofeed)
* **Scraping:** [`goquery`](https://github.com/PuerkitoBio/goquery)
* **LLM Integration:** Google Gemini API SDK
* **Data Store:** Redis

---
## ⚠️ NOTE
There is a possibility of failure when attempting to resume YouTube videos due to frequent changes in how YouTube structures and delivers its public-facing data. These changes often require adjustments to the scraping logic to maintain compatibility

---
## 📝 TODO
Future updates may be as follows :
1. Implement a security mechanism for the `POST` endpoints to protect against unauthorized access.
2. Develop a method to generate more consistent and accurate summaries using the LLM.
3. Fix various bugs and handle LLM API failures gracefully to prevent unnecessary token.
4. Handle potential errors more properly.
---

## 📝 Changelog

- Updated: YouTube's `ytInitialPlayerResponse` structure changed; adjusted RegExp to use greedy quantifier `.*` instead of `.*?` to fully capture JSON.
- Updated: `Reddit .rss` to `feedSources`.
- Added: `Proxy.go` to proxy Reddit RSS feeds via `gofeed` (for Indonesian ISP blocks).
- Fixed: Prevent sending requests to LLM when text input is blank.
