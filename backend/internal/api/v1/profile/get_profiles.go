package profile

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/buffi-buchi/invest-compass/backend/internal/api"
	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

func (i *Implementation) GetProfiles(w http.ResponseWriter, r *http.Request) {
	claims, _ := model.AuthClaimsValue(r.Context())

	// TODO: Use limit offset.

	profiles, err := i.service.GeByUserID(r.Context(), claims.UserID, 10, 0)
	if err != nil {
		i.logger.Error("Get profiles by user ID", zap.Error(err))

		api.EncodeErrorf(w, http.StatusInternalServerError, "Get profiles by user ID")

		return
	}

	response := GetProfilesResponse{
		Profiles: make([]Profile, 0, len(profiles)),
	}

	for _, profile := range profiles {
		response.Profiles = append(response.Profiles, Profile{
			Id:         profile.ID,
			UserId:     profile.UserID,
			Name:       profile.Name,
			CreateTime: profile.CreateTime,
		})
	}

	api.EncodeSuccess(w, http.StatusOK, response)
}
