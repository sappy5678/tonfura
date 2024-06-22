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
	"go.mongodb.org/mongo-driver/mongo/options"
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
	clockStr := fmt.Sprintf("0001-01-01T%02d:%02d:%02dZ", h, m, s)
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
	ReserveTimeMin, _ := time.Parse(time.RFC3339, "0001-01-01T00:55:00Z")
	ReserveTimeMax, _ := time.Parse(time.RFC3339, "0001-01-01T22:59:00Z")
	SnatchTimeMin, _ := time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")
	SnatchTimeMax, _ := time.Parse(time.RFC3339, "0001-01-01T23:01:00Z")
	return couponConfig{
		CouponProbability: 10000,

		ReserveTimeMin: ReserveTimeMin,
		ReserveTimeMax: ReserveTimeMax,
		SnatchTimeMin:  SnatchTimeMin,
		SnatchTimeMax:  SnatchTimeMax,
	}
}

// Reserve a coupon
func Reserve(userID string) error {

	config := getCouponConfig()
	t := time.Now()

	if !config.isReserveTime(t) {
		return fmt.Errorf("not reserve time")
	}

	couponID := uuid.NewString()
	coupon := db.NewCoupon(couponID, userID, db.CouponStatusNotActive)
	isWin := config.isWinCoupon()
	isReserved, err := checkAndSetReserveUser(userID, couponID)
	if err != nil {
		return err
	} else if isReserved {
		return nil
	} else if !isWin {
		return nil
	}

	err = mgm.Coll(coupon).Create(coupon)
	if err != nil {
		removeReserveUserRecord(userID) // if mongo write fail, give users the opportunity to retry
		return err
	}

	return nil
}

// Snatch activate a coupon
func Snatch(userID string) (*db.Coupon, error) {
	config := getCouponConfig()
	t := time.Now()

	if !config.isSnatchTime(t) {
		return nil, fmt.Errorf("not snatch time")
	}
	couponID, err := checkReserveUserCouponID(userID)
	if err != nil {
		return nil, err
	} else if couponID == "" {
		return nil, nil
	}

	coupon := &db.Coupon{}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := mgm.Coll(coupon).FindOneAndUpdate(context.Background(),
		bson.M{"couponID": couponID, "status": db.CouponStatusNotActive},
		bson.M{"$set": bson.M{"status": db.CouponStatusActive}, "$currentDate": bson.M{"updated_at": true}},
		opt,
	)

	if res.Err() != nil {
		// if coupon already be used or activated, get coupon info and return
		err := mgm.Coll(coupon).First(bson.M{"couponID": couponID}, coupon)
		if err == nil {
			return coupon, nil
		}
		return nil, res.Err()
	}

	err = res.Decode(coupon)
	if err != nil {
		return nil, err
	}

	return coupon, nil
}

func getReserveKey(userID string) string {
	return "coupon:" + userID
}

func removeReserveUserRecord(userID string) error {
	if !Config.UseRedis {
		return fmt.Errorf("redis cannot used")
	}

	key := getReserveKey(userID)

	ret := GetRedisDefaultClient().Del(context.Background(), key)

	return ret.Err()
}

func checkAndSetReserveUser(userID string, couponID string) (bool, error) {
	if !Config.UseRedis {
		return false, fmt.Errorf("redis cannot used")
	}

	key := getReserveKey(userID)

	ret := GetRedisDefaultClient().SetNX(context.Background(), key, couponID, 20*time.Minute)

	return !ret.Val(), ret.Err()
}

func checkReserveUserCouponID(userID string) (string, error) {
	if !Config.UseRedis {
		return "", fmt.Errorf("redis cannot used")
	}

	key := getReserveKey(userID)

	ret := GetRedisDefaultClient().Get(context.Background(), key)
	if ret.Err() != nil {
		return "", ret.Err()
	}
	val := ret.Val()

	return val, nil
}
