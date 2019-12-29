package account

type accountAux struct {
	fileName   string
	Name       string `yaml:"name"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Hostname   string `yaml:"hostname"`
	Port       int    `yaml:"port"`
	Debug      bool   `yaml:"debug"`
	SkipVerify bool   `yaml:"skipVerify"`
}

func (a accountAux) convert() Account {
	return Account{a.Name, a.Username, a.Password, a.Hostname, a.Port, a.Debug, a.SkipVerify}
}

func (a accountAux) IsEmpty() bool {
	return a.fileName == "" || a.Name == "" || a.Hostname == "" || a.Port == 0
}
