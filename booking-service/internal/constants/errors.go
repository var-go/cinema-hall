package constants

import "errors"

var ErrBookingNotFound = errors.New("booking not found")
var ErrBookingFailedUpdate = errors.New("failed update booking")
var ErrBookingAlreadyCancelled = errors.New("booking already cancelled")
var ErrBookingExpired = errors.New("the reservation time has expired")
var ErrInvalidID = errors.New("invalid id")
var ErrBookingAlreadyConfirmed = errors.New("booking already confirmed")
var ErrInvalidBookingStatus = errors.New("invalid booking status")
