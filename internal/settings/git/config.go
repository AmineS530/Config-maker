package git

type Config struct {
	ConfigureGit bool   `json:"configure_git"`
	GitName      string `json:"git_name"`
	GitEmail     string `json:"git_email"`
}
