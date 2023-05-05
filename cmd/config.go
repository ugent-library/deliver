package cmd

type Config struct {
	Production bool
	Host       string
	Port       int
	Admins     []string
	Storage    struct {
		Backend string
		Conn    string
	}
	Repo struct {
		Conn string
	}
	OIDC struct {
		URL         string
		ID          string
		Secret      string
		RedirectURL string `mapstructure:"redirect_url"`
	}
	Cookies struct {
		Secret string
	}
	MaxFileSize int64 `mapstructure:"max_file_size"`
}
