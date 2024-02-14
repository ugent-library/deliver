package cli

import "fmt"

//go:generate go run github.com/g4s8/envdoc@v0.1.2 --output ../CONFIG.md --all

// Version info
type Version struct {
	Branch string `json:"branch" env:"SOURCE_BRANCH"`
	Commit string `json:"commit" env:"SOURCE_COMMIT"`
	Image  string `json:"image" env:"IMAGE_NAME"`
}

type Config struct {
	// Env must be local, development, test or production
	Env      string   `env:"ENV" envDefault:"production"`
	Timezone string   `env:"TIMEZONE" envDefault:"Europe/Brussels"`
	Host     string   `env:"HOST"`
	Port     int      `env:"PORT" envDefault:"3000"`
	Admins   []string `env:"ADMINS,notEmpty"`
	Storage  struct {
		Backend string `env:"BACKEND" envDefault:"s3"`
		Conn    string `env:"CONN,notEmpty"`
	} `envPrefix:"STORAGE_"`
	Repo struct {
		Conn string `env:"CONN,notEmpty"`
	} `envPrefix:"REPO_"`
	OIDC struct {
		URL         string `env:"URL,notEmpty"`
		ID          string `env:"ID,notEmpty"`
		Secret      string `env:"SECRET,notEmpty"`
		RedirectURL string `env:"REDIRECT_URL,notEmpty"`
	} `envPrefix:"OIDC_"`
	Cookie struct {
		Secret string `env:"SECRET,notEmpty"`
	} `envPrefix:"COOKIE_"`
	MaxFileSize int64 `env:"MAX_FILE_SIZE" envDefault:"2000000000"`
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
