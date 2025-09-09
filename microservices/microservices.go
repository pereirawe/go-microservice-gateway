package microservices

import (
	"log"

	"github.com/pereirawe/go-microservice-gateway/config"
)

// GetMicroservices returns a map of microservices and their URLs.
func GetMicroservices() map[string]string {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	microservices := map[string]string{
		"ms-ai":       cfg.MSAi,
		"ms-payments": cfg.MSPayments,
		"ms-bi":       cfg.MSBi,
		"ms-reports":  cfg.MSReports,
	}
	return microservices
}
