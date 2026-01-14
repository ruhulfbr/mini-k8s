package requests

type BuildConfigRequest struct {
	GitRepo           string `json:"git_repo" validate:"required,url"`
	GitBranch         string `json:"git_branch" validate:"required"`
	DockerContextPath string `json:"docker_context_path" validate:"required"`
	DockerfileName    string `json:"dockerfile_name" validate:"required"`
}

type CreateClusterRequest struct {
	Name       string              `json:"name" validate:"required"`
	IP         string              `json:"ip" validate:"required,ip"`
	Port       int                 `json:"port"`
	Replicas   int                 `json:"replicas" validate:"gte=1"`
	CPU        int                 `json:"cpu" validate:"required,min=1,max=4"`
	Memory     int                 `json:"memory" validate:"required,min=512,max=4096"`
	Path       string              `json:"path"`
	Type       string              `json:"type" validate:"required,oneof=http worker"`
	DeployMode int                 `json:"deploy_mode" validate:"required,oneof=1 2"`                          // 1 = image, 2 = build
	Image      *string             `json:"image" validate:"required_if=DeployMode 1,excluded_if=DeployMode 2"` // Required when DeployMode == 1 (image)
	Build      *BuildConfigRequest `json:"build,omitempty" validate:"required_if=DeployMode 2"`                // Required when DeployMode == 2 (build)
	Envs       map[string]string   `json:"envs,omitempty"`
}

type UpdateClusterRequest = CreateClusterRequest
