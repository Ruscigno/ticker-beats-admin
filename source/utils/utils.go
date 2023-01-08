package utils

import (
	"crypto/sha256"
	"fmt"
	"hash/fnv"
	"regexp"
	"strconv"
	"time"
)

const (
	ServerPort = "31015"
	//GlobalDBVersion ...
	GlobalDBVersion            int = 1
	ConfigDatabaseDSN              = "CONFIG_DATABASE_DSN"
	ConfigDbMaxOpenConnections     = "CONFIG_DB_MAX_OPEN_CONNECTIONS"
	ConfigDbIdleConnections        = "CONFIG_DB_IDLE_CONNECTIONS"
	ConfigServerPort               = "CONFIG_SERVER_PORT"
)

type Database struct {
	DSN                string
	MaxOpenConnections string
	MaxIdleConnections string
}

type Config struct {
	Database
}

var DefaultFieldsToIgnore = []string{"created", "updated", "deleted"}

var MT5_MIN_DATE = time.Now().UTC().Add(time.Hour * 24 * 18000 * -1).Unix()

// AsSha256 hash function
func AsSha256(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))

	return fmt.Sprintf("%x", h.Sum(nil))
}

// AsSha256OfOrderedMap hash function
// My tests showed this is unnecessary because apparently, the order of the key is the same
// if the content of the row is the same. So, two rows could have different field order, but
// it should be always the same for the specific row
func AsSha256OfOrderedMap(o map[string]interface{}, orderedKey []string) string {
	type ordered struct {
		key   string
		value interface{}
	}
	var list []ordered
	for _, j := range orderedKey {
		a := o[j]
		list = append(list, ordered{
			key:   j,
			value: a,
		})
	}
	return AsSha256(list)
}

// CleanAndTrim just alphanumeric chars
func CleanAndTrim(text string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(text, "")
}

func NumericHash(ac interface{}) int64 {
	// :P
	rash := AsSha256(ac)
	h := fnv.New32a()
	h.Write([]byte(rash))
	//TODO: is there a better way to convert uint32 into int64
	str := fmt.Sprint(h.Sum32())
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}
