package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danieeelfr/swd-challenge/internal/config"
	response "github.com/danieeelfr/swd-challenge/internal/httpserver/api/apimodels/response"
	"github.com/danieeelfr/swd-challenge/internal/httpserver/utils/requestvalidator"
	"github.com/danieeelfr/swd-challenge/internal/models"
	"github.com/danieeelfr/swd-challenge/pkg/wait"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type LoginHandlerTestSuite struct {
	suite.Suite
	cfg *config.Config
	e   *echo.Echo
	h   *LoginHandler
	rec *httptest.ResponseRecorder
}

func (s *LoginHandlerTestSuite) SetupSuite() {
	s.e = echo.New()
	s.cfg = &config.Config{
		HTTPServerConfig: &config.HTTPServerConfig{
			HTTPServerHost: "8060",
			WaitToShutdown: 2,
		},
		// without database mock it is an integration tests
		MySQLRepositoryConfig: &config.MySQLRepositoryConfig{
			DBUser:     "dev",
			DBPassword: "dev",
			DBName:     "dev",
			DBHost:     "0.0.0.0",
			DBPort:     "3306",
		},
	}

	h, err := NewLoginHandler(s.cfg, s.e, wait.New())
	s.NoError(err)
	s.h = h
}

func (s *LoginHandlerTestSuite) TestAuthorizeWithSuccess() {

	type testCase struct {
		description        string
		loginModel         models.Login
		expectedStatusCode int
		expectedResponse   response.LoginResponse
	}

	testCases := []testCase{
		{
			description: "success-manager",
			loginModel: models.Login{
				User:     "manager1",
				Password: "12345",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: response.LoginResponse{
				Message: "authorized with success!",
				Token:   "random-token",
			},
		},
		{
			description: "success-technician",
			loginModel: models.Login{
				User:     "technician1",
				Password: "12345",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: response.LoginResponse{
				Message: "authorized with success!",
				Token:   "random-token",
			},
		},
	}

	for _, tc := range testCases {

		s.Run(tc.description, func() {
			payload, err := json.Marshal(tc.loginModel)
			s.NoError(err)

			req := httptest.NewRequest(
				http.MethodPost,
				"/login/authorize",
				strings.NewReader(string(payload)),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			ctx := s.e.NewContext(req, s.rec)
			ctx.Echo().Validator = &requestvalidator.CustomValidator{Validator: validator.New()}

			h, err := NewLoginHandler(s.cfg, s.e, wait.New())
			s.NoError(err)

			if s.NoError(h.Authorize(ctx)) {
				s.Equal(tc.expectedStatusCode, s.rec.Code)

				resp, err := ioutil.ReadAll(s.rec.Body)
				s.NoError(err)

				var r response.LoginResponse

				s.NoError(json.Unmarshal(resp, &r))

				s.Equal(tc.expectedResponse.Message, r.Message)
				s.NotEmpty(r.Token)

			}
		})
	}
}

func (s *LoginHandlerTestSuite) TestAuthorizeWithErrors() {

	type testCase struct {
		description        string
		loginModel         models.Login
		expectedStatusCode int
	}

	testCases := []testCase{
		{
			description: "badrequest-missing-user",
			loginModel: models.Login{
				User:     "",
				Password: "12345",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "badrequest-missing-password",
			loginModel: models.Login{
				User:     "technician1",
				Password: "",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "not-found",
			loginModel: models.Login{
				User:     "technician123",
				Password: "123",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {

		s.Run(tc.description, func() {
			s.rec = httptest.NewRecorder()

			payload, err := json.Marshal(tc.loginModel)
			s.NoError(err)

			req := httptest.NewRequest(
				http.MethodPost,
				"/login/authorize",
				strings.NewReader(string(payload)),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			ctx := s.e.NewContext(req, s.rec)
			ctx.Echo().Validator = &requestvalidator.CustomValidator{Validator: validator.New()}

			h, err := NewLoginHandler(s.cfg, s.e, wait.New())
			s.NoError(err)

			err = h.Authorize(ctx)
			he, ok := err.(*echo.HTTPError)
			if ok {
				s.Equal(tc.expectedStatusCode, he.Code)
			}
		})
	}
}

func TestLoginHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(LoginHandlerTestSuite))
}
