package main

// Captcha stores the POST data used to verify a TikTok captcha.
type Captcha struct {
	ModifiedImgWidth int    `json:"modified_img_width"`
	ID               string `json:"id"`
	Mode             string `json:"mode"`
	Reply            []struct {
		X            int `json:"x"`
		Y            int `json:"y"`
		RelativeTime int `json:"relative_time"`
	} `json:"reply"`
	Models struct {
		X struct {
			Time int64   `json:"time"`
			X    float64 `json:"x"`
			Y    int     `json:"y"`
		} `json:"x"`
		Y struct {
			X    float64 `json:"x"`
			Y    int     `json:"y"`
			Time int64   `json:"time"`
		} `json:"y"`
		Z []struct {
			X    float64 `json:"x"`
			Y    int     `json:"y"`
			Time int64   `json:"time"`
		} `json:"z"`
		T []interface{} `json:"t"`
		M []struct {
			X    float64 `json:"x"`
			Y    int     `json:"y"`
			Time int64   `json:"time"`
		} `json:"m"`
	} `json:"models"`
	LogParams string `json:"log_params"`
	Reply2    []struct {
		X            int `json:"x"`
		Y            int `json:"y"`
		RelativeTime int `json:"relative_time"`
	} `json:"reply2"`
	Models2 struct {
		X struct {
			Time int64   `json:"time"`
			X    float64 `json:"x"`
			Y    int     `json:"y"`
		} `json:"x"`
		Y struct {
			X    float64 `json:"x"`
			Y    int     `json:"y"`
			Time int64   `json:"time"`
		} `json:"y"`
		Z []struct {
			X    float64 `json:"x"`
			Y    int     `json:"y"`
			Time int64   `json:"time"`
		} `json:"z"`
		T []interface{} `json:"t"`
		M []struct {
			X    float64 `json:"x"`
			Y    int     `json:"y"`
			Time int64   `json:"time"`
		} `json:"m"`
	} `json:"models2"`
}
