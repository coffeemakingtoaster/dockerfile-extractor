package ghapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

var BASE_URL = "https://api.github.com"

func GetDockerfilesFrom(repo string) int {
	tree := getFileTreeContent(repo, "main")
	getDockerfilesInTree(repo, tree)
	return 0
}

func getFileTreeContent(repo, branch string) GitTree {
	res, err := http.Get(fmt.Sprintf("%s/repos/%s/git/trees/%s?recursive=1", BASE_URL, repo, branch))
	if err != nil {
		panic(err)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var tree GitTree
	err = json.Unmarshal([]byte(resBody), &tree)
	if err != nil {
		panic(err)
	}
	return tree
}

func getDockerfilesInTree(repo string, tree GitTree) {
	for _, item := range tree.Tree {
		if strings.Contains(item.Path, "Dockerfile") {
			log.Debug().Msgf("Found: %s", item.Path)
			scrapeDockerfile(repo, item.Path)
		}
	}
}

func scrapeDockerfile(repo, dockerfilePath string) {
}
