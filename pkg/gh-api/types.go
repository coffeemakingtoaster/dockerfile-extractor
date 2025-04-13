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

type ContentTree struct {
	Type        string  `json:"type"`
	Size        int     `json:"size"`
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	SHA         string  `json:"sha"`
	Content     string  `json:"content,omitempty"`
	URL         string  `json:"url"`
	GitURL      *string `json:"git_url"`
	HTMLURL     *string `json:"html_url"`
	DownloadURL *string `json:"download_url"`
	Entries     []Entry `json:"entries,omitempty"`
	Links       Links   `json:"_links"`
}

type Entry struct {
	Type        string  `json:"type"`
	Size        int     `json:"size"`
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	SHA         string  `json:"sha"`
	URL         string  `json:"url"`
	GitURL      *string `json:"git_url"`
	HTMLURL     *string `json:"html_url"`
	DownloadURL *string `json:"download_url"`
	Links       Links   `json:"_links"`
}

type Links struct {
	Git  *string `json:"git"`
	HTML *string `json:"html"`
	Self string  `json:"self"`
}

type RepoInfo struct {
	DefaultBranch string `json:"default_branch"`
}

type ContributerInfo struct {
	Id   int    `json:"id"`
	Name string `json:"login"`
}

type RepositoryOverviewInfo struct {
	Name string `json:"full_name"`
}
