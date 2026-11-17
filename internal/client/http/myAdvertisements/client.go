package myAdvertisements

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpGetMyAdvertisements выполняет GET-запрос для получения списка объявлений авторизованного пользователя
func HttpGetMyAdvertisements(t *testing.T, token string, limit int, search ...string) *http.Response {
	r := apiRunner.GetRunner().Auth(token).Create().Get(path.MyAdvertisementsPath).
		Query("limit", strconv.Itoa(limit))
	if len(search) > 0 {
		r.Query("search", search[0])
	}
	return r.Expect(t).End().Response
}
