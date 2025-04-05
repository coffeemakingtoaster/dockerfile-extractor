package ghapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/coffeemakingtoaster/dockerfile-extractor/pkg/util"
	"github.com/rs/zerolog/log"
)

var BASE_URL = "https://api.github.com"

type DockerFileInformation struct {
	Repo    string
	Path    string
	Content string
}

func CanUseAuth() bool {
	return os.Getenv("GH_TOKEN") != ""
}

func doRequest(url string) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if CanUseAuth() {
		req.Header = http.Header{
			"Authorization": {fmt.Sprintf("Bearer %s", os.Getenv("GH_TOKEN"))},
		}
	}

	res, err := client.Do(req)

	if res.StatusCode != 200 {
		log.Error().Msg(url)
		//bdy, _ := io.ReadAll(res.Body)
		//panic(fmt.Sprintf("Did not get expected status code %d: %s", res.StatusCode, string(bdy)))
		return nil, err
	}

	// rate limiting?
	time.Sleep(1)

	return res, err
}

func (dfi *DockerFileInformation) PopulateContent() error {
	url, err := getFileURL(*dfi)
	if err != nil {
		return err
	}
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	dfi.Content = string(resBody)
	return nil
}

func (dfi DockerFileInformation) HashContent() (string, error) {
	if dfi.Content == "" {
		return "", errors.New("No content present. File has not been fetched yet")
	}
	return util.HashString(dfi.Content), nil
}

func (dfi DockerFileInformation) SaveToFile(targetPath string) error {
	hash, err := dfi.HashContent()
	if err != nil {
		return err
	}
	dest := path.Join(targetPath, fmt.Sprintf("%s.Dockerfile", hash))
	signature := fmt.Sprintf("# %s - %s", dfi.Repo, dfi.Path)
	content := []byte(fmt.Sprintf("%s\n%s", signature, dfi.Content))
	util.CreateDirIfNotExist(targetPath)
	err = os.WriteFile(dest, content, 0644)
	return err
}

func GetDockerfilesFrom(repo string) []DockerFileInformation {
	defaultBranch := getRepositoryDefaultBranch(repo)
	tree := getFileTreeContent(repo, defaultBranch)
	return getDockerfilesInTree(repo, tree)
}

func getRepositoryDefaultBranch(repo string) string {
	res, err := doRequest(fmt.Sprintf("%s/repos/%s", BASE_URL, repo))
	if err != nil {
		panic(err)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var info RepoInfo
	err = json.Unmarshal([]byte(resBody), &info)
	if err != nil {
		panic(err)
	}
	return info.DefaultBranch

}

func getFileTreeContent(repo, branch string) GitTree {
	res, err := doRequest(fmt.Sprintf("%s/repos/%s/git/trees/%s?recursive=1", BASE_URL, repo, branch))
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

func getDockerfilesInTree(repo string, tree GitTree) []DockerFileInformation {
	res := []DockerFileInformation{}
	for _, item := range tree.Tree {
		if strings.Contains(item.Path, "Dockerfile") {
			res = append(res,
				DockerFileInformation{
					Repo: repo,
					Path: item.Path,
				})
		}
	}
	return res
}

func getFileURL(info DockerFileInformation) (string, error) {
	res, err := doRequest(fmt.Sprintf("%s/repos/%s/contents/%s", BASE_URL, info.Repo, info.Path))
	if err != nil {
		return "", err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var tree ContentTree
	err = json.Unmarshal([]byte(resBody), &tree)
	if err != nil {
		return "", err
	}
	if tree.DownloadURL == nil {
		return "", errors.New("Download url was nil Pointer")
	}
	return *tree.DownloadURL, nil
}
