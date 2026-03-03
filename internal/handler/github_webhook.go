package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"

	"github.com/tritrongnguyen/repo-reviewer.git/pkg/logger"
	"go.uber.org/zap"
)

func GithubWebhook(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if secret == "" {
		logger.Log.Error("Missing webhook secret")
		http.Error(w, "Server misconfigured", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("Failed reading body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	signature := r.Header.Get("X-Hub-Signature-256")
	if !verifySignature(secret, body, signature) {
		logger.Log.Warn("Invalid signature")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	eventType := r.Header.Get("X-GitHub-Event")

	logger.Log.Info("Webhook received",
		zap.String("event", eventType),
		zap.String("delivery", r.Header.Get("X-GitHub-Delivery")),
	)

	// TODO: handle event types
	switch eventType {
	case "pull_request":
		handlePullRequest(body)

	default:
		logger.Log.Info("Unhandled event", zap.String("type", eventType))
	}

	w.WriteHeader(http.StatusOK)
}

func verifySignature(secret string, payload []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMac := mac.Sum(nil)

	expectedSignature := "sha256=" + hex.EncodeToString(expectedMac)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func handlePullRequest(body []byte) {
	logger.Log.Info("Pull request event received")
	// TODO: parse JSON into struct
}
