package advertisementcreate

import (
	"net/http"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	cases := []struct {
		name    string
		token   string
		message string
	}{
		{name: "без токена", message: fixtures.MissingAuthMessage},
		{name: "неверный токен", token: fixtures.InvalidToken, message: fixtures.InvalidAuthMessage},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			request := fixtures.AdvertisementRequest(s.T())
			s.T().Logf("[STEP] Пробуем создать валидное объявление: %s", tc.name)
			actual := advertisements.CreateAdvertisement(s.T(), tc.token, request, http.StatusUnauthorized)
			assertions.APIError(s.T(), actual, fixtures.Unauthorized, tc.message)
		})
	}
}
