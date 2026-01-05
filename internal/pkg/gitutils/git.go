package gitutils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
)

var baseDir string = "applications"

// ValidateRepoAndBranch checks if repo URL is valid and branch exists
func ValidateRepoAndBranch(repoURL, branch string) error {
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
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
func CloneApplication(repoURL, branch, appName string) error {
	targetDir := filepath.Join(baseDir, appName)

	// Ensure applications directory exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {

		logger.Error(nil, "Failed to create application directory for "+appName, err)
		return err
	}

	// Prevent overwriting existing app
	if _, err := os.Stat(targetDir); err == nil {
		logger.Error(nil, "application directory already exist for: "+appName, err)

		return fmt.Errorf("application '%s' already exists", appName)
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

func PullApplication(appName, branch string) error {
	repoPath := filepath.Join(baseDir, appName)

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

	// 4. Checkout branch
	ref := plumbing.NewBranchReferenceName(branch)
	if err := wt.Checkout(&git.CheckoutOptions{
		Branch: ref,
	}); err != nil {

		logger.Error(nil, "Failed to checkout", err)

		return appErrors.GitFailedToCheckout
	}

	// 5. Pull latest changes
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

func RemoveApplicationDir(appName string) error {
	return os.RemoveAll(filepath.Join(baseDir, appName))
}
