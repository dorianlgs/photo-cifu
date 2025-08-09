package handlers

import (
	"io/fs"
	"net/http"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

// RegisterRoutes registers all API routes
func (h *Handlers) RegisterRoutes(router *router.Router[*core.RequestEvent]) {
	// API prefix for all PhotoCifu endpoints
	const apiPrefix = "/api/photocifu"

	// Gallery routes
	router.POST(apiPrefix+"/gallery/create", h.CreateGallery).
		Bind(apis.RequireAuth())

	// Workflow routes
	router.POST(apiPrefix+"/workflow/create", h.CreateWorkflow).
		Bind(apis.RequireAuth())

	// Signal routes
	router.POST(apiPrefix+"/signal/send", h.SendSignal).
		Bind(apis.RequireAuth())

	// Settings routes
	router.POST(apiPrefix+"/settings", h.UpdateSettings).
		Bind(apis.RequireAuth())
}

// RegisterStaticRoutes registers static file serving
func RegisterStaticRoutes(router *router.Router[*core.RequestEvent], distFS fs.FS, indexFallback bool) {
	if !router.HasRoute(http.MethodGet, "/{path...}") {
		router.GET("/{path...}", apis.Static(distFS, indexFallback)).
			Bind(apis.Gzip())
	}
}