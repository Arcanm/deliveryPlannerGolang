package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Arcanm/deliveryPlannerGolang/internal/application/services"
	"github.com/Arcanm/deliveryPlannerGolang/internal/domain/models"
)

// RouteHandler handles HTTP requests for routes
type RouteHandler struct {
	service *services.RouteService
}

// NewRouteHandler creates a new route handler
func NewRouteHandler(service *services.RouteService) *RouteHandler {
	return &RouteHandler{
		service: service,
	}
}

// RegisterRoutes registers the route routes
func (h *RouteHandler) RegisterRoutes(router *gin.Engine) {
	routes := router.Group("/routes")
	{
		routes.POST("", h.CreateRoute)
		routes.GET("/:id", h.GetRoute)
		routes.GET("", h.ListRoutes)
		routes.PUT("/:id", h.UpdateRoute)
		routes.PATCH("/:id/status", h.UpdateRouteStatus)
		routes.POST("/:id/packages", h.AddPackagesToRoute)
		routes.PATCH("/:id/packages/:package_id/delivered", h.UpdatePackageDeliveryStatus)
		routes.DELETE("/:id", h.DeleteRoute)
	}
}

// CreateRouteRequest represents the request body for creating a route
type CreateRouteRequest struct {
	DriverID primitive.ObjectID `json:"driver_id" binding:"required"`
	Date     time.Time          `json:"date" binding:"required"`
}

// CreateRoute handles the creation of a new route
func (h *RouteHandler) CreateRoute(c *gin.Context) {
	var req CreateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := h.service.CreateRoute(c.Request.Context(), req.DriverID, req.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, route)
}

// GetRoute handles retrieving a route by ID
func (h *RouteHandler) GetRoute(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route ID"})
		return
	}

	route, err := h.service.GetRoute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if route == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
		return
	}

	c.JSON(http.StatusOK, route)
}

// ListRoutes handles retrieving all routes
func (h *RouteHandler) ListRoutes(c *gin.Context) {
	routes, err := h.service.ListRoutes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routes)
}

// UpdateRouteRequest represents the request body for updating a route
type UpdateRouteRequest struct {
	DriverID            primitive.ObjectID `json:"driver_id" binding:"required"`
	Date                time.Time          `json:"date" binding:"required"`
	EstimatedDistanceKm float64            `json:"estimated_distance_km"`
	EstimatedTimeMin    int                `json:"estimated_time_min"`
}

// UpdateRoute handles updating a route
func (h *RouteHandler) UpdateRoute(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route ID"})
		return
	}

	var req UpdateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := h.service.GetRoute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if route == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
		return
	}

	route.DriverID = req.DriverID
	route.Date = req.Date
	route.EstimatedDistanceKm = req.EstimatedDistanceKm
	route.EstimatedTimeMin = req.EstimatedTimeMin

	if err := h.service.UpdateRoute(c.Request.Context(), route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, route)
}

// UpdateRouteStatusRequest represents the request body for updating a route's status
type UpdateRouteStatusRequest struct {
	Status models.RouteStatus `json:"status" binding:"required"`
}

// UpdateRouteStatus handles updating a route's status
func (h *RouteHandler) UpdateRouteStatus(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route ID"})
		return
	}

	var req UpdateRouteStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateRouteStatus(c.Request.Context(), id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// AddPackagesToRouteRequest represents the request body for adding packages to a route
type AddPackagesToRouteRequest struct {
	PackageIDs []primitive.ObjectID `json:"package_ids" binding:"required"`
}

// AddPackagesToRoute handles adding packages to a route
func (h *RouteHandler) AddPackagesToRoute(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route ID"})
		return
	}

	var req AddPackagesToRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddPackagesToRoute(c.Request.Context(), id, req.PackageIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// UpdatePackageDeliveryStatus handles updating a package's delivery status in a route
func (h *RouteHandler) UpdatePackageDeliveryStatus(c *gin.Context) {
	routeID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route ID"})
		return
	}

	packageID, err := primitive.ObjectIDFromHex(c.Param("package_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid package ID"})
		return
	}

	if err := h.service.UpdatePackageDeliveryStatus(c.Request.Context(), routeID, packageID, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// DeleteRoute handles deleting a route
func (h *RouteHandler) DeleteRoute(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route ID"})
		return
	}

	if err := h.service.DeleteRoute(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
