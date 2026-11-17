package advertisements

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/steinfletcher/apitest"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpGetAdvertisements выполняет поиск объявлений по подстроке в названии.
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).Status(http.StatusOK).
		End().Response
}

// HttpPostAdvertisement отправляет поля и фотографии объявления в multipart-запросе.
func HttpPostAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) *http.Response {
	r := apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath)
	for name, value := range request.Fields {
		r.MultipartFormData(name, value)
	}
	for _, photo := range request.Photos {
		r.MultipartFile("photos", photo)
	}
	return r.Expect(t).End().Response
}

// HttpGetAdvertisement получает объявление по идентификатору из query-параметра.
func HttpGetAdvertisement(t *testing.T, id string) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

// HttpDeleteAdvertisement удаляет объявление от имени его владельца.
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

// HttpGetAdvertisementPhotos получает список фотографий объявления.
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(fmt.Sprintf(path.AdvertisementPhotosPath, id)).
		Expect(t).End().Response
}

// HttpGetPhoto скачивает файл фотографии по URL.
func HttpGetPhoto(t *testing.T, url string) *http.Response {
	return apitest.New().EnableNetworking(&http.Client{Timeout: 30 * time.Second}).Get(url).Expect(t).End().Response
}
