package daemon

type Config struct {
	PidFile string `json:"pid"`
	LogFile string `json:"log"`
	WorkDir string `json:"dir"`
}
