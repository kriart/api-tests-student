package fixtures

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/favorites"
	favoriteModels "api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"
)

// PhotoFiles возвращает пути к готовым PNG; четвёртый файл нужен для проверки превышения лимита.
func PhotoFiles(t *testing.T, count int) []string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok, "Не удалось определить папку с тестовыми фотографиями")
	files := make([]string, count)
	for i := range count {
		files[i] = filepath.Join(filepath.Dir(file), "testdata", fmt.Sprintf("photo-%d.png", i%MaxPhotos+1))
	}
	return files
}

// AdvertisementRequest формирует валидные поля с уникальным названием и одной фотографией.
func AdvertisementRequest(t *testing.T) models.CreateAdvertisementRequest {
	t.Helper()
	return models.CreateAdvertisementRequest{
		Fields: map[string]string{
			"title":       TitlePrefix + utils.RandomString(16) + TitleMiddle + utils.RandomString(8) + TitleEnd,
			"description": "Объявление для автоматического API-теста. Удаляется после проверки.",
			"price":       TestPriceText,
			"quantity":    TestQuantityText,
		},
		Photos: PhotoFiles(t, 1),
	}
}

// CreateAdvertisement создаёт фикстуру и сразу регистрирует удаление даже при падении теста.
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	t.Helper()
	t.Logf("[PREPARE] Создаём тестовое объявление %q", request.Fields["title"])
	body := advertisements.CreateAdvertisement(t, token, request, http.StatusCreated)
	ad := models.Advertisement{}
	require.NoError(t, json.Unmarshal([]byte(body), &ad))
	require.NotEmpty(t, ad.ID, "Не заполнен ID созданного объявления")
	t.Cleanup(func() { advertisements.DeleteAdvertisementIfExists(t, token, ad.ID) })
	t.Logf("[RESULT] Создано объявление %s", ad.ID)
	return ad
}

// AddFavorite добавляет тестовую запись и регистрирует удаление только этого избранного.
func AddFavorite(t *testing.T, token, id string) string {
	t.Helper()
	t.Cleanup(func() { favorites.DeleteFavoriteIfExists(t, token, id) })
	t.Logf("[STEP] Добавляем объявление %s в избранное", id)
	return favorites.AddFavorite(t, token, favoriteModels.AddFavoriteRequest{AdvertisementID: id}, http.StatusCreated)
}
