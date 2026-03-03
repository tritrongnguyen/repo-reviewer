package handler

import (
	"net/http"

	"github.com/tritrongnguyen/repo-reviewer.git/pkg/helper"
)

type healthResponse struct {
	Status string `json:"status"`
}

func Health(w http.ResponseWriter, r *http.Request) {

	helper.RespondWithJson(w, 200, healthResponse{
		Status: "ok",
	})
}
