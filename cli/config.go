package cli

type Config struct {
	Production bool     `env:"PRODUCTION"`
	Host       string   `env:"HOST"`
	Port       int      `env:"PORT" envDefault:"3000"`
	Admins     []string `env:"ADMINS,notEmpty"`
	Storage    struct {
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

	Banner string `env:"BANNER"`
}

func (c *Config) AfterLoad() {
	if !c.Production && c.Banner == "" {
		c.Banner = "development"
	}
}