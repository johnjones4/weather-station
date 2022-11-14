package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func jsonResponse(w http.ResponseWriter, status int, info any) {
	bytes, err := json.Marshal(info)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	w.Write(bytes)
}

func errorResponse(w http.ResponseWriter, status int, err error) {
	msg := map[string]any{
		"ok":      false,
		"message": err.Error(),
	}
	jsonResponse(w, status, msg)
}

func readJson(r *http.Request, readTo any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, readTo)
	if err != nil {
		return err
	}

	return nil
}

func readDateParameter(r *http.Request, param string, defaultValue time.Time) time.Time {
	str := r.URL.Query().Get(param)
	if str == "" {
		return defaultValue
	}

	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return time.Time{}
	}

	return t
}
