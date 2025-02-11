package services

import (
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)

func Settings(e *core.RequestEvent) error {

	return e.JSON(http.StatusOK, map[string]bool{"success": true})
}
