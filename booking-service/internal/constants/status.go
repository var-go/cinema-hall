package constants

type BookingStatus string

const (
	Pending   BookingStatus = "pending"
	Confirmed BookingStatus = "confirmed"
	Cancelled BookingStatus = "cancelled"
	Expired   BookingStatus = "expired"
)

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "pending"
	PaymentPaid    PaymentStatus = "paid"
)

const (
	BookingTimeoutMinutes = 15
)
