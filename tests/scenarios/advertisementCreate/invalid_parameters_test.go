package advertisementcreate

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

const (
	titleRequired       = "title is required"
	descriptionRequired = "description is required"
	invalidPrice        = "price must be between 0 and 1000000 inclusive"
	invalidQuantity     = "quantity must be between 1 and 100 inclusive"
)

func (s *TestSuite) TestCreateAdvertisementInvalidParameters() {
	cases := []struct {
		name    string
		field   string
		value   string
		omit    bool
		photos  int
		message string
	}{
		{name: "нет названия", field: "title", omit: true, photos: 1, message: titleRequired},
		{name: "пустое название", field: "title", value: "", photos: 1, message: titleRequired},
		{
			name: "название длиннее допустимого", field: "title",
			value: strings.Repeat("я", fixtures.MaxTitleLength+1), photos: 1,
			message: fmt.Sprintf("title must not exceed %d characters provided %d", fixtures.MaxTitleLength, fixtures.MaxTitleLength+1),
		},
		{name: "нет описания", field: "description", omit: true, photos: 1, message: descriptionRequired},
		{name: "пустое описание", field: "description", value: "", photos: 1, message: descriptionRequired},
		{
			name: "описание длиннее допустимого", field: "description",
			value: strings.Repeat("я", fixtures.MaxDescriptionLength+1), photos: 1,
			message: "description must be between 1 and 500 characters",
		},
		{name: "отрицательная цена", field: "price", value: "-1", photos: 1, message: invalidPrice},
		{name: "слишком большая цена", field: "price", value: strconv.Itoa(fixtures.MaxPrice + 1), photos: 1, message: invalidPrice},
		{name: "цена не число", field: "price", value: "not-a-price", photos: 1, message: "cannot convert string value not-a-price to float"},
		{name: "нулевое количество", field: "quantity", value: "0", photos: 1, message: invalidQuantity},
		{name: "отрицательное количество", field: "quantity", value: "-1", photos: 1, message: invalidQuantity},
		{name: "слишком большое количество", field: "quantity", value: strconv.Itoa(fixtures.MaxQuantity + 1), photos: 1, message: invalidQuantity},
		{name: "количество не число", field: "quantity", value: "abc", photos: 1, message: "cannot convert string value to int"},
		{name: "нет фотографий", photos: 0, message: "at least 1 photo is required"},
		{name: "слишком много фотографий", photos: fixtures.MaxPhotos + 1, message: "maximum 3 photos allowed"},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			request := fixtures.AdvertisementRequest(s.T())
			if tc.field != "" {
				if tc.omit {
					delete(request.Fields, tc.field)
				} else {
					request.Fields[tc.field] = tc.value
				}
			}
			request.Photos = fixtures.PhotoFiles(s.T(), tc.photos)
			s.T().Logf("[STEP] Создаём объявление с невалидными параметрами: %s", tc.name)
			actual := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, request, http.StatusBadRequest)
			assertions.APIError(s.T(), actual, fixtures.InvalidRequest, tc.message)
		})
	}
}
