package data

type CaptionTrack struct {
	BaseUrl string `json:"baseUrl"`
	Lang    string `json:"languageCode"`
}

type PlayerCaptionsTracklistRenderer struct {
	CaptionTracks []CaptionTrack `json:"captionTracks"`
}

type Captions struct {
	PlayerCaptionsTracklistRenderer PlayerCaptionsTracklistRenderer `json:"playerCaptionsTracklistRenderer"`
}

type InitialPlayerResponse struct {
	Captions Captions `json:"captions"`
}
