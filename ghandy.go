package ghandy

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

// AsString gets a parameter coming from a http request as string, truncated to maxLength
// Only maxLength >= 1 is considered. Otherwise, it's ignored
func AsString(r *http.Request, key string, maxLength int) string {
	if err := r.ParseForm(); err != nil {
		log.Println(r.RequestURI, err)
		return ""
	}

	s := r.FormValue(key)

	if s == "" {
		s = r.URL.Query().Get(key)
	}

	if s == "" {
		params := mux.Vars(r)
		s = params[key]
	}

	if s != "" && (maxLength > 0) && utf8.RuneCountInString(s) >= maxLength {
		return s[0:maxLength]
	}

	return s
}

// AsInt gets a parameter coming from a http request as an integer
// It tries to guess if it's a signed/negative integer
func AsInt(r *http.Request, key string) int {
	if err := r.ParseForm(); err != nil {
		log.Println(r.RequestURI, err)
		return 0
	}

	s := r.FormValue(key)

	if s == "" {
		s = r.URL.Query().Get(key)
	}

	if s == "" {
		params := mux.Vars(r)
		s = params[key]
	}

	if s == "" {
		return 0
	}

	neg := s[0:1] == "-"

	i, _ := strconv.Atoi(s)

	if neg && (i > 0) {
		return i * -1
	}

	return i
}

// AsFloat gets a parameter coming from a http request as float64 number
// You have to inform the decimal separator symbol.
// If decimalSeparator is period, engine considers thousandSeparator is comma, and vice-versa.
func AsFloat(r *http.Request, key string, decimalSeparator rune) float64 {
	if err := r.ParseForm(); err != nil {
		log.Println(r.RequestURI, err)
		return 0
	}

	s := r.FormValue(key)

	if s == "" {
		s = r.URL.Query().Get(key)
	}

	if s == "" {
		params := mux.Vars(r)
		s = params[key]
	}

	if s == "" {
		return 0
	}

	thousandSeparator := ','

	if decimalSeparator == ',' {
		thousandSeparator = '.'
	}

	s = strings.ReplaceAll(s, string(thousandSeparator), "")

	f, _ := strconv.ParseFloat(s, 64)

	return f
}

// JSONAsStruct decode json to a given anatomically compatible struct
func JSONAsStruct(r *http.Request, targetStruct interface{}, closeBody bool) error {
	err := json.NewDecoder(r.Body).Decode(targetStruct)

	if closeBody {
		defer func() {
			if err0 := r.Body.Close(); err0 != nil {
				log.Println(err0)
			}
		}()
	}

	return err
}
