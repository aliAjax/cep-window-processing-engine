package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr                 string
	CheckpointDir        string
	MaxPatternComplexity int
	AllowedLateness      time.Duration
	QuotaPerSecond       float64
}

func Load() Config {
	c := Config{Addr: ":8080", CheckpointDir: "./data/checkpoints", MaxPatternComplexity: 100, AllowedLateness: 30 * time.Second, QuotaPerSecond: 100}
	if v := os.Getenv("CEP_ADDR"); v != "" {
		c.Addr = v
	}
	if v := os.Getenv("CEP_CHECKPOINT_DIR"); v != "" {
		c.CheckpointDir = v
	}
	if v := os.Getenv("CEP_MAX_COMPLEXITY"); v != "" {
		if n, e := strconv.Atoi(v); e == nil {
			c.MaxPatternComplexity = n
		}
	}
	return c
}
func (c Config) Validate() error {
	if c.Addr == "" {
		return fmt.Errorf("addr required")
	}
	if c.MaxPatternComplexity < 1 {
		return fmt.Errorf("complexity invalid")
	}
	if c.QuotaPerSecond <= 0 {
		return fmt.Errorf("quota invalid")
	}
	return nil
}
