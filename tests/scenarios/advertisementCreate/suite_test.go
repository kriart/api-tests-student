package advertisementcreate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

// SetupSuite загружает настройки и авторизует владельца объявлений перед сценариями.
func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	apiRunner.EnsureConnection(s.T())
	base.Precondition("Авторизация владельца тестовых объявлений")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

// TearDownSuite завершает набор сценариев создания объявлений.
func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}
