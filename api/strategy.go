package api

import (
	"github.com/urlooker/web/model"
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
