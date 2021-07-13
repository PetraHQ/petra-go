package petra

import "fmt"

type InvoiceService service

type Invoice struct {
	ID int `json:"id:,omitempty"`
	CreatedAt      string         `json:"createdAt,omitempty"`
	UpdatedAt      string         `json:"updatedAt,omitempty"`
	Domain         string         `json:"domain,omitempty"`
	Integration    int            `json:"integration,omitempty"`
	Email          string         `json:"email,omitempty"`
}

type InvoiceList struct {
	Meta   ListMeta
	Values []Invoice `json:"data"`
}

func (s *InvoiceService)Create(invoice *Invoice)(*Invoice,error)  {
	u := fmt.Sprintf("/invoice")
	inv := &Invoice{}
	err := s.client.Call("POST", u, invoice, inv)

	return inv, err
}


func (s *InvoiceService) Update(customer *Invoice) (*Invoice, error) {
	u := fmt.Sprintf("/invoice/%d", customer.ID)
	inv := &Invoice{}
	err := s.client.Call("PUT", u, customer, inv)

	return inv, err
}


func (s *InvoiceService) Get(id string) (*Invoice, error) {
	u := fmt.Sprintf("/invoice/%s", id)
	inv := &Invoice{}
	err := s.client.Call("GET", u, nil, inv)

	return inv, err
}

func (s *InvoiceService) List() (*InvoiceList, error) {
	return s.ListN(10, 0)
}


func (s *InvoiceService) ListN(count, offset int) (*InvoiceList, error) {
	u := paginateURL("/invoice", count, offset)
	inv := &InvoiceList{}
	err := s.client.Call("GET", u, nil, inv)
	return inv, err
}