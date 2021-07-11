package petra

import "fmt"

type CustomerService service


type Customer struct {
	ID             int            `json:"id,omitempty"`
	CreatedAt      string         `json:"createdAt,omitempty"`
	UpdatedAt      string         `json:"updatedAt,omitempty"`
	Domain         string         `json:"domain,omitempty"`
	Integration    int            `json:"integration,omitempty"`
	FirstName      string         `json:"first_name,omitempty"`
	LastName       string         `json:"last_name,omitempty"`
	Email          string         `json:"email,omitempty"`
	Phone          string         `json:"phone,omitempty"`
	Metadata       Metadata       `json:"metadata,omitempty"`
	CustomerCode   string         `json:"customer_code,omitempty"`
}

type CustomerList struct {
	Meta   ListMeta
	Values []Customer `json:"data"`
}


func (s *CustomerService) Create(customer *Customer) (*Customer, error) {
	u := fmt.Sprintf("/customer")
	cust := &Customer{}
	err := s.client.Call("POST", u, customer, cust)

	return cust, err
}

func (s *CustomerService) Update(customer *Customer) (*Customer, error) {
	u := fmt.Sprintf("customer/%d", customer.ID)
	cust := &Customer{}
	err := s.client.Call("PUT", u, customer, cust)

	return cust, err
}


func (s *CustomerService) Get(customerCode string) (*Customer, error) {
	u := fmt.Sprintf("/customer/%s", customerCode)
	cust := &Customer{}
	err := s.client.Call("GET", u, nil, cust)

	return cust, err
}

func (s *CustomerService) List() (*CustomerList, error) {
	return s.ListN(10, 0)
}


func (s *CustomerService) ListN(count, offset int) (*CustomerList, error) {
	u := paginateURL("/customer", count, offset)
	cust := &CustomerList{}
	err := s.client.Call("GET", u, nil, cust)
	return cust, err
}