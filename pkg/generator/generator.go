package generator

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

type Generator struct {
	gitlabProjectBaseURL string
	gitlabRepo           string
	gitlabBranch         string
}

func New(gitlabBaseProjectURL, gitlabRepoName, gitlabBranch string) *Generator {
	return &Generator{gitlabBaseProjectURL,
		getGitlabRepo(gitlabRepoName),
		gitlabBranch}
}

func (g *Generator) Generate(vars map[string]interface{}) (string, error) {
	uri, _ := url.Parse(fmt.Sprintf("%s/%s/pipelines/new", g.gitlabProjectBaseURL, g.gitlabRepo))
	query := url.Values{}
	if g.gitlabBranch != "" {
		query.Set("ref", g.gitlabBranch)
	}
	for k, v := range vars {
		key := parseVarsKey(k)
		val, err := parseVarsValue(v.(string))
		if err != nil {
			return "", err
		}
		query.Set(key, val)
	}
	uri.RawQuery = query.Encode()
	return uri.String(), nil
}

func getGitlabRepo(repo string) string {
	return strings.ReplaceAll(repo, "${GIT_REPO}", getGitRepoDir())
}

func getGitRepoDir() string {
	getwd, _ := os.Getwd()
	splits := strings.Split(getwd, string(os.PathSeparator))
	return splits[len(splits)-1]
}

func parseVarsKey(key string) string {
	return fmt.Sprintf("var[%s]", key)
}

func parseVarsValue(val string) (string, error) {
	head, err := getGitHead()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(val, "${GIT_HEAD}", head), nil
}

func getGitHead() (string, error) {
	gitDir, err := exec.Command("git", "rev-parse --git-dir").Output()
	if err != nil {
		gitDir = []byte(".git")
	}
	f, err := os.ReadFile(string(gitDir) + "/HEAD")
	if err != nil {
		return "", err
	}
	head := strings.Split(string(f), "ref: refs/heads/")[1]
	return strings.TrimSpace(head), nil
}
