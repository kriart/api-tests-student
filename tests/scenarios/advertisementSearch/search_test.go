package advertisementsearch

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	advList "api-tests-template/internal/managers/advertisements"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

type TestSuite struct {
	suite.Suite
	loginData       models.LoginOkResponse
	advertisement   adModels.Advertisement
	searchBeginning string
	searchMiddle    string
	searchEnd       string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

// SetupSuite создаёт уникальное объявление и подстроки, независимые от ручных данных аккаунта.
func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	apiRunner.EnsureConnection(s.T())
	base.Precondition("Авторизация владельца тестовых объявлений")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	request := fixtures.AdvertisementRequest(s.T())
	unique := utils.RandomString(12)
	request.Fields["title"] = fixtures.TitlePrefix + unique + fixtures.TitleMiddle + unique + fixtures.TitleEnd
	s.advertisement = fixtures.CreateAdvertisement(s.T(), s.loginData.Token, request)
	s.searchBeginning = fixtures.TitlePrefix + unique
	s.searchMiddle = strings.TrimSpace(fixtures.TitleMiddle) + " " + unique
	s.searchEnd = unique + fixtures.TitleEnd
}

// TearDownSuite завершает набор сценариев.
func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestSearchByTitle() {
	s.T().Log("[STEP] Ищем тестовое объявление по подстроке из середины названия")
	list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, s.searchMiddle)
	require.NotEmpty(s.T(), list.Items, "Поиск не вернул тестовое объявление")
	assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.advertisement.Title)
	assertions.ContainsAdvertisementByID(s.T(), list.Items, s.advertisement.ID)
}

func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{name: "начало", search: s.searchBeginning},
		{name: "середина", search: s.searchMiddle},
		{name: "конец", search: s.searchEnd},
		{name: "нижний регистр", search: strings.ToLower(s.searchMiddle)},
		{name: "верхний регистр", search: strings.ToUpper(s.searchMiddle)},
		{name: "смешанный регистр", search: utils.MixCase(s.searchMiddle)},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.T().Logf("[STEP] Ищем по подстроке %q", tc.search)
			list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.advertisement.Title)
			assertions.ContainsAdvertisementByID(s.T(), list.Items, s.advertisement.ID)
		})
	}
}
