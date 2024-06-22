package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	db "github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models/db"
	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type couponConfig struct {
	CouponProbability int

	ReserveTimeMin time.Time
	ReserveTimeMax time.Time
	SnatchTimeMin  time.Time
	SnatchTimeMax  time.Time
}

func (c *couponConfig) isBetweenTime(n time.Time, min time.Time, max time.Time) bool {
	h, m, s := n.Clock()
	clockStr := fmt.Sprintf("0001-01-01T%d:%d:%dZ", h, m, s)
	clock, err := time.Parse(time.RFC3339, clockStr)
	if err != nil {
		return false
	}
	return min.Before(clock) && clock.Before(max)
}

func (c *couponConfig) isReserveTime(n time.Time) bool {
	return c.isBetweenTime(n, c.ReserveTimeMin, c.ReserveTimeMax)
}
func (c *couponConfig) isSnatchTime(n time.Time) bool {
	return c.isBetweenTime(n, c.SnatchTimeMin, c.SnatchTimeMax)
}

func (c *couponConfig) isWinCoupon() bool {
	return rand.Intn(10000) < c.CouponProbability
}

// We can get coupon config from DB, but in there we simplify the problem, hardcode in there.
func getCouponConfig() couponConfig {
	ReserveTimeMin, _ := time.Parse(time.RFC3339, "0001-01-01T22:55:00Z")
	ReserveTimeMax, _ := time.Parse(time.RFC3339, "0001-01-01T22:59:00Z")
	SnatchTimeMin, _ := time.Parse(time.RFC3339, "0001-01-01T23:00:00Z")
	SnatchTimeMax, _ := time.Parse(time.RFC3339, "0001-01-01T23:01:00Z")
	return couponConfig{
		CouponProbability: 2000,

		ReserveTimeMin: ReserveTimeMin,
		ReserveTimeMax: ReserveTimeMax,
		SnatchTimeMin:  SnatchTimeMin,
		SnatchTimeMax:  SnatchTimeMax,
	}
}

// Reserve a coupon
func Reserve(userID string) (*db.Coupon, error) {

	config := getCouponConfig()
	t := time.Now()

	if !config.isReserveTime(t) {
		return nil, fmt.Errorf("not reserve time")
	}

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
	config := getCouponConfig()
	t := time.Now()

	if !config.isSnatchTime(t) {
		return nil, fmt.Errorf("not snatch time")
	}

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
