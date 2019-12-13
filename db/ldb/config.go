package ldb

type Config struct {
	Path        string `json:"path"`
	Compression bool   `json:"compression"`
	FileSize    int    `json:"filesize"`
	ReadOnly    bool   `json:"readonly"`
}

func (c *Config) Clone() *Config {
	return &Config{
		Path:        c.Path,
		Compression: c.Compression,
		FileSize:    c.FileSize,
		ReadOnly:    c.ReadOnly,
	}
}
