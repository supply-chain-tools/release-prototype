package count

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/supply-chain-tools/go-sandbox/gitkit"
	"os"
)

func Commits(state *gitkit.RepoState) int {
	return len(state.CommitMap)
}

func LoadRepoStateFromCurrentDirectory() (*gitkit.RepoState, error) {
	repo, err := openRepoFromCurrentDirectory()
	if err != nil {
		return nil, err
	}

	return gitkit.LoadRepoState(repo), nil
}

func openRepoFromCurrentDirectory() (*git.Repository, error) {
	basePath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repoDir, found, err := gitkit.GetRootPathOfLocalGitRepo(basePath)
	if err != nil {
		return nil, fmt.Errorf("unable infer git root from %s: %w", basePath, err)
	}

	if !found {
		return nil, fmt.Errorf("not in a git repo %s", basePath)
	}

	repo, err := gitkit.OpenRepoInLocalPath(repoDir)
	if err != nil {
		return nil, fmt.Errorf("unable to open repo %s: %w", repoDir, err)
	}

	return repo, nil
}
