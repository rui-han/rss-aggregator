package main

import "net/http"

func HandlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Something went wrong...")
}
