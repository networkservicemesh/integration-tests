package exec

import "time"

type Config struct {
	Timeout time.Duration `default:"1m" desc:"Timeout to wait for successful finishing of command." split_words:"true"`
}
