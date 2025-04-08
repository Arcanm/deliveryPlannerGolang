package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Arcanm/deliveryPlannerGolang/internal/application/services"
)

// PackageHandler handles HTTP requests for packages
type PackageHandler struct {
	packageService *services.PackageService
	routeService   *services.RouteService
}

// NewPackageHandler creates a new package handler
func NewPackageHandler(packageService *services.PackageService, routeService *services.RouteService) *PackageHandler {
	return &PackageHandler{
		packageService: packageService,
		routeService:   routeService,
	}
}

// RegisterRoutes registers the package routes
func (h *PackageHandler) RegisterRoutes(router *gin.Engine) {
	packages := router.Group("/api/v1/packages")
	{
		packages.POST("", h.CreatePackage)
		packages.GET("", h.ListPackages)
		packages.GET("/:id", h.GetPackage)
		packages.PUT("/:id", h.UpdatePackage)
		packages.DELETE("/:id", h.DeletePackage)
		packages.POST("/:id/assign", h.AssignToRoute)
		packages.POST("/:id/deliver", h.MarkAsDelivered)
		packages.GET("/route/:route_id", h.GetPackagesByRoute)
	}
}

// CreatePackageRequest represents the request body for creating a package
type CreatePackageRequest struct {
	TrackingNumber  string  `json:"tracking_number" binding:"required"`
	CustomerName    string  `json:"customer_name" binding:"required"`
	CustomerAddress string  `json:"customer_address" binding:"required"`
	CustomerPhone   string  `json:"customer_phone" binding:"required"`
	WeightKg        float64 `json:"weight_kg" binding:"required,gt=0"`
	VolumeM3        float64 `json:"volume_m3" binding:"required,gt=0"`
}

// CreatePackage handles the creation of a new package
func (h *PackageHandler) CreatePackage(c *gin.Context) {
	var req CreatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pkg, err := h.packageService.CreatePackage(c.Request.Context(), req.TrackingNumber, req.CustomerName, req.CustomerAddress, req.CustomerPhone, req.WeightKg, req.VolumeM3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pkg)
}

// GetPackage handles retrieving a package by ID
func (h *PackageHandler) GetPackage(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid package id"})
		return
	}

	pkg, err := h.packageService.GetPackage(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "package not found"})
		return
	}

	c.JSON(http.StatusOK, pkg)
}

// ListPackages handles retrieving all packages
func (h *PackageHandler) ListPackages(c *gin.Context) {
	packages, err := h.packageService.ListPackages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, packages)
}

// UpdatePackageRequest represents the request body for updating a package
type UpdatePackageRequest struct {
	TrackingNumber  string  `json:"tracking_number" binding:"required"`
	CustomerName    string  `json:"customer_name" binding:"required"`
	CustomerAddress string  `json:"customer_address" binding:"required"`
	CustomerPhone   string  `json:"customer_phone" binding:"required"`
	WeightKg        float64 `json:"weight_kg" binding:"required,gt=0"`
	VolumeM3        float64 `json:"volume_m3" binding:"required,gt=0"`
}

// UpdatePackage handles updating a package
func (h *PackageHandler) UpdatePackage(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid package id"})
		return
	}

	var req UpdatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pkg, err := h.packageService.UpdatePackage(c.Request.Context(), id, req.TrackingNumber, req.CustomerName, req.CustomerAddress, req.CustomerPhone, req.WeightKg, req.VolumeM3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pkg)
}

// DeletePackage handles deleting a package
func (h *PackageHandler) DeletePackage(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid package id"})
		return
	}

	if err := h.packageService.DeletePackage(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignToRouteRequest represents the request body for assigning a package to a route
type AssignToRouteRequest struct {
	RouteID string `json:"route_id" binding:"required"`
}

// AssignToRoute handles assigning a package to a route
func (h *PackageHandler) AssignToRoute(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid package id"})
		return
	}

	var req AssignToRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	routeID, err := primitive.ObjectIDFromHex(req.RouteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route id"})
		return
	}

	err = h.routeService.AddPackagesToRoute(c.Request.Context(), routeID, []primitive.ObjectID{id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// MarkAsDelivered handles marking a package as delivered
func (h *PackageHandler) MarkAsDelivered(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid package id"})
		return
	}

	pkg, err := h.packageService.MarkAsDelivered(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pkg)
}

// GetPackagesByRoute handles retrieving packages by route
func (h *PackageHandler) GetPackagesByRoute(c *gin.Context) {
	routeID, err := primitive.ObjectIDFromHex(c.Param("route_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route id"})
		return
	}

	route, err := h.routeService.GetRoute(c.Request.Context(), routeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, route.Packages)
}
