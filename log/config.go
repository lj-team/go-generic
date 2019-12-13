package log

type Config struct {
	Template string `json:"template"`
	Period   int    `json:"period"`
	Save     int    `json:"save"`
	Level    string `json:"level"`
	StdOut   bool   `json:"stdout"`
	StdErr   bool   `json:"stderr"`
}
