package services

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
	"github.com/ruhulfbr/mini-k8s/internal/utils/fsUtils"
)

type GitService struct {
	dockerConfig config.DockerConfig
}

func NewGitService() *GitService {
	return &GitService{dockerConfig: config.GetDockerConfig()}
}

// ValidateRepoAndBranch checks if repo URL is valid and branch exists
func (gs *GitService) ValidateRepoAndBranch(repoURL, branch string) error {
	remote := git.NewRemote(memory.NewStorage(), &gitConfig.RemoteConfig{
		Name: "origin",
		URLs: []string{repoURL},
	})

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return appErrors.GitRepoNotFound
	}

	branchRef := plumbing.NewBranchReferenceName(branch)

	for _, ref := range refs {
		if ref.Name() == branchRef {
			return nil
		}
	}

	return appErrors.GitBranchNotExist
}

// CloneApplication clones repo into applications/<appName>
func (gs *GitService) CloneApplication(repoURL string, branch string, appName string) error {
	targetDir := fsUtils.Join(gs.dockerConfig.ApplicationPath, appName)

	// Ensure applications directory exists
	if err := fsUtils.EnsureDir(gs.dockerConfig.ApplicationPath); err != nil {
		logger.Error(nil, "Failed to create application directory for "+appName, err)
		return err
	}

	// Prevent overwriting existing app
	if fsUtils.DirExists(targetDir) {
		logger.Error(nil, "application directory already exist for: "+appName, nil)

		return fmt.Errorf("application '%s' already exists in system", appName)
	}

	_, err := git.PlainClone(targetDir, false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Depth:         1,
		Progress:      os.Stdout,
	})

	if err != nil {
		logger.Error(nil, "Failed to clone repository", err)

		return appErrors.GitFailedToClone
	}

	return nil
}

func (gs *GitService) PullApplication(appName, branch string) error {
	repoPath := fsUtils.Join(gs.dockerConfig.ApplicationPath, appName)

	if _, err := os.Stat(repoPath); err != nil {

		logger.Error(nil, "Repository not found to pull", err)
		return appErrors.GitRepoNotFound
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		logger.Error(nil, "Failed to open repo", err)

		return appErrors.GitRepoFailedToOpen
	}

	wt, err := repo.Worktree()
	if err != nil {
		logger.Error(nil, "Failed to open git worktree", err)
		return appErrors.GitFailedToOpenWorkTree
	}

	ref := plumbing.NewBranchReferenceName(branch)
	if err := wt.Checkout(&git.CheckoutOptions{
		Branch: ref,
	}); err != nil {

		logger.Error(nil, "Failed to checkout", err)

		return appErrors.GitFailedToCheckout
	}

	err = wt.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: ref,
		SingleBranch:  true,
	})

	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return nil
		}

		logger.Error(nil, "git pull failed", err)
		return appErrors.GitFailedToPull
	}

	return nil
}

func (gs *GitService) RemoveApplicationDir(appName string) error {
	return fsUtils.RemoveDir(fsUtils.Join(gs.dockerConfig.ApplicationPath, appName))
}
