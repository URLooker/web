package api

import (
	"github.com/peng19940915/urlooker/web/model"
)

type StrategyResponse struct {
	Message string
	Data    []*model.Strategy
}

func (this *Web) GetStrategies(req interface{}, resp *StrategyResponse) error {
	strategies, err := model.GetAllStrategyByCron()
	if err != nil {
		resp.Message = err.Error()
	}
	resp.Data = strategies

	return nil
}

type PortStrategyResponse struct {
	Message string
	Data    []*model.PortStrategy
}

func (this *Web) GetPortStrategies(req interface{}, resp *PortStrategyResponse) error {
	strategies, err := model.GetAllPortStrategyByCron()
	if err != nil {
		resp.Message = err.Error()
	}
	resp.Data = strategies

	return nil
}

