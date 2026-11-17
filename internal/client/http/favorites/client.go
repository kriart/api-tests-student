package favorites

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpPostFavorite отправляет JSON-запрос на добавление объявления в избранное.
func HttpPostFavorite(t *testing.T, token string, request any) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		JSON(request).Expect(t).End().Response
}

// HttpGetFavorites получает страницу избранного с фильтром по названию.
func HttpGetFavorites(t *testing.T, token, search string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("search", search).Query("page_size", "50").
		Expect(t).End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного текущего пользователя.
func HttpDeleteFavorite(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(fmt.Sprintf(path.FavoritePath, id)).
		Expect(t).End().Response
}

// HttpPostFavoritesBulk отправляет запрос объединения избранного для проверки авторизации.
func HttpPostFavoritesBulk(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesBulkPath).
		JSON(map[string][]string{"advertisement_ids": {id}}).Expect(t).End().Response
}
