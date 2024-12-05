package api

type InfoService interface {
	// GetInfo returns the info of the server
	getInfo() error
	getDetail() error
}
