package models

import (
	"github.com/kamva/mgm/v3"
)

type CouponStatus int

const (
	CouponStatusCancel    CouponStatus = -1
	CouponStatusNotActive CouponStatus = 0
	CouponStatusActive    CouponStatus = 1
	CouponStatusUsed      CouponStatus = 2
)

type Coupon struct {
	mgm.DefaultModel `bson:",inline"`
	CouponID         string       `json:"couponID" bson:"couponID"`
	UserID           string       `json:"userID" bson:"userID"`
	Status           CouponStatus `json:"status" bson:"status"`
}

func NewCoupon(couponID string, userID string, status CouponStatus) *Coupon {
	return &Coupon{
		CouponID: couponID,
		UserID:   userID,
		Status:   status,
	}
}

func (model *Coupon) CollectionName() string {
	return "coupon"
}

// You can override Collection functions or CRUD hooks
// https://github.com/Kamva/mgm#a-models-hooks
// https://github.com/Kamva/mgm#collections
