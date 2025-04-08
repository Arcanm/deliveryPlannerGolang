package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Arcanm/deliveryPlannerGolang/internal/application/services"
	"github.com/Arcanm/deliveryPlannerGolang/internal/domain/models"
)

// DriverHandler handles HTTP requests for drivers
type DriverHandler struct {
	service *services.DriverService
}

// NewDriverHandler creates a new driver handler
func NewDriverHandler(service *services.DriverService) *DriverHandler {
	return &DriverHandler{
		service: service,
	}
}

// RegisterRoutes registers the driver routes
func (h *DriverHandler) RegisterRoutes(router *gin.Engine) {
	drivers := router.Group("/api/v1/drivers")
	{
		drivers.POST("", h.CreateDriver)
		drivers.GET("", h.ListDrivers)
		drivers.GET("/:id", h.GetDriver)
		drivers.PUT("/:id", h.UpdateDriver)
		drivers.DELETE("/:id", h.DeleteDriver)
		drivers.GET("/:id/routes", h.GetDriverRoutes)
	}
}

// CreateDriverRequest represents the request body for creating a driver
type CreateDriverRequest struct {
	Name        string `json:"name" binding:"required"`
	VehicleType string `json:"vehicle_type" binding:"required,oneof=bike van truck"`
}

// CreateDriver handles the creation of a new driver
func (h *DriverHandler) CreateDriver(c *gin.Context) {
	var req CreateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehicleType := models.VehicleType(req.VehicleType)
	driver, err := h.service.CreateDriver(c.Request.Context(), req.Name, vehicleType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, driver)
}

// GetDriver handles retrieving a driver by ID
func (h *DriverHandler) GetDriver(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver id"})
		return
	}

	driver, err := h.service.GetDriver(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "driver not found"})
		return
	}

	c.JSON(http.StatusOK, driver)
}

// ListDrivers handles retrieving all drivers
func (h *DriverHandler) ListDrivers(c *gin.Context) {
	drivers, err := h.service.ListDrivers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, drivers)
}

// UpdateDriverRequest represents the request body for updating a driver
type UpdateDriverRequest struct {
	Name        string             `json:"name" binding:"required"`
	VehicleType models.VehicleType `json:"vehicle_type" binding:"required"`
	Active      bool               `json:"active"`
}

// UpdateDriver handles updating a driver
func (h *DriverHandler) UpdateDriver(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver ID"})
		return
	}

	var req UpdateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	driver, err := h.service.GetDriver(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if driver == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "driver not found"})
		return
	}

	driver.Name = req.Name
	driver.VehicleType = req.VehicleType
	driver.Active = req.Active

	updatedDriver, err := h.service.UpdateDriver(c.Request.Context(), id, req.Name, req.VehicleType, req.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedDriver)
}

// DeleteDriver handles deleting a driver
func (h *DriverHandler) DeleteDriver(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver ID"})
		return
	}

	if err := h.service.DeleteDriver(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetDriverRoutes handles retrieving all routes for a driver
func (h *DriverHandler) GetDriverRoutes(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver ID"})
		return
	}

	routes, err := h.service.GetDriverRoutes(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routes)
}
