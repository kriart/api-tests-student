package favorites

import (
	"net/http"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/assertions"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestCannotFavoriteAdvertisementTwice() {
	ad := fixtures.CreateAdvertisement(s.T(), s.loginData.Token, fixtures.AdvertisementRequest(s.T()))
	fixtures.AddFavorite(s.T(), s.favoritesLoginData.Token, ad.ID)

	s.T().Log("[STEP] Повторно добавляем то же объявление в избранное")
	actual := manager.AddFavorite(s.T(), s.favoritesLoginData.Token, models.AddFavoriteRequest{AdvertisementID: ad.ID}, http.StatusConflict)
	assertions.APIError(s.T(), actual, fixtures.AlreadyExists, "Advertisement is already in favorites")
	list := manager.GetFavorites(s.T(), s.favoritesLoginData.Token, ad.Title, http.StatusOK)
	assertions.ContainsFavorite(s.T(), list, ad.ID)
	require.EqualValues(s.T(), 1, gjson.Get(list, "total").Int(), "После повторного добавления появились дубликаты")
}
