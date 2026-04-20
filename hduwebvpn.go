package hduwebvpn

import (
	"github.com/U1traVeno/hduwebvpn/client"
	"github.com/U1traVeno/hduwebvpn/handler"
	"github.com/U1traVeno/hduwebvpn/request"
	"github.com/U1traVeno/hduwebvpn/service"
	"github.com/U1traVeno/hduwebvpn/transport"
)

// Re-export types from transport package
type Mode = transport.Mode

const (
	WebVPNMode = transport.WebVPNMode
	DirectMode = transport.DirectMode
)

// Re-export types from client package
type ClientOption = client.ClientOption

// Re-export types from request package
type Request = request.Request
type Response = request.Response
type RealRequest = request.RealRequest

// Re-export types from service package
type Service = service.Service

// Re-export types from middleware package
type Handler = handler.Handler
type Context = handler.Context

// Re-export Transport interface
type Transport = transport.Transport

// Re-export functions from client package
var WithUsername = client.WithUsername
var WithPassword = client.WithPassword
var WithMode = client.WithMode
var NewClient = client.NewClient

// Re-export functions from middleware package
var NewContext = handler.NewContext

// Re-export functions from service package
// (Service methods are accessed via the type, so no re-export needed)