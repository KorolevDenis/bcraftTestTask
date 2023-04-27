package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func Message(status bool, message string, name string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message, "name": name}
}

func Respond(w http.ResponseWriter, data map[string]interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}

// GetIdFromUrl function that gets the identifier from the URL according to the given pattern
func GetIdFromUrl(url, pattern string) (uint, error) {
	re := regexp.MustCompile(pattern + "/([1-9]{1}[0-9]*)")
	str := re.FindString(url)
	if len(str) == 0 {
		return 0, fmt.Errorf("id must be natural number")
	}

	id, _ := strconv.ParseUint(strings.Split(str, "/")[1], 10, 64)

	return uint(id), nil
}
