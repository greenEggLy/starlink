package utils

import (
	"context"
	pb "starlink/pb"
	"strconv"
	"sync"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Redis struct {
	client       *redis.Client
	expire_sec   int
	maxTimeStamp int64
	m            sync.Mutex
}

func NewRedis(expire_seconds ...int) *Redis {
	if len(expire_seconds) == 0 {
		expire_seconds = append(expire_seconds, 60)
	}
	r := &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		expire_sec:   expire_seconds[0],
		maxTimeStamp: 0,
	}
	return r
}

func (r *Redis) SetPosition(value *pb.PositionInfo) {
	// log.Printf("[redis] set position, timestamp: %s", value.Timestamp)
	ts, err1 := strconv.ParseInt(value.Timestamp, 10, 64)
	if err1 != nil {
		panic("parse time error")
	}
	timestamp_time := time.Unix(ts, 0)
	expire_time := timestamp_time.Add(time.Duration(r.expire_sec) * time.Second)
	if ts > r.maxTimeStamp {
		r.maxTimeStamp = ts
	}
	key := r.appendString(value.TargetName, ts)
	r.m.Lock()
	r.client.HSet(context.Background(), key, "ALT", value.Alt, "LAT", value.Lat, "LNG", value.Lng)
	r.client.ExpireAt(context.Background(), key, expire_time)
	r.m.Unlock()
	// log.Printf("[redis] set position")

}

// get all positions by target name and timestamp
func (r *Redis) GetAllPos(targetNames []*string) []*pb.PositionInfo {
	var pos []*pb.PositionInfo
	time_now := time.Now().Unix()
	r.m.Lock()
	for _, name := range targetNames {
		for ts := time_now; ts <= r.maxTimeStamp; ts++ {
			// key in format: "targetName-timestamp"
			key := r.appendString(*name, ts)

			value, err := r.client.HGetAll(context.Background(), key).Result()
			if err != nil {
				continue
			}
			if value["ALT"] == "" || value["LAT"] == "" || value["LNG"] == "" {
				continue
			}
			alt, err1 := strconv.ParseFloat(value["ALT"], 32)
			lat, err2 := strconv.ParseFloat(value["LAT"], 32)
			lng, err3 := strconv.ParseFloat(value["LNG"], 32)
			if err1 != nil || err2 != nil || err3 != nil {
				panic("cannot parse position info")
			}
			p := pb.PositionInfo{
				Timestamp:  strconv.FormatInt(ts, 10),
				Alt:        float32(alt),
				Lat:        float32(lat),
				Lng:        float32(lng),
				TargetName: *name,
			}
			pos = append(pos, &p)
			// log.Printf("[redis] get position info: ALT: %s, LAT: %s, LNG: %s", value["ALT"], value["LAT"], value["LNG"])

		}
	}
	r.m.Unlock()
	return pos
}

func (r *Redis) appendString(targetName string, ts int64) string {
	buf := make([]byte, 0)
	buf = append(buf, []byte(targetName)...)
	buf = append(buf, "-"...)
	buf = append(buf, []byte(strconv.FormatInt(ts, 10))...)
	return string(buf)
}
