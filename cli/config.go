package cli

import "fmt"

type Config struct {
	// Version info
	Version struct {
		Branch string `env:"SOURCE_BRANCH"`
		Commit string `env:"SOURCE_COMMIT"`
		Image  string `env:"IMAGE_NAME"`
	}
	// Env must be local, development, test or production
	Env     string   `env:"DELIVER_ENV" envDefault:"production"`
	Host    string   `env:"DELIVER_HOST"`
	Port    int      `env:"DELIVER_PORT" envDefault:"3000"`
	Admins  []string `env:"DELIVER_ADMINS,notEmpty"`
	Storage struct {
		Backend string `env:"BACKEND" envDefault:"s3"`
		Conn    string `env:"CONN,notEmpty"`
	} `envPrefix:"DELIVER_STORAGE_"`
	Repo struct {
		Conn string `env:"CONN,notEmpty"`
	} `envPrefix:"DELIVER_REPO_"`
	OIDC struct {
		URL         string `env:"URL,notEmpty"`
		ID          string `env:"ID,notEmpty"`
		Secret      string `env:"SECRET,notEmpty"`
		RedirectURL string `env:"REDIRECT_URL,notEmpty"`
	} `envPrefix:"DELIVER_OIDC_"`
	Cookie struct {
		Secret string `env:"SECRET,notEmpty"`
	} `envPrefix:"DELIVER_COOKIE_"`
	MaxFileSize int64 `env:"DELIVER_MAX_FILE_SIZE" envDefault:"2000000000"`
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
