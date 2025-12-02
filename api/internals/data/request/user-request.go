package request

type UpdateAddressRequest struct {
	AddressID   int64  `json:"address_id"`
	Name        string `json:"name"`
	PinCode     int32  `json:"pin_code"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	PhoneNumber int64  `json:"phone_number"`
}
