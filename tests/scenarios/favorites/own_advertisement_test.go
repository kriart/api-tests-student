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

func (s *TestSuite) TestCannotFavoriteOwnAdvertisement() {
	ad := fixtures.CreateAdvertisement(s.T(), s.loginData.Token, fixtures.AdvertisementRequest(s.T()))
	s.T().Log("[STEP] Добавляем собственное объявление в избранное")
	actual := manager.AddFavorite(s.T(), s.loginData.Token, models.AddFavoriteRequest{AdvertisementID: ad.ID}, http.StatusForbidden)
	assertions.APIError(s.T(), actual, fixtures.Forbidden, "Cannot add own advertisement to favorites")
	list := manager.GetFavorites(s.T(), s.loginData.Token, ad.Title, http.StatusOK)
	assertions.NotContainsFavorite(s.T(), list, ad.ID)
	require.Zero(s.T(), gjson.Get(list, "total").Int(), "Собственное объявление попало в избранное")
}
