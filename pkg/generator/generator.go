package generator

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Generator struct {
	gitlabProjectBaseURL string
	gitlabRepo           string
	gitlabBranch         string
	sys                  sysWrapper
}

func New(gitlabBaseProjectURL, gitlabRepoName, gitlabBranch string) *Generator {
	sys := &sysWrap{}
	return &Generator{gitlabBaseProjectURL,
		getGitlabRepo(gitlabRepoName, sys),
		gitlabBranch,
		sys}
}

func (g *Generator) Generate(vars map[string]interface{}) (string, error) {
	uri, _ := url.Parse(fmt.Sprintf("%s/%s/pipelines/new", g.gitlabProjectBaseURL, g.gitlabRepo))
	query := url.Values{}
	if g.gitlabBranch != "" {
		query.Set("ref", g.gitlabBranch)
	}
	for k, v := range vars {
		key := parseVarsKey(k)
		val, err := g.parseVarsValue(v.(string))
		if err != nil {
			return "", err
		}
		query.Set(key, val)
	}
	uri.RawQuery = query.Encode()
	return uri.String(), nil
}

func getGitlabRepo(repo string, sys sysWrapper) string {
	return strings.ReplaceAll(repo, "${GIT_REPO}", getGitRepoDir(sys))
}

func getGitRepoDir(sys sysWrapper) string {
	getwd, _ := sys.osGetwd()
	splits := strings.Split(getwd, string(os.PathSeparator))
	return splits[len(splits)-1]
}

func parseVarsKey(key string) string {
	return fmt.Sprintf("var[%s]", key)
}

func (g *Generator) parseVarsValue(val string) (string, error) {
	head, err := g.getGitHead()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(val, "${GIT_HEAD}", head), nil
}

func (g *Generator) getGitHead() (string, error) {
	gitDir, err := g.sys.execCommand("git", "rev-parse --git-dir")
	if err != nil {
		gitDir = []byte(".git")
	}
	f, err := g.sys.osReadFile(string(gitDir) + "/HEAD")
	if err != nil {
		return "", err
	}
	head := strings.Split(string(f), "ref: refs/heads/")[1]
	return strings.TrimSpace(head), nil
}
