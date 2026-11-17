package favorites

import (
	"net/http"

	"api-tests-template/internal/assertions"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestFavoritesUnauthorized() {
	ad := fixtures.CreateAdvertisement(s.T(), s.loginData.Token, fixtures.AdvertisementRequest(s.T()))
	fixtures.AddFavorite(s.T(), s.favoritesLoginData.Token, ad.ID)

	cases := []struct {
		name    string
		token   string
		message string
	}{
		{name: "без токена", message: fixtures.MissingAuthMessage},
		{name: "неверный токен", token: fixtures.InvalidToken, message: fixtures.InvalidAuthMessage},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.Run("получение списка", func() {
				s.T().Log("[STEP] GET /favorites без действующей авторизации")
				assertions.APIError(s.T(), manager.GetFavorites(s.T(), tc.token, "", http.StatusUnauthorized), fixtures.Unauthorized, tc.message)
			})
			s.Run("добавление", func() {
				s.T().Log("[STEP] POST /favorites без действующей авторизации")
				actual := manager.AddFavorite(s.T(), tc.token, models.AddFavoriteRequest{AdvertisementID: ad.ID}, http.StatusUnauthorized)
				assertions.APIError(s.T(), actual, fixtures.Unauthorized, tc.message)
			})
			s.Run("удаление", func() {
				s.T().Log("[STEP] DELETE /favorites без действующей авторизации")
				assertions.APIError(s.T(), manager.DeleteFavorite(s.T(), tc.token, ad.ID, http.StatusUnauthorized), fixtures.Unauthorized, tc.message)
			})
			s.Run("пакетное добавление", func() {
				s.T().Log("[STEP] POST /favorites/bulk без действующей авторизации")
				assertions.APIError(s.T(), manager.BulkAdd(s.T(), tc.token, ad.ID, http.StatusUnauthorized), fixtures.Unauthorized, tc.message)
			})
		})
	}
	list := manager.GetFavorites(s.T(), s.favoritesLoginData.Token, ad.Title, http.StatusOK)
	assertions.ContainsFavorite(s.T(), list, ad.ID)
}
