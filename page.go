package petra



type PageService service


type Page struct {
	ID           int                 `json:"id,omitempty"`
	CreatedAt    string              `json:"createdAt,omitempty"`
	UpdatedAt    string              `json:"updatedAt,omitempty"`
	Domain       string              `json:"domain,omitempty"`
	Integration  int                 `json:"integration,omitempty"`
	Name         string              `json:"name,omitempty"`
	Slug         string              `json:"slug,omitempty"`
	Description  string              `json:"description,omitempty"`
	Amount       float32             `json:"amount,omitempty"`
	Currency     string              `json:"currency,omitempty"`
	Active       bool                `json:"active,omitempty"`
	RedirectURL  string              `json:"redirect_url,omitempty"`
	CustomFields []map[string]string `json:"custom_fields,omitempty"`
}
