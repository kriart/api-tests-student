package advertisements

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

// GetAdvertisementsList возвращает результат поиска объявлений.
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	t.Helper()
	body := utils.ReadBody(t, advertisements.HttpGetAdvertisements(t, token, text), http.StatusOK)
	result := models.AdvertisementsList{}
	require.NoError(t, json.Unmarshal([]byte(body), &result))
	return result
}

// CreateAdvertisement создаёт объявление и проверяет ожидаемый статус.
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest, status int) string {
	t.Helper()
	raw := advertisements.HttpPostAdvertisement(t, token, request)
	if raw.StatusCode == http.StatusCreated && status != http.StatusCreated {
		body := utils.ReadBody(t, raw, http.StatusCreated)
		id := gjson.Get(body, "id").String()
		t.Cleanup(func() { DeleteAdvertisementIfExists(t, token, id) })
		t.Fatalf("Ожидался HTTP %d, но создано объявление %s", status, id)
	}
	return utils.ReadBody(t, raw, status)
}

// GetAdvertisement получает объявление по ID и проверяет статус.
func GetAdvertisement(t *testing.T, id string, status int) string {
	t.Helper()
	return utils.ReadBody(t, advertisements.HttpGetAdvertisement(t, id), status)
}

// DeleteAdvertisement удаляет объявление владельца.
func DeleteAdvertisement(t *testing.T, token, id string) {
	t.Helper()
	utils.ReadBody(t, advertisements.HttpDeleteAdvertisement(t, token, id), http.StatusNoContent)
}

// DeleteAdvertisementIfExists удаляет фикстуру, проверяя отсутствие уже удалённого объявления.
func DeleteAdvertisementIfExists(t *testing.T, token, id string) {
	t.Helper()
	t.Logf("[CLEANUP] Удаляем объявление %s", id)
	raw := advertisements.HttpDeleteAdvertisement(t, token, id)
	if raw.StatusCode == http.StatusNotFound || raw.StatusCode == http.StatusForbidden {
		utils.ReadBody(t, raw, raw.StatusCode)
		GetAdvertisement(t, id, http.StatusNotFound)
		return
	}
	utils.ReadBody(t, raw, http.StatusNoContent)
}

// GetAdvertisementPhotos получает фотографии объявления.
func GetAdvertisementPhotos(t *testing.T, token, id string) string {
	t.Helper()
	return utils.ReadBody(t, advertisements.HttpGetAdvertisementPhotos(t, token, id), http.StatusOK)
}

// DownloadPhoto возвращает содержимое загруженной фотографии.
func DownloadPhoto(t *testing.T, url string) []byte {
	t.Helper()
	return []byte(utils.ReadBody(t, advertisements.HttpGetPhoto(t, url), http.StatusOK))
}
