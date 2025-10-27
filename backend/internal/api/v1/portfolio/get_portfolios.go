package portfolio

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/buffi-buchi/invest-compass/backend/internal/api"
	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

func (i *Implementation) GetPortfolios(w http.ResponseWriter, r *http.Request) {
	claims, _ := model.AuthClaimsValue(r.Context())

	// TODO: Use limit offset.

	portfolios, err := i.service.GeByUserID(r.Context(), claims.UserID, 10, 0)
	if err != nil {
		i.logger.Error("Get portfolios by user ID", zap.Error(err))

		api.EncodeErrorf(w, http.StatusInternalServerError, "Get portfolios by user ID")

		return
	}

	response := GetPortfoliosResponse{
		Portfolios: make([]Portfolio, 0, len(portfolios)),
	}

	for _, portfolio := range portfolios {
		response.Portfolios = append(response.Portfolios, Portfolio{
			Id:         portfolio.ID,
			UserId:     portfolio.UserID,
			Name:       portfolio.Name,
			CreateTime: portfolio.CreateTime,
		})
	}

	api.EncodeSuccess(w, http.StatusOK, response)
}
