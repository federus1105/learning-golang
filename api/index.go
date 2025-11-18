package handler

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	Serve(w, r) // panggil fungsi utama Gin
}
