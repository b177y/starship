package nebutils

type Config struct {
	Pki struct {
		Ca   string `yaml:"ca"`
		Cert string `yaml:"cert"`
		Key  string `yaml:"key"`
	} `yaml:"pki"`
	StaticHostmap struct {
	}
}
