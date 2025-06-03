package utils

// PageData holds common metadata for pages.
type PageData struct {
	Lang               string
	Theme              string
	CurrentURL         string // The full URL of the current page, including query parameters
	Icon               string
	Title              string
	Description        string
	Keywords           string
	OgType             string
	OgImage            string
	OgImageWidth       string
	OgImageHeight      string
	UrlPath            string // The relative path of the page (e.g., "/sketches")
	CanonicalUrl       string
	SupportedLanguages []LanguageInfo
	CurrentLanguage    string
}

// GetDefaultPageData initializes PageData with default values.
// urlPath should be the relative path for the page (e.g., "/", "/sketches").
// lang and theme are the current UI language and theme.
// requestURI is r.RequestURI which includes path and query.
func GetDefaultPageData(urlPath, currentLang, theme, requestURI string) *PageData {
	return &PageData{
		Lang:               currentLang,
		Theme:              theme,
		CurrentURL:         GetFullURL(requestURI),
		Icon:               "🪄",
		Title:              "Creative Coding Bookclub",
		Description:        "Creative coding sketches and projects by the bookclub members.",
		Keywords:           "Creative Coding, Bookclub, Sketches, Art, Code, Generative Art, p5js, threejs",
		OgType:             "article",
		OgImage:            GetFullURL("/assets/ogimage.jpeg"),
		OgImageWidth:       "1200",
		OgImageHeight:      "630",
		UrlPath:            urlPath,
		CanonicalUrl:       GetFullURL(urlPath),
		SupportedLanguages: GetSupportedLanguages(),
		CurrentLanguage:    currentLang,
	}
}
