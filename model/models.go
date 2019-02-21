package model

type Routes struct {
	Routes []Route `json:"routes"`
}

type Route struct {
	Path     string   `json:"path"`
	Method   string   `json:"method"`
	Contents Contents `json:"contents"`
}

type Contents struct {
	Default  Content   `json:"default"`
	Handlers []Handler `json:"handlers"`
}

type Content struct {
	Body   string `json:"body"`
	Status int    `json:"status"`
}

type Handler struct {
	Content Content                `json:"content"`
	Param   map[string]interface{} `json:"param"`
	Header  map[string]interface{} `json:"header"`
}
