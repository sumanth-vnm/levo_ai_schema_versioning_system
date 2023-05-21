
package api

import (
	"github.com/gorilla/mux"
	"example.com/levo_app/controller"
)

// RegisterRoutes registers the API routes
func RegisterRoutes(handler *controller.APIHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/upload/schema", handler.UploadSchemaHandler).Methods("POST")
	r.HandleFunc("/getSchemaByVersion/{filename}/{version}", handler.GetSchemaHandler).Methods("GET")
	r.HandleFunc("/getLatestSchema/{filename}", handler.GetLatestSchemaHandler).Methods("GET")
	r.HandleFunc("/getAllVersions/{filename}", handler.GetAllVersionsHandler).Methods("GET")

	return r
}
