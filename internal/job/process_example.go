package job

import (
	"time"
)

type ProcessExample struct {
	Data string `json:"data"`
}

func (p *ProcessExample) Handle() error {
	// Process the job here.
	time.Sleep(1 * time.Second)
	return nil
}
