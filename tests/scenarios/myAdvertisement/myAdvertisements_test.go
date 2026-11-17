package myAdvertisement

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/assertions"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/myAdvertisements"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

type TestSuite struct {
	suite.Suite
	loginData     models.LoginOkResponse
	request       adModels.CreateAdvertisementRequest
	advertisement adModels.Advertisement
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

// SetupSuite подготавливает собственное объявление с автоматической очисткой.
func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	apiRunner.EnsureConnection(s.T())
	base.Precondition("Авторизация владельца тестовых объявлений")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	s.request = fixtures.AdvertisementRequest(s.T())
	s.advertisement = fixtures.CreateAdvertisement(s.T(), s.loginData.Token, s.request)
}

// TearDownSuite завершает набор сценариев.
func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestGetMyAdvertisementsPositive() {
	s.T().Log("[STEP] Получаем собственное тестовое объявление по уникальному названию")
	body := myAdvertisements.GetMyAdvertisements(s.T(), s.loginData.Token, http.StatusOK, 50, s.advertisement.Title)
	item := gjson.Get(body, `items.#(id=="`+s.advertisement.ID+`")`)
	require.True(s.T(), item.Exists(), "Созданное объявление отсутствует в профиле")
	require.Equal(s.T(), s.loginData.User.Id, item.Get("user_id").String())
	assertions.AdvertisementMatchesRequest(s.T(), item.Raw, s.request)
}

func (s *TestSuite) TestGetMyAdvertisementsIncorrectToken() {
	s.T().Log("[STEP] Запрашиваем собственные объявления с неверным токеном")
	body := myAdvertisements.GetMyAdvertisements(s.T(), fixtures.InvalidToken, http.StatusUnauthorized, 50)
	assertions.APIError(s.T(), body, fixtures.Unauthorized, fixtures.InvalidAuthMessage)
}
