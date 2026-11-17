package myAdvertisements

import (
	"testing"

	"api-tests-template/internal/client/http/myAdvertisements"
	"api-tests-template/internal/utils"
)

// GetMyAdvertisements возвращает список объявлений пользователя и проверяет статус.
func GetMyAdvertisements(t *testing.T, token string, status, limit int, search ...string) string {
	t.Helper()
	return utils.ReadBody(t, myAdvertisements.HttpGetMyAdvertisements(t, token, limit, search...), status)
}
