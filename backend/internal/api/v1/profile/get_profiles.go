package profile

import (
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/buffi-buchi/invest-compass/backend/internal/api"
)

func (i *Implementation) GetProfiles(w http.ResponseWriter, r *http.Request) {
	userID := uuid.New()

	// TODO: Implement authentication.
	// TODO: Get userID from context.

	profiles, err := i.service.GetProfilesByUserID(r.Context(), userID)
	if err != nil {
		i.logger.Errorw("Get profiles by user ID", zap.Error(err))

		api.EncodeErrorf(w, http.StatusInternalServerError, "Get profiles by user ID")

		return
	}

	response := GetProfilesResponse{
		Profiles: make([]Profile, 0, len(profiles)),
	}

	for _, profile := range profiles {
		response.Profiles = append(response.Profiles, Profile{
			UserId: profile.UserID,
			Ticker: profile.Ticker,
		})
	}

	api.EncodeSuccess(w, http.StatusOK, response)
}
