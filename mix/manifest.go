package mix

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type Config struct {
	ManifestFile string
	PublicPath   string
}

type Manifest map[string]string

func New(c Config) (Manifest, error) {
	data, err := os.ReadFile(c.ManifestFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't read mix manifest '%s': %w", c.ManifestFile, err)
	}

	manifest := make(map[string]string)
	if err = json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("couldn't parse mix manifest '%s': %w", c.ManifestFile, err)
	}

	if c.PublicPath != "" {
		for asset, p := range manifest {
			manifest[asset] = path.Join(c.PublicPath, p)
		}
	}

	return manifest, nil
}

func (m Manifest) AssetPath(asset string) (string, error) {
	p, ok := m[asset]
	if !ok {
		return "", fmt.Errorf("asset '%s' not found in mix manifest", asset)
	}
	return p, nil
}
