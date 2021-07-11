package petra

import "fmt"

type TransactionService service

type TransactionList struct {
	Meta   ListMeta
	Values []Transaction `json:"data"`
}

type Transaction struct {
	ID              int                    `json:"id,omitempty"`
	CreatedAt       string                 `json:"createdAt,omitempty"`
	Domain          string                 `json:"domain,omitempty"`
	Metadata        string                 `json:"metadata,omitempty"` //TODO: why is transaction metadata a string?
	Status          string                 `json:"status,omitempty"`
	Reference       string                 `json:"reference,omitempty"`
	Amount          float32                `json:"amount,omitempty"`
	Message         string                 `json:"message,omitempty"`
	GatewayResponse string                 `json:"gateway_response,omitempty"`
	PaidAt          string                 `json:"piad_at,omitempty"`
	Channel         string                 `json:"channel,omitempty"`
	Currency        string                 `json:"currency,omitempty"`
	IPAddress       string                 `json:"ip_address,omitempty"`
	Log             map[string]interface{} `json:"log,omitempty"` // TODO: same as timeline?
	Fees            int                    `json:"int,omitempty"`
	FeesSplit       string                 `json:"fees_split,omitempty"` // TODO: confirm data type
	Customer        Customer               `json:"customer,omitempty"`
}

type CreateTransactionRequest struct {
	Amount            float32  `json:"amount,omitempty"`
	Email             string   `json:"email,omitempty"`
	Bearer            string   `json:"bearer,omitempty"`
}


func (s *TransactionService) Initialize(txn *CreateTransactionRequest) (Response, error) {
	u := fmt.Sprintf("/transaction")
	resp := Response{}
	err := s.client.Call("POST", u, txn, &resp)
	return resp, err
}


func (s *TransactionService) List() (*TransactionList, error) {
	return s.ListN(10, 1)
}

func (s *TransactionService) ListN(count, offset int) (*TransactionList, error) {
	u := paginateURL("/transaction", count, offset)
	txns := &TransactionList{}
	err := s.client.Call("GET", u, nil, txns)
	return txns, err
}