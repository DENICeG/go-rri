package messagetype

type LoginMessage struct {
	Version  string `yaml:"Version"`
	Action   string `yaml:"Action"`
	Username string `yaml:"User"`
	Password string `yaml:"Password"`
}

type LogoutMessage struct {
	Version string `yaml:"Version"`
	Action  string `yaml:"Action"`
}
