package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"
)

func GenerateSignature(timestamp string, secret string) string {
	data := fmt.Sprintf("%s%s", timestamp, secret)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func GetTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func ErrorMessage(err error) string {
	if err != nil {
		return err.Error()
	}
	return "Unknown error"
}

// HashCreate generates a signature hash from a request map and a key (appKey or tokenKey)
// Use appKey for Provider Endpoints (/games and /launch)
// Use tokenKey for validating callback requests from provider (/player_info, /bet, /win and /rollback)
func HashCreate(request map[string]interface{}, key string) string {
	// Step 1: Sort the keys in the request
	keys := getKeys(request)
	sort.Strings(keys)

	var hashkey string

	// Step 2: Iterate through sorted keys
	for _, k := range keys {
		value := request[k]

		switch v := value.(type) {
		case map[string]interface{}: // Step 3a: Handle nested maps
			nestedKeys := getKeys(v)
			sort.Strings(nestedKeys)

			for _, nestedKey := range nestedKeys {
				nestedValue := v[nestedKey]
				serialized, _ := json.Marshal(nestedValue)
				md5Hash := md5.Sum(serialized)
				hashkey += "&" + nestedKey + "=" + hex.EncodeToString(md5Hash[:])
			}

		case []interface{}: // Step 3b: Handle arrays
			for index, arrayValue := range v {
				serialized, _ := json.Marshal(arrayValue)
				md5Hash := md5.Sum(serialized)
				hashkey += "&" + strconv.Itoa(index) + "=" + hex.EncodeToString(md5Hash[:])
			}

		default: // Step 3c: Handle primitive types
			// Convert primitive value directly to string and append
			hashkey += "&" + k + "=" + fmt.Sprintf("%v", v)
		}
	}

	// Step 4: Trim leading "&"
	if len(hashkey) > 0 && hashkey[0] == '&' {
		hashkey = hashkey[1:]
	}

	// Step 5: Append the tokenKey/appKey to the concatenated string
	finalString := hashkey + key

	// Step 6: Compute final MD5 hash of the entire string
	finalHash := md5.Sum([]byte(finalString))
	signatureKey := hex.EncodeToString(finalHash[:])

	return signatureKey
}

// GenerateSignatureKey generates signature key as per Betkraft specification (legacy wrapper)
// Deprecated: Use HashCreate instead
func GenerateSignatureKey(request map[string]interface{}, appKey string) string {
	return HashCreate(request, appKey)
}

// getKeys extracts and returns all keys from a map
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
