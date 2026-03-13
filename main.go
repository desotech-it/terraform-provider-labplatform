package main

import (
	"context"
	"log"

	"github.com/desotech-it/terraform-provider-labplatform/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/desotech-it/labplatform",
	})
	if err != nil {
		log.Fatal(err)
	}
}
