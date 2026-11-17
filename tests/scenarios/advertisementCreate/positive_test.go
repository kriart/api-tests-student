package advertisementcreate

import (
	"net/http"
	"os"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestCreateAdvertisementAllParameters() {
	request := fixtures.AdvertisementRequest(s.T())
	request.Photos = fixtures.PhotoFiles(s.T(), fixtures.MaxPhotos)
	created := fixtures.CreateAdvertisement(s.T(), s.loginData.Token, request)
	require.Len(s.T(), created.Photos, len(request.Photos), "В ответе создания отсутствуют фотографии")

	s.T().Log("[STEP] Проверяем сохранённые параметры объявления")
	body := advertisements.GetAdvertisement(s.T(), created.ID, http.StatusOK)
	require.Equal(s.T(), created.ID, gjson.Get(body, "id").String())
	require.Equal(s.T(), s.loginData.User.Id, gjson.Get(body, "user_id").String())
	assertions.AdvertisementMatchesRequest(s.T(), body, request)

	s.T().Log("[STEP] Ищем объявление по полному названию")
	list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, created.Title)
	assertions.ContainsAdvertisementByID(s.T(), list.Items, created.ID)

	s.T().Log("[STEP] Проверяем фотографии и содержимое файлов на сервере")
	body = advertisements.GetAdvertisementPhotos(s.T(), s.loginData.Token, created.ID)
	photos := gjson.Parse(body).Array()
	require.Len(s.T(), photos, len(request.Photos), "Сохранены не все фотографии")
	for i, photo := range photos {
		require.Equal(s.T(), created.Photos[i].ID, photo.Get("id").String())
		require.EqualValues(s.T(), i, photo.Get("sort_order").Int())
		require.NotEmpty(s.T(), photo.Get("url").String(), "Не заполнен URL фотографии")
		expected, err := os.ReadFile(request.Photos[i])
		require.NoError(s.T(), err)
		require.Equal(s.T(), expected, advertisements.DownloadPhoto(s.T(), photo.Get("url").String()), "Не совпадает содержимое фотографии")
	}
	s.T().Log("[RESULT] Объявление найдено, параметры и фотографии сохранены")
}
