package favorites

import (
	"net/http"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestFavoritesNonexistentAdvertisement() {
	ad := fixtures.CreateAdvertisement(s.T(), s.loginData.Token, fixtures.AdvertisementRequest(s.T()))
	advertisements.DeleteAdvertisement(s.T(), s.loginData.Token, ad.ID)

	s.Run("добавление несуществующего объявления", func() {
		s.T().Log("[STEP] Добавляем в избранное заведомо удалённое объявление")
		actual := manager.AddFavorite(s.T(), s.favoritesLoginData.Token, models.AddFavoriteRequest{AdvertisementID: ad.ID}, http.StatusNotFound)
		assertions.APIError(s.T(), actual, fixtures.NotFound, "Advertisement not found")
	})
	s.Run("удаление несуществующего объявления", func() {
		s.T().Log("[STEP] Удаляем из избранного заведомо удалённое объявление")
		actual := manager.DeleteFavorite(s.T(), s.favoritesLoginData.Token, ad.ID, http.StatusNotFound)
		assertions.APIError(s.T(), actual, fixtures.NotFound, "Favorite entry not found")
	})
}
