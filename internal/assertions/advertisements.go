package assertions

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/advertisements/models"
)

// ContainsAdvertisementByTitle проверяет полное название среди найденных объявлений.
func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	titles := make([]string, len(items))
	for i, item := range items {
		titles[i] = item.Title
	}

	require.Contains(t, titles, title)
}

// ContainsAdvertisementByID проверяет, что поиск вернул именно созданное в тесте объявление.
func ContainsAdvertisementByID(t *testing.T, items []models.Advertisement, id string) {
	t.Helper()
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	require.Contains(t, ids, id, "Созданное объявление отсутствует в результатах поиска")
	t.Logf("[RESULT] Объявление %s найдено в поиске", id)
}

// APIError проверяет код и сообщение об ошибке.
func APIError(t *testing.T, body, code, message string) {
	t.Helper()
	require.Equal(t, code, gjson.Get(body, "error").String(), "Неверный код ошибки")
	require.Equal(t, message, gjson.Get(body, "message").String(), "Неверное сообщение об ошибке")
	t.Logf("[RESULT] Получена ошибка %s: %s", code, message)
}

// AdvertisementMatchesRequest сравнивает сохранённые поля объявления с запросом.
func AdvertisementMatchesRequest(t *testing.T, body string, request models.CreateAdvertisementRequest) {
	t.Helper()
	price, err := strconv.ParseFloat(request.Fields["price"], 64)
	require.NoError(t, err)
	quantity, err := strconv.Atoi(request.Fields["quantity"])
	require.NoError(t, err)
	require.Equal(t, request.Fields["title"], gjson.Get(body, "title").String(), "Не совпадает название")
	require.Equal(t, request.Fields["description"], gjson.Get(body, "description").String(), "Не совпадает описание")
	require.InDelta(t, price, gjson.Get(body, "price").Float(), 0.001, "Не совпадает цена")
	require.EqualValues(t, quantity, gjson.Get(body, "quantity").Int(), "Не совпадает количество")
	require.Len(t, gjson.Get(body, "photos").Array(), len(request.Photos), "Сохранены не все фотографии")
	t.Log("[RESULT] Параметры объявления совпадают с запросом")
}

// ContainsFavorite проверяет наличие объявления в списке избранного.
func ContainsFavorite(t *testing.T, body, id string) {
	t.Helper()
	require.True(t, gjson.Get(body, `items.#(advertisement.id=="`+id+`")`).Exists(), "В избранном нет объявления %s", id)
}

// NotContainsFavorite проверяет отсутствие объявления в списке избранного.
func NotContainsFavorite(t *testing.T, body, id string) {
	t.Helper()
	require.False(t, gjson.Get(body, `items.#(advertisement.id=="`+id+`")`).Exists(), "Объявление %s осталось в избранном", id)
}
