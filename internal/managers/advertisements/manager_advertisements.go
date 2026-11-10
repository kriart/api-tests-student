package advertisements

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	advertisementsListRawResponse := advertisements.HttpGetAdvertisements(t, token, text)
	defer advertisementsListRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(advertisementsListRawResponse.Body)
	require.NoError(t, errResponse)

	advertisementsListResponse := models.AdvertisementsList{}
	err := json.Unmarshal(respArray, &advertisementsListResponse)
	require.NoError(t, err)

	return advertisementsListResponse
}
