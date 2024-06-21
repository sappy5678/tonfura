package services

import (
	"context"
	"fmt"
	"time"

	db "github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models/db"
	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

// Reserve a coupon
func Reserve(userID string) (*db.Coupon, error) {
	couponID := uuid.NewString()
	coupon := db.NewCoupon(couponID, userID, db.CouponStatusNotActive)
	isReserved, err := checkAndSetReserveUser(userID)
	if err != nil {
		return nil, err
	} else if isReserved {
		err := mgm.Coll(coupon).First(bson.M{"userID": userID}, coupon)
		return coupon, err
	}
	err = mgm.Coll(coupon).Create(coupon)
	if err != nil {
		return nil, err
	}

	return coupon, nil
}

// Snatch activate a coupon
func Snatch(userID string) (*db.Coupon, error) {
	coupon := &db.Coupon{}
	res := mgm.Coll(coupon).FindOneAndUpdate(context.Background(), bson.M{"userID": userID}, bson.M{"$set": bson.M{"status": db.CouponStatusActive}})
	if res.Err() != nil {
		return nil, res.Err()
	}

	err := res.Decode(coupon)
	if err != nil {
		return nil, err
	}

	return coupon, nil
}

func getReserveKey(userID string) string {
	return "coupon:bf:" + userID
}

func checkAndSetReserveUser(userID string) (bool, error) {
	if !Config.UseRedis {
		return false, fmt.Errorf("redis cannot used")
	}

	key := getReserveKey(userID)

	ret := GetRedisDefaultClient().SetNX(context.TODO(), key, "1", 20*time.Minute)

	return !ret.Val(), ret.Err()
}
