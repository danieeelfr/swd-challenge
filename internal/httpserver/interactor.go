package app

import (
	"github.com/danieeelfr/swd-challenge/internal/config"
	"github.com/danieeelfr/swd-challenge/internal/httpserver/api"
	"github.com/danieeelfr/swd-challenge/pkg/wait"
)

// Interactor to expose the http server methods
type Interactor interface {
	Start() error
	Shutdown()
}

// New http server
func New(cfg *config.Config, wg *wait.Wait) (Interactor, error) {
	httpSrv, err := api.New(cfg, wg)
	if err != nil {
		return nil, err
	}

	return httpSrv, nil
}
