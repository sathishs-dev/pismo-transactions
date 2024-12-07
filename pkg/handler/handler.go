package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/sathishs-dev/pismo-transactions/internal/meta/writer"
	"github.com/sathishs-dev/pismo-transactions/pkg/enums"
	"github.com/sathishs-dev/pismo-transactions/pkg/repository"
)

type handler struct {
	repo repository.PismoRepo
}

type Handler interface {
	CreateAccount() http.HandlerFunc
	GetAccount() http.HandlerFunc
	CreateTransaction() http.HandlerFunc
}

func NewHandler(repo repository.PismoRepo) Handler {
	return &handler{
		repo,
	}
}

// CreateAccount handler function handles account creation request
func (h *handler) CreateAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateAccountReqPayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg("failed to decode body")
			errorWriter(w, http.StatusBadRequest, "failed to decode body")
			return
		}

		if req.DocumentNumber == "" {
			errorWriter(w, http.StatusBadRequest, "document_number required")
			return
		}

		// check unique document_number
		isExists, err := h.repo.GetAccountByDocumentNo(r.Context(), req.DocumentNumber)
		if err != nil {
			log.Error().Err(err).Msg("failed to retrive the account")
			errorWriter(w, http.StatusInternalServerError, "please try again later.")
			return
		}

		if isExists {
			errorWriter(w, http.StatusConflict, "document_number already associated with an account.")
			return
		}

		// create account
		err = h.repo.CreateAccount(r.Context(), req.DocumentNumber)
		if err != nil {
			log.Error().Err(err).Msg("failed to store the account")
			errorWriter(w, http.StatusInternalServerError, "please try again later.")
			return
		}

		err = writer.WriteJSON(w, http.StatusAccepted, nil)
		if err != nil {
			log.Error().Err(err).Msg("failed to write")
			return
		}
	}
}

// GetAccount handler function handles fetch account requests
func (h *handler) GetAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accIdParam := chi.URLParam(r, "accountId")
		if accIdParam == "" {
			errorWriter(w, http.StatusBadRequest, "accountID required")
			return
		}

		accID, err := strconv.Atoi(accIdParam)
		if err != nil {
			log.Error().Err(err).Msg("failed to convert accountID")
			errorWriter(w, http.StatusInternalServerError, "please try again later.")
			return
		}

		if accID <= 0 {
			errorWriter(w, http.StatusBadRequest, "invalid accountId")
			return
		}

		account, err := h.repo.GetAccountByAccountID(r.Context(), accID)
		if err != nil {
			log.Error().Err(err).Msg("failed to retrieve data from store")
			errorWriter(w, http.StatusInternalServerError, "please try again later.")
			return
		}

		if account == nil {
			errorWriter(w, http.StatusBadRequest, "account not found")
			return
		}

		if err := writer.WriteJSON(w, http.StatusOK, GetAccountResPaylaod{
			AccountID:      account.AccountID,
			DocumentNumber: account.DocumentNo,
		}); err != nil {
			log.Error().Err(err).Msg("failed to write")
			return
		}
	}
}

// CreateTransaction handler function handles create txn requests
func (h *handler) CreateTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateTransactionReqPayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg("failed to decode body")
			errorWriter(w, http.StatusBadRequest, "failed to decode body")
			return
		}

		if errs := validateCreateTransactionReq(&req); len(errs) > 0 {
			errorWriter(w, http.StatusBadRequest, strings.Join(errs, "/"))
			return
		}

		operationType, err := enums.ParseOperationType(req.OperationTypeID)
		if err != nil {
			errorWriter(w, http.StatusBadRequest, "invalid operation_type_id")
			return
		}

		if req.Amount < 0 && !enums.AllowNegative(operationType) {
			errorWriter(w, http.StatusBadRequest, "negative transactions not allowed for the operation_type_id")
			return
		}

		if req.Amount > 0 && enums.AllowNegative(operationType) {
			errorWriter(w, http.StatusBadRequest, "positive transactions not allowed for the operation_type_id")
			return
		}

		// fetch account
		acc, err := h.repo.GetAccountByAccountID(r.Context(), req.AccountID)
		if err != nil {
			log.Error().Err(err).Msg("failed to retrieve the account")
			errorWriter(w, http.StatusInternalServerError, "please try again later.")
			return
		}

		if acc == nil {
			errorWriter(w, http.StatusBadRequest, "account not found")
			return
		}

		// create transaction
		err = h.repo.CreateTransaction(
			r.Context(),
			repository.Transaction{
				AccountID:       acc.AccountID,
				OperationTypeID: int(operationType),
				Amount:          req.Amount,
			},
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to store the transaction")
			errorWriter(w, http.StatusInternalServerError, "please try again later.")
			return
		}

		if err := writer.WriteJSON(w, http.StatusAccepted, nil); err != nil {
			log.Error().Err(err).Msg("failed to write")
			return
		}

	}
}

func validateCreateTransactionReq(req *CreateTransactionReqPayload) (errs []string) {
	if req.AccountID <= 0 {
		errs = append(errs, "invalid account_id")
	}
	if req.Amount == 0 {
		errs = append(errs, "invalid amount")
	}

	if req.OperationTypeID <= 0 {
		errs = append(errs, "invalid operation_type_id")
	}

	return
}

// errorWriter writes error response to the caller
func errorWriter(w http.ResponseWriter, status int, errMsg string) {
	errResp := GenericErrRespPayload{
		Message: errMsg,
	}

	if err := writer.WriteJSON(w, status, errResp); err != nil {
		log.Error().Err(err).Msg("failed writting to the client")
		return
	}
}
