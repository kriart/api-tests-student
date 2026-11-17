package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/utils"
)

// AddFavorite добавляет объявление в избранное и проверяет статус.
func AddFavorite(t *testing.T, token string, request any, status int) string {
	t.Helper()
	return utils.ReadBody(t, favorites.HttpPostFavorite(t, token, request), status)
}

// GetFavorites возвращает избранное с фильтром по названию.
func GetFavorites(t *testing.T, token, search string, status int) string {
	t.Helper()
	body := utils.ReadBody(t, favorites.HttpGetFavorites(t, token, search), status)
	if status == http.StatusOK {
		require.True(t, gjson.Get(body, "items").IsArray(), "В ответе нет массива избранного")
	}
	return body
}

// DeleteFavorite удаляет объявление из избранного и проверяет статус.
func DeleteFavorite(t *testing.T, token, id string, status int) string {
	t.Helper()
	return utils.ReadBody(t, favorites.HttpDeleteFavorite(t, token, id), status)
}

// DeleteFavoriteIfExists удаляет тестовую запись избранного, допуская её отсутствие.
func DeleteFavoriteIfExists(t *testing.T, token, id string) {
	t.Helper()
	raw := favorites.HttpDeleteFavorite(t, token, id)
	if raw.StatusCode == http.StatusNotFound {
		utils.ReadBody(t, raw, http.StatusNotFound)
		return
	}
	utils.ReadBody(t, raw, http.StatusNoContent)
}

// BulkAdd проверяет ответ пакетного добавления в избранное.
func BulkAdd(t *testing.T, token, id string, status int) string {
	t.Helper()
	return utils.ReadBody(t, favorites.HttpPostFavoritesBulk(t, token, id), status)
}
