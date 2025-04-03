package ghapi

type GitTree struct {
	SHA       string     `json:"sha"`
	URL       string     `json:"url"`
	Truncated bool       `json:"truncated"`
	Tree      []TreeItem `json:"tree"`
}

type TreeItem struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	SHA  string `json:"sha"`
	Size *int   `json:"size,omitempty"`
	URL  string `json:"url"`
}
