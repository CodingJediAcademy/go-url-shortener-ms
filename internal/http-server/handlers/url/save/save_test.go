package save

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go-url-shortener-ms/internal/http-server/handlers/url/save/mocks"
	sl "go-url-shortener-ms/internal/lib/logger/handlers/slogdiscard"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
			// Тут поля respError и mockError оставляем пустыми,
			// т.к. это успешный запрос
		},
		// Другие кейсы ...
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Создаем объект мока стораджа
			urlSaverMock := mocks.NewURLSaver(t)

			// Если ожидается успешный ответ, значит к моку точно будет вызов
			// Либо даже если в ответе ожидаем ошибку,
			// но мок должен ответить с ошибкой, к нему тоже будет запрос:
			if tc.respError == "" || tc.mockError != nil {
				// Сообщаем моку, какой к нему будет запрос, и что надо вернуть
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once() // Запрос будет ровно один
			}

			// Создаем наш хэндлер
			handler := New(sl.NewDiscardLogger(), urlSaverMock)

			// Формируем тело запроса
			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			// Создаем объект запроса
			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			// Создаем ResponseRecorder для записи ответа хэндлера
			rr := httptest.NewRecorder()
			// Обрабатываем запрос, записывая ответ в рекордер
			handler.ServeHTTP(rr, req)

			// Проверяем, что статус ответа корректный
			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp Response

			// Анмаршаллим тело, и проверяем что при этом не возникло ошибок
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			// Проверяем наличие требуемой ошибки в ответе
			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
