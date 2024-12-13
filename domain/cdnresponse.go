package domain

type CdnResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		FileId      string `json:"fileId"`
		Name        string `json:"name"`
		Size        int    `json:"size"`
		VersionInfo struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"versionInfo"`
		FilePath     string      `json:"filePath"`
		Url          string      `json:"url"`
		FileType     string      `json:"fileType"`
		Height       int         `json:"height"`
		Width        int         `json:"width"`
		ThumbnailUrl string      `json:"thumbnailUrl"`
		AITags       interface{} `json:"AITags"`
	} `json:"data"`
}
