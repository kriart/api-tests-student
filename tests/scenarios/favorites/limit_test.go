package favorites

import (
	"fmt"
	"net/http"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/assertions"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestFavoritesLimit() {
	before := manager.GetFavorites(s.T(), s.favoritesLoginData.Token, "", http.StatusOK)
	require.LessOrEqual(s.T(), gjson.Get(before, "total").Int(), int64(fixtures.MaxFavorites), "Исходное избранное уже превышает лимит")
	test := s.T()
	test.Cleanup(func() {
		restored := manager.GetFavorites(test, s.favoritesLoginData.Token, "", http.StatusOK)
		require.Equal(test, gjson.Get(before, "total").Int(), gjson.Get(restored, "total").Int(), "Очистка не восстановила исходное число записей избранного")
		test.Log("[RESULT] Исходное избранное восстановлено после проверки лимита")
	})
	request := fixtures.AdvertisementRequest(s.T())

	s.T().Logf("[PREPARE] Заполняем избранное до лимита %d, исходное число: %d", fixtures.MaxFavorites, gjson.Get(before, "total").Int())
	for i := gjson.Get(before, "total").Int(); i < int64(fixtures.MaxFavorites); i++ {
		request.Fields["title"] = fixtures.TitlePrefix + utils.RandomString(24)
		ad := fixtures.CreateAdvertisement(s.T(), s.loginData.Token, request)
		fixtures.AddFavorite(s.T(), s.favoritesLoginData.Token, ad.ID)
	}
	full := manager.GetFavorites(s.T(), s.favoritesLoginData.Token, "", http.StatusOK)
	require.EqualValues(s.T(), fixtures.MaxFavorites, gjson.Get(full, "total").Int(), "Не удалось подготовить лимит избранного")

	request.Fields["title"] = fixtures.TitlePrefix + utils.RandomString(24)
	extra := fixtures.CreateAdvertisement(s.T(), s.loginData.Token, request)
	s.T().Log("[STEP] Добавляем новое чужое объявление сверх лимита")
	actual := manager.AddFavorite(s.T(), s.favoritesLoginData.Token, models.AddFavoriteRequest{AdvertisementID: extra.ID}, http.StatusBadRequest)
	assertions.APIError(s.T(), actual, fixtures.LimitReached, fmt.Sprintf("Favorites limit of %d reached", fixtures.MaxFavorites))
	after := manager.GetFavorites(s.T(), s.favoritesLoginData.Token, "", http.StatusOK)
	require.Equal(s.T(), gjson.Get(full, "total").Int(), gjson.Get(after, "total").Int(), "Запрос сверх лимита изменил число записей")
	assertions.NotContainsFavorite(s.T(), after, extra.ID)
}
