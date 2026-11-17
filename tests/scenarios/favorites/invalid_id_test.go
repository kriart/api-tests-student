package favorites

import (
	"net/http"

	"api-tests-template/internal/assertions"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestFavoritesInvalidAdvertisementID() {
	cases := []struct {
		name    string
		request any
		message string
	}{
		{
			name: "нет advertisement_id", request: map[string]string{},
			message: "Key: 'AddFavoriteRequest.AdvertisementID' Error:Field validation for 'AdvertisementID' failed on the 'required' tag",
		},
		{name: "пустой advertisement_id", request: models.AddFavoriteRequest{}, message: "invalid UUID length: 0"},
		{
			name: "неверный формат advertisement_id", request: models.AddFavoriteRequest{AdvertisementID: "not-a-uuid"},
			message: "invalid UUID length: 10",
		},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.T().Logf("[STEP] POST /favorites: %s", tc.name)
			actual := manager.AddFavorite(s.T(), s.favoritesLoginData.Token, tc.request, http.StatusBadRequest)
			assertions.APIError(s.T(), actual, fixtures.InvalidRequest, tc.message)
		})
	}
	s.Run("удаление с неверным форматом ID", func() {
		s.T().Log("[STEP] DELETE /favorites с невалидным UUID")
		actual := manager.DeleteFavorite(s.T(), s.favoritesLoginData.Token, "not-a-uuid", http.StatusBadRequest)
		assertions.APIError(s.T(), actual, fixtures.InvalidRequest, "Invalid advertisement ID")
	})
}
