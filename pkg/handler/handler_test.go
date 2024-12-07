package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sathishs-dev/pismo-transactions/pkg/mocks"
	"github.com/sathishs-dev/pismo-transactions/pkg/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type handlerTestSuite struct {
	suite.Suite
	recorder *httptest.ResponseRecorder
	router   *chi.Mux
	repo     *mocks.PismoRepo
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}

func (h *handlerTestSuite) SetupTest() {
	h.recorder = httptest.NewRecorder()
	h.router = chi.NewRouter()
	h.repo = new(mocks.PismoRepo)

	handler := NewHandler(h.repo)

	h.router.Post("/accounts", handler.CreateAccount())
	h.router.Get("/accounts/{accountId}", handler.GetAccount())
	h.router.Post("/transactions", handler.CreateTransaction())
}

func (h *handlerTestSuite) TestCreateAccount() {
	tcs := []struct {
		name               string
		reqBody            string
		expectedMocks      func(h *handlerTestSuite)
		expectedStatusCode int
	}{
		{
			name:    "Valid Create Account Request",
			reqBody: `{"document_number": "1234567890"}`,
			expectedMocks: func(h *handlerTestSuite) {
				h.repo.On("GetAccountByDocumentNo", mock.Anything, "1234567890").
					Return(false, nil)
				h.repo.On("CreateAccount", mock.Anything, "1234567890").
					Return(nil)
			},
			expectedStatusCode: http.StatusAccepted,
		},
		{
			name:               "Invalid Create Account Request - Empty Payload",
			reqBody:            ``,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid Create Account Request - Invalid Payload",
			reqBody:            `{"dc": "1234567890"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:    "Invalid Create Account Request - Existing Alreay Exists",
			reqBody: `{"document_number": "1234567890"}`,
			expectedMocks: func(h *handlerTestSuite) {
				h.repo.On("GetAccountByDocumentNo", mock.Anything, "1234567890").
					Return(true, nil)
			},
			expectedStatusCode: http.StatusConflict,
		},
		{
			name:               "Invalid Create Account Request - Get Account Fails",
			reqBody:            `{"document_number": "1234567890"}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedMocks: func(h *handlerTestSuite) {
				h.repo.On("GetAccountByDocumentNo", mock.Anything, "1234567890").
					Return(false, errors.New("err"))
			},
		},
		{
			name:               "Invalid Create Account Request - Store Account Fails",
			reqBody:            `{"document_number": "1234567890"}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedMocks: func(h *handlerTestSuite) {
				h.repo.On("GetAccountByDocumentNo", mock.Anything, "1234567890").
					Return(false, nil)
				h.repo.On("CreateAccount", mock.Anything, "1234567890").
					Return(errors.New("err"))
			},
		},
	}

	for _, tc := range tcs {
		h.T().Run(tc.name, func(t *testing.T) {
			h.recorder = httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(tc.reqBody))

			if tc.expectedMocks != nil {
				tc.expectedMocks(h)
			}

			h.router.ServeHTTP(h.recorder, req)
			h.Equal(tc.expectedStatusCode, h.recorder.Code)
			h.repo.ExpectedCalls = nil
		})
	}

}

func (h *handlerTestSuite) TestGetAccount() {
	tcs := []struct {
		name               string
		accID              int
		expectedMocks      func(h *handlerTestSuite)
		expectedStatusCode int
	}{
		{
			name:  "Valid Get Account Request",
			accID: 1,
			expectedMocks: func(h *handlerTestSuite) {
				h.repo.On("GetAccountByAccountID", mock.Anything, 1).
					Return(&repository.Account{
						AccountID:  1,
						DocumentNo: "1234567890",
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid Get Account Request - Empty AccountId",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid Get Account Request - No Account Found",
			accID:              100,
			expectedStatusCode: http.StatusBadRequest,
			expectedMocks: func(h *handlerTestSuite) {
				h.repo.On("GetAccountByAccountID", mock.Anything, 100).
					Return(nil, nil)
			},
		},
		{
			name:               "Invalid Get Account Request - Invalid Account ID",
			accID:              0,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Invalid Get Account Request - Fetching DataStore failed",
			accID: 100,
			expectedMocks: func(h *handlerTestSuite) {
				h.repo.On("GetAccountByAccountID", mock.Anything, 100).
					Return(&repository.Account{}, errors.New("failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tcs {
		h.T().Run(tc.name, func(t *testing.T) {
			h.recorder = httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%d", tc.accID), nil)

			if tc.expectedMocks != nil {
				tc.expectedMocks(h)
			}

			h.router.ServeHTTP(h.recorder, req)
			h.Equal(tc.expectedStatusCode, h.recorder.Code)
			h.repo.ExpectedCalls = nil
		})
	}

}

func (h *handlerTestSuite) TestCreateTransaction() {
	tcs := []struct {
		name               string
		reqBody            string
		expectedMocks      func(h *handlerTestSuite)
		expectedStatusCode int
	}{
		{
			name:    "Valid Create Transaction Request",
			reqBody: `{"account_id": 1, "operation_type_id": 1, "amount": -500.00}`,
			expectedMocks: func(h *handlerTestSuite) {
				h.repo.On("GetAccountByAccountID", mock.Anything, 1).
					Return(&repository.Account{
						AccountID:  1,
						DocumentNo: "1234567890",
					}, nil)
				h.repo.On("CreateTransaction", mock.Anything,
					repository.Transaction{AccountID: 1, OperationTypeID: 1, Amount: -500.00},
				).Return(nil)
			},
			expectedStatusCode: http.StatusAccepted,
		},
		{
			name:               "Invalid Create Transaction Request - Empty Payload",
			reqBody:            `{}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid Create Transaction Request - Invalid Account",
			reqBody:            `{"account_id": 2000, "operation_type_id": 1, "amount": -500.00}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedMocks: func(h *handlerTestSuite) {
				h.repo.On("GetAccountByAccountID", mock.Anything, 2000).
					Return(nil, nil)
			},
		},
		{
			name:               "Invalid Create Transaction Request - Invalid Operation Type ID",
			reqBody:            `{"account_id": 1, "operation_type_id": 5, "amount": -500.00}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid Create Transaction Request - invalid Amount",
			reqBody:            `{"account_id": 1, "operation_type_id": 1, "amount": 0.00}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid Create Transaction Request - Unsupported transaction for Operation Type",
			reqBody:            `{"account_id": 1, "operation_type_id": 1, "amount": 1000.00}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid Create Transaction Request - Unsupported transaction for Operation Type",
			reqBody:            `{"account_id": 1, "operation_type_id": 4, "amount": -1000.00}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid Create Transaction Request - Get Account Fails",
			reqBody:            `{"account_id": 1, "operation_type_id": 4, "amount": 1000.00}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedMocks: func(h *handlerTestSuite) {
				h.repo.On("GetAccountByAccountID", mock.Anything, 1).
					Return(nil, errors.New("err"))
			},
		},
		{
			name:               "Invalid Create Transaction Request - Create Transaction Fails",
			reqBody:            `{"account_id": 1, "operation_type_id": 4, "amount": 1000.00}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedMocks: func(h *handlerTestSuite) {
				h.repo.On("GetAccountByAccountID", mock.Anything, 1).
					Return(&repository.Account{
						AccountID:  1,
						DocumentNo: "1234567890",
					}, errors.New("err"))
				h.repo.On("CreateTransaction", mock.Anything,
					repository.Transaction{AccountID: 1, OperationTypeID: 4, Amount: 1000.00},
				).Return(errors.New("err"))
			},
		},
	}

	for _, tc := range tcs {
		h.T().Run(tc.name, func(t *testing.T) {
			h.recorder = httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(tc.reqBody))

			if tc.expectedMocks != nil {
				tc.expectedMocks(h)
			}

			h.router.ServeHTTP(h.recorder, req)
			h.Equal(tc.expectedStatusCode, h.recorder.Code)
			h.repo.ExpectedCalls = nil
		})
	}

}
