package utils

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// ReadBody проверяет статус, читает тело ответа и закрывает его.
func ReadBody(t *testing.T, response *http.Response, status int) string {
	t.Helper()
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать ответ")
	require.Equal(t, status, response.StatusCode, "Неожиданный HTTP-статус: %s", body)
	if status == http.StatusNoContent {
		require.Empty(t, body, "Ожидалось пустое тело ответа")
	}
	return string(body)
}
