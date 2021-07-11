package petra



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
