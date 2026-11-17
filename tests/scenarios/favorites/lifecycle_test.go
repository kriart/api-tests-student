package favorites

import (
	"net/http"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestFavoritesLifecycle() {
	ad := fixtures.CreateAdvertisement(s.T(), s.loginData.Token, fixtures.AdvertisementRequest(s.T()))

	favorite := fixtures.AddFavorite(s.T(), s.favoritesLoginData.Token, ad.ID)
	require.NotEmpty(s.T(), gjson.Get(favorite, "favorite_id").String(), "Сервер не вернул ID записи избранного")
	require.NotEmpty(s.T(), gjson.Get(favorite, "added_at").String(), "Не заполнено время добавления в избранное")
	require.Equal(s.T(), ad.ID, gjson.Get(favorite, "advertisement.id").String(), "В избранное добавлено другое объявление")

	s.T().Log("[STEP] Проверяем появление объявления в списке избранного")
	list := manager.GetFavorites(s.T(), s.favoritesLoginData.Token, ad.Title, http.StatusOK)
	assertions.ContainsFavorite(s.T(), list, ad.ID)
	require.EqualValues(s.T(), 1, gjson.Get(list, "total").Int(), "По уникальному названию должна быть одна запись")

	s.T().Log("[STEP] Удаляем объявление из избранного")
	manager.DeleteFavorite(s.T(), s.favoritesLoginData.Token, ad.ID, http.StatusNoContent)
	list = manager.GetFavorites(s.T(), s.favoritesLoginData.Token, ad.Title, http.StatusOK)
	assertions.NotContainsFavorite(s.T(), list, ad.ID)
	require.Zero(s.T(), gjson.Get(list, "total").Int(), "После удаления в избранном осталась запись")

	s.T().Log("[STEP] Повторно добавляем объявление в избранное и удаляем само объявление")
	fixtures.AddFavorite(s.T(), s.favoritesLoginData.Token, ad.ID)
	advertisements.DeleteAdvertisement(s.T(), s.loginData.Token, ad.ID)
	actual := advertisements.GetAdvertisement(s.T(), ad.ID, http.StatusNotFound)
	assertions.APIError(s.T(), actual, fixtures.NotFound, "Advertisement not found")

	list = manager.GetFavorites(s.T(), s.favoritesLoginData.Token, ad.Title, http.StatusOK)
	assertions.NotContainsFavorite(s.T(), list, ad.ID)
	require.Zero(s.T(), gjson.Get(list, "total").Int(), "Избранное не очищено после удаления самого объявления")
	s.T().Log("[RESULT] Полный цикл избранного и каскадное удаление проверены")
}
