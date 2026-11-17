package favorites

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData          models.LoginOkResponse
	favoritesLoginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

// SetupSuite авторизует владельца объявлений и отдельного пользователя избранного.
func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	apiRunner.EnsureConnection(s.T())
	base.Precondition("Авторизация владельца тестовых объявлений")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	base.Precondition("Авторизация пользователя избранного")
	require.NotEmpty(s.T(), os.Getenv("FAVORITES_LOGIN"), "Заполните FAVORITES_LOGIN в .env.override")
	require.NotEmpty(s.T(), os.Getenv("FAVORITES_PASSWORD"), "Заполните FAVORITES_PASSWORD в .env.override")
	s.favoritesLoginData = auth.Login(s.T(), os.Getenv("FAVORITES_LOGIN"), os.Getenv("FAVORITES_PASSWORD"))
	require.NotEqual(s.T(), s.loginData.User.Id, s.favoritesLoginData.User.Id, "Для избранного нужны два разных пользователя")
}

// TearDownSuite завершает набор сценариев избранного.
func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}
