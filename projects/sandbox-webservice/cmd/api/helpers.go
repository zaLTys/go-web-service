package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type envelope map[string]any

// -----------------------------------------------------------------------------
// JSON Response Writer
// -----------------------------------------------------------------------------

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(js)
	return err
}

// -----------------------------------------------------------------------------
// JSON Reader With Good Error Messages
// -----------------------------------------------------------------------------

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Malformed JSON
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON at character %d", syntaxError.Offset)

		// Unknown field
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", field)

		// Wrong type
		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)

		// Empty body
		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")

		// Request too large
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		default:
			return err
		}
	}

	// Prevent multiple JSON objects
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain only a single JSON object")
	}

	return nil
}

// -----------------------------------------------------------------------------
// Unified Error Responses
// -----------------------------------------------------------------------------

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	// TODO: app.logger.Error(err)  // recommended
	app.errorResponse(w, nil, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (app *application) badRequest(w http.ResponseWriter, err error) {
	app.errorResponse(w, nil, http.StatusBadRequest, err.Error())
}

func (app *application) notFound(w http.ResponseWriter) {
	app.errorResponse(w, nil, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}
