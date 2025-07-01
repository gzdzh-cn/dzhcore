package defineType

import "time"

// 运行日志
type OutputsForLogger struct {
	Time       time.Time `json:"time"`
	Host       string    `json:"host"`
	RequestURI string    `json:"requestURI"`
	Params     string    `json:"params"`
	RunTime    float64   `json:"runTime"`
	Prefix     string    `json:"prefix"`
	Suffix     string    `json:"suffix"`
	File       string    `json:"file"`
	FileRule   string    `json:"fileRule"`
	RotateSize string    `json:"rotateSize"`
	Stdout     bool      `json:"stdout"`
	Path       string    `json:"path"`
	Throughput float64   `json:"throughput"`
	MemUsed    uint64    `json:"memUsed"`
}

// 运行日志配置
type RunLogger struct {
	Path       string `json:"path"`
	Enable     bool   `json:"enable"`
	File       string `json:"file"`
	RotateSize string `json:"rotateSize"`
	Stdout     bool   `json:"stdout"`
}
