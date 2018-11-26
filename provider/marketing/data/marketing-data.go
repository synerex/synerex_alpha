package mkdata

type Request struct {
	Command  string    `json:"command"`
	Contents []Content `json:"contents"`
}

type Content struct {
	Type   string `json:"type"`
	Data   string `json:"data"`
	Period int    `json:"period"`
}

type Result struct {
	Command string `json:"command"`
	Results []struct {
		Data string `json:"data"`
	}
}
