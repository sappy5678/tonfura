package services

import (
	"errors"

	db "github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models/db"
	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
)

// CreateUser create a user record
func Reserve(userID string) (*db.Coupon, error) {
	couponID := uuid.NewString()
	coupon := db.NewCoupon(couponID, userID, db.CouponStatusNotActive)
	err := mgm.Coll(coupon).Create(coupon)
	if err != nil {
		return nil, errors.New("cannot create new user")
	}

	return coupon, nil
}
