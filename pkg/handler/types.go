package handler

type (
	CreateAccountReqPayload struct {
		DocumentNumber string `json:"document_number"`
	}

	GetAccountResPaylaod struct {
		AccountID      int    `json:"account_id"`
		DocumentNumber string `json:"document_number"`
	}

	CreateTransactionReqPayload struct {
		AccountID       int     `json:"account_id"`
		OperationTypeID int     `json:"operation_type_id"`
		Amount          float64 `json:"amount"`
	}

	GenericErrRespPayload struct {
		Message string `json:"message"`
	}
)
