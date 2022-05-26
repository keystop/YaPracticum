package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/keystop/YaPracticum.git/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RepoMock struct {
	mock.Mock
}

func (m *RepoMock) SaveURL(ctx context.Context, url, baseURL, userID string) (string, error) {
	args := m.Called(ctx, url, baseURL, userID)
	return args.String(0), args.Error(1)
}

func (m *RepoMock) SaveURLs(ctx context.Context, urls map[string]string, baseURL, userID string) (map[string]string, error) {
	args := m.Called(ctx, urls, baseURL, userID)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *RepoMock) GetURL(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

// func (m *RepoMock) Get() map[string][]string {
// 	return nil
// }

// func (m *RepoMock) ToSet() *map[string][]string {
// 	return nil
// }

func (m *RepoMock) FindUser(context.Context, string) bool {
	return false
}

func (m *RepoMock) CreateUser(context.Context) (string, error) {
	return "", nil
}

func (m *RepoMock) GetUserURLs(context.Context, string) ([]models.URLs, error) {
	return nil, nil
}

func (m *RepoMock) CheckDBConnection(context.Context) error {
	return nil
}

func (m *RepoMock) SetURLsToDel(context.Context, []string, string) error {
	return nil
}

type OptsMock struct {
	mock.Mock
}

func (o *OptsMock) ServAddr() string {
	args := o.Called()
	return args.String(0)
}

func (o *OptsMock) RespBaseURL() string {
	args := o.Called()
	return args.String(0)
}

func (o *OptsMock) RepoFileName() string {
	args := o.Called()
	return args.String(0)
}

func (o *OptsMock) DBConnString() string {
	args := o.Called()
	return args.String(0)
}

func newOptsMock() *OptsMock {
	optMock := new(OptsMock)
	optMock.On("ServAddr").Return("http://localhost:8080")
	optMock.On("RespBaseURL").Return("http://localhost")
	optMock.On("RepoFileName").Return("local.db")
	optMock.On("DBConnString").Return("user=kseikseich dbname=yap sslmode=disable")
	return optMock
}

var repoMock *RepoMock
var optsMock *OptsMock
var opt models.Options

func TestHandlerUrlGet(t *testing.T) {
	InitMocks()
	dataTests := map[string]map[string]interface{}{
		"test1": {
			"reqID":       "123123asdasd",
			"result":      "www.example.com",
			"resStatus":   http.StatusTemporaryRedirect,
			"mockReturn1": "www.example.com",
		},
		"test2": {
			"reqID":       "123123",
			"result":      "",
			"resStatus":   http.StatusBadRequest,
			"mockReturn1": "",
			"mockReturn2": errors.New("not found"),
		},
	}

	handler := http.HandlerFunc(HandlerURLGet)

	for key, value := range dataTests {

		ctx := context.WithValue(context.Background(), models.URLID, value["reqID"].(string))
		log.Println("start test", key)
		var err error
		if value["mockReturn2"] != nil {
			err = value["mockReturn2"].(error)
		}
		repoMock.On("GetURL", ctx, value["reqID"].(string)).Return(value["mockReturn1"].(string), err)

		r := httptest.NewRequest("GET", "/"+value["reqID"].(string), strings.NewReader(""))
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r.WithContext(ctx))

		res := w.Result()
		assert.Equal(t, value["resStatus"].(int), res.StatusCode, "Не верный код ответа GET")
		assert.Equal(t, w.Header().Get("Location"), value["result"].(string), "Не верный ответ GET")
		defer res.Body.Close()
	}
}

func TestHandlerUrlPost(t *testing.T) {
	ctx := context.WithValue(context.Background(), models.UserKey, "asdasd")
	repoMock.On("SaveURL", ctx, "www.example.com", opt.RespBaseURL()+"/", "asdasd").Return(opt.RespBaseURL()+"/123123asdasd", nil)

	handler := http.HandlerFunc(HandlerURLPost)
	r := httptest.NewRequest("POST", "http://localhost:8080", strings.NewReader("www.example.com"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r.WithContext(ctx))

	res := w.Result()
	b, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	assert.Equal(t, http.StatusCreated, res.StatusCode, "Не верный код ответа POST")
	assert.Equal(t, opt.RespBaseURL()+"/123123asdasd", string(b), "Не верный ответ POST")

}

func TestHandlerApiUrlPost(t *testing.T) {
	str := &struct {
		URL string
	}{
		URL: "www.example.com",
	}
	bOut, err := json.Marshal(str)
	if err != nil {
		t.Error("Ошибка серилизации")
	}
	ctx := context.WithValue(context.Background(), models.UserKey, "aasdasdSQW")
	repoMock.On("SaveURL", ctx, "www.example.com", opt.RespBaseURL()+"/", "aasdasdSQW").Return(opt.RespBaseURL()+"/123123asdasd", nil)
	handler := http.HandlerFunc(HandlerAPIURLPost)
	r := httptest.NewRequest("POST", "http://localhost:8080", bytes.NewBuffer(bOut))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r.WithContext(ctx))
	res := w.Result()
	b, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	assert.Equal(t, http.StatusCreated, res.StatusCode, "Не верный код ответа POST")
	assert.Equal(t, `{"result":"`+opt.RespBaseURL()+`/123123asdasd"}`, string(b), "Не верный ответ POST")

}

func InitMocks() {
	repoMock = new(RepoMock)
	optsMock = newOptsMock()
	opt = optsMock
	NewHandlers(repoMock, optsMock)
}
