package project

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/getantibody/antibody/naming"
)

type gitProject struct {
	URL     string
	Version string
	folder  string
}

// NewClonedGit is a git project that was already cloned, so, only Update
// will work here.
func NewClonedGit(home, folder string) Project {
	version := "master"
	version, _ = branch(folder)
	url := naming.FolderToURL(folder)
	return gitProject{
		folder:  filepath.Join(home, folder),
		Version: version,
		URL:     url,
	}
}

// NewGit A git project can be any repository in any given branch. It will
// be downloaded to the provided cwd
func NewGit(cwd, repo, version string) Project {
	var url string
	switch {
	case strings.HasPrefix(repo, "http://"):
		fallthrough
	case strings.HasPrefix(repo, "https://"):
		fallthrough
	case strings.HasPrefix(repo, "git://"):
		fallthrough
	case strings.HasPrefix(repo, "ssh://"):
		fallthrough
	case strings.HasPrefix(repo, "git@github.com:"):
		url = repo
	default:
		url = "https://github.com/" + repo
	}
	folder := filepath.Join(cwd, naming.URLToFolder(url))
	return gitProject{
		Version: version,
		URL:     url,
		folder:  folder,
	}
}

func (g gitProject) Download() error {
	if _, err := os.Stat(g.folder); os.IsNotExist(err) {
		cmd := exec.Command(
			"git", "clone", "--depth", "1", "-b", g.Version, g.URL, g.folder,
		)
		if bts, err := cmd.CombinedOutput(); err != nil {
			log.Println("git clone failed for", g.URL, string(bts))
			return err
		}
	}
	return nil
}

func (g gitProject) Update() error {
	if bts, err := exec.Command(
		"git", "-C", g.folder, "pull", "origin", g.Version,
	).CombinedOutput(); err != nil {
		log.Println("git update failed for", g.folder, string(bts))
		return err
	}
	return nil
}

func branch(folder string) (string, error) {
	branch, err := exec.Command(
		"git", "-C", folder, "rev-parse", "--abbrev-ref", "HEAD",
	).Output()
	return strings.Replace(string(branch), "\n", "", -1), err
}

func (g gitProject) Folder() string {
	return g.folder
}