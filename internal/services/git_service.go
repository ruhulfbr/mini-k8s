package services

import (
	"context"
	"errors"
	"os"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
	"github.com/ruhulfbr/mini-k8s/internal/utils/fsUtils"
)

type GitService struct {
	dockerConfig config.DockerConfig
}

func NewGitService() *GitService {
	return &GitService{dockerConfig: config.GetDockerConfig()}
}

// ValidateRepoAndBranch checks if clusterRepo URL is valid and branch exists
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

// CloneApplication clones clusterRepo into clusters/<appName>
func (gs *GitService) CloneApplication(appName string, clusterName string, repoURL string, branch string) error {
	repoPath := gs.repoPath(appName, clusterName)
	ctx := context.Background()

	// Ensure clusters directory exists
	if err := fsUtils.EnsureDir(gs.dockerConfig.ClusterPath); err != nil {
		logger.Error(ctx, "Failed to create clusters base directory", err,
			"application", appName,
			"cluster", clusterName,
			"repository", repoPath,
		)
		return err
	}

	// Prevent overwriting existing app
	if fsUtils.DirExists(repoPath) {
		logger.Error(ctx, "Cluster repository already cloned", nil,
			"application", appName,
			"cluster", clusterName,
			"repository", repoPath,
		)
		return appErrors.GitRepoAlreadyCloned
	}

	_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Depth:         1,
		Progress:      os.Stdout,
	})

	if err != nil {
		logger.Error(ctx, "Failed to clone repository", err,
			"application", appName,
			"cluster", clusterName,
			"repository", repoPath,
		)
		return appErrors.GitFailedToClone
	}

	return nil
}

func (gs *GitService) PullApplication(appName string, clusterName string, buildConfig *entities.ClusterBuildConfig) error {
	repoPath := gs.repoPath(appName, clusterName)

	if !fsUtils.DirExists(repoPath) {
		return gs.CloneApplication(appName, clusterName, buildConfig.GitRepo, buildConfig.GitBranch)
	}

	ctx := context.Background()

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		logger.Error(ctx, "Failed to open clusterRepo", err,
			"application", appName,
			"cluster", clusterName,
			"repository", repoPath,
		)
		return appErrors.GitRepoFailedToOpen
	}

	wt, err := repo.Worktree()
	if err != nil {
		logger.Error(nil, "Failed to open git worktree", err,
			"application", appName,
			"cluster", clusterName,
			"repository", repoPath,
		)
		return appErrors.GitFailedToOpenWorkTree
	}

	ref := plumbing.NewBranchReferenceName(buildConfig.GitBranch)
	if err := wt.Checkout(&git.CheckoutOptions{
		Branch: ref,
	}); err != nil {
		logger.Error(nil, "Failed to checkout", err,
			"application", appName,
			"cluster", clusterName,
			"repository", repoPath,
		)
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

		logger.Error(nil, "git pull failed", err,
			"application", appName,
			"cluster", clusterName,
			"repository", repoPath,
		)
		return appErrors.GitFailedToPull
	}

	return nil
}

func (gs *GitService) repoPath(appName string, clusterName string) string {
	return fsUtils.Join(gs.dockerConfig.ClusterPath, appName+"-"+clusterName)
}

func (gs *GitService) RemoveApplicationDir(appName string) error {
	return fsUtils.RemoveDir(fsUtils.Join(gs.dockerConfig.ClusterPath, appName))
}
