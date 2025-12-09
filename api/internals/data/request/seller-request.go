package request

type SellerKYC struct {
	BusinessName      string           `json:"business_name" binding:"required"`
	BusinessAddress   string           `json:"business_address" binding:"required"`
	Website           string           `json:"website" binding:"required,url"`
	ContactNumber     int64            `json:"contact_number" binding:"required"`
	ContactPerson     string           `json:"contact_person" binding:"required"`
	GstNumber         string           `json:"gst_number" binding:"required"`
	BankAccountNumber string           `json:"bank_account_number" binding:"required"`
	IfscCode          string           `json:"bank_ifsc_code" binding:"required"`
	Documents         []SellerDocument `json:"documents" binding:"required,dive"`
}
type SellerDocument struct {
	Document    string `json:"document" binding:"required"`
	DocumentURL string `json:"document_url" binding:"required,url"`
}
type RegisterSellerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Username string `json:"username" binding:"required"`
}
