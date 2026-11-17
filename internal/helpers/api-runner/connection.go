package api_runner

import (
	"errors"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"api-tests-template/internal/utils"

	"github.com/stretchr/testify/require"
)

const (
	connectionAttempts = 5
	connectionTimeout  = 5 * time.Second
	defaultAPIURL      = "https://testboard.avito.com/api/v1"
)

// EnsureConnection проверяет доступность TestBoard перед набором; повторяет только сетевые таймауты безопасного GET.
func EnsureConnection(t *testing.T) {
	t.Helper()
	if strings.TrimRight(os.Getenv("API_URL"), "/") != defaultAPIURL {
		return
	}
	client := &http.Client{Timeout: connectionTimeout}
	var lastErr error
	for attempt := 1; attempt <= connectionAttempts; attempt++ {
		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, defaultAPIURL+"/advertisements?limit=1", nil)
		require.NoError(t, err, "Не удалось подготовить проверку соединения")
		request.Header.Set("User-Agent", defaultUserAgent)
		raw, err := client.Do(request)
		if err != nil {
			var timeout net.Error
			if errors.As(err, &timeout) && timeout.Timeout() && t.Context().Err() == nil {
				lastErr = err
				t.Logf("[Precondition] Таймаут подключения, попытка %d/%d", attempt, connectionAttempts)
				continue
			}
			require.NoError(t, err, "Не удалось подключиться к API")
		}
		utils.ReadBody(t, raw, http.StatusOK)
		return
	}
	t.Fatalf("Соединение с API не установлено за %d попыток: %v", connectionAttempts, lastErr)
}
