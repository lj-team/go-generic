package daemon

type Config struct {
	PidFile string `json:"pid"`
	LogFile string `json:"log"`
	WordDir string `json:"dir"`
}
