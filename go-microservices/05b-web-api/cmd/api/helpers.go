package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// Retrieve the "id" URL parameter from the current request context.
func (app *application) readIDParam(r *http.Request) (int64, error) {
    params := httprouter.ParamsFromContext(r.Context())

    id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
    if err != nil || id < 1 {
        return 0, errors.New("invalid id parameter")
    }

    return id, nil
}

func (app *application) readUuidParam(r *http.Request) (uuid.UUID, error) {
    params := httprouter.ParamsFromContext(r.Context())

    id, err := uuid.Parse(params.ByName("id"))
    if err != nil {
        return uuid.Nil, errors.New("invalid id parameter")
    }

    return id, nil
}

// Define a writeJSON() helper for sending responses.
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
    js, err := json.Marshal(data)
    if err != nil {
        return err
    }

    // Append a newline to make it easier to view in terminal applications.
    js = append(js, '\n')

    // Writer headers
    for key, value := range headers {
        w.Header()[key] = value
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(js)

    return nil
}

// newUUID returns a new UUID and ignores potential errors
func (app *application) newUUID() uuid.UUID {
	r, _ := uuid.NewUUID()
	return r
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
    // Decode the request body into the target destination. 
    err := json.NewDecoder(r.Body).Decode(dst)
    if err != nil {
        // If there is an error during decoding, start the triage...
        var syntaxError *json.SyntaxError
        var unmarshalTypeError *json.UnmarshalTypeError
        var invalidUnmarshalError *json.InvalidUnmarshalError

        switch {
        // Use the errors.As() function to check whether the error has the type 
        // *json.SyntaxError. If it does, then return a plain-english error message 
        // which includes the location of the problem.
        case errors.As(err, &syntaxError):
            return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

        // In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
        // for syntax errors in the JSON. So we check for this using errors.Is() and
        // return a generic error message. There is an open issue regarding this at
        // https://github.com/golang/go/issues/25956.
        case errors.Is(err, io.ErrUnexpectedEOF):
            return errors.New("body contains badly-formed JSON")

        // Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
        // JSON value is the wrong type for the target destination. If the error relates
        // to a specific field, then we include that in our error message to make it 
        // easier for the client to debug.
        case errors.As(err, &unmarshalTypeError):
            if unmarshalTypeError.Field != "" {
                return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
            }
            return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

        // An io.EOF error will be returned by Decode() if the request body is empty. We
        // check for this with errors.Is() and return a plain-english error message 
        // instead.
        case errors.Is(err, io.EOF):
            return errors.New("body must not be empty")

        // A json.InvalidUnmarshalError error will be returned if we pass a non-nil 
        // pointer to Decode(). We catch this and panic, rather than returning an error
        // to our handler. At the end of this chapter we'll talk about panicking
        // versus returning errors, and discuss why it's an appropriate thing to do in 
        // this specific situation.
        case errors.As(err, &invalidUnmarshalError):
            panic(err)

        // For anything else, return the error message as-is.
        default:
            return err
        }
    }

    return nil
}