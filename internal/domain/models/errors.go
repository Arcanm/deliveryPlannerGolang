package models

import "errors"

var (
	// ErrRouteHasPendingPackages is returned when trying to complete a route with pending packages
	ErrRouteHasPendingPackages = errors.New("cannot complete route: there are pending packages")
)
