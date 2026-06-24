package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const maxTimestampDrift = 30

func GetUserId (c echo.Context) (uuid.UUID, error) {
	// Достает секрет из контекста, положенный туда мидлварем
	secret, ok := c.Get("gateway_secret").(string)
	if !ok || secret == "" {
		return uuid.Nil, echo.NewHTTPError(500, "Gateway secret not configured")
	}

	idStr := c.Request().Header.Get("X-User-ID")
	role := c.Request().Header.Get("X-User-Role")
	timestamp := c.Request().Header.Get("X-Gateway-Timestamp")
	signature := c.Request().Header.Get("X-Gateway-Signature")
	


	if idStr == "" || timestamp == "" || signature == "" {
		return uuid.Nil, echo.NewHTTPError(401, "Missing gateway headers") 
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(400, "Invalid Timestamp")
	}
	if math.Abs(float64(time.Now().Unix()-ts)) > maxTimestampDrift {
		return uuid.Nil, echo.NewHTTPError(403, "Request expired")
	}

	path := c.Request().Header.Get("X-Gateway-Path")
	if path == "" {
		path = c.Request().URL.Path
	}
	message := fmt.Sprintf("%s:%s:%s:%s", idStr, role, timestamp, path)
	
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return uuid.Nil, echo.NewHTTPError(401, "Invalid gateway signature")
	}

	uid, err := uuid.Parse(idStr)
	
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(401, "Invalid User ID format")
	}

	return uid, nil
}