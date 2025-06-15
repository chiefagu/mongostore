package mongostore

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

func TestMongoDB(t *testing.T) {
	ctx := context.Background()

	// Set up MongoDB container
	mongoC, err := mongodb.Run(ctx, "iswprodacr.azurecr.io/mongo:7")
	if err != nil {
		t.Fatalf("Failed to start MongoDB container: %s", err)
	}
	defer mongoC.Terminate(ctx)

	// Get container's address
	host, err := mongoC.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get MongoDB container host: %s", err)
	}

	port, err := mongoC.MappedPort(ctx, "27017")
	if err != nil {
		t.Fatalf("Failed to get MongoDB container port: %s", err)
	}

	dsn := "mongodb://" + host + ":" + port.Port()

	viper.Set("mongodb_dsn", dsn)
	defer viper.Reset()

	t.Run("initMongo", func(t *testing.T) {
		dsn := viper.GetString("mongodb_dsn")

		Init(ctx, dsn)

		if Conn == nil {
			t.Errorf("Expected Conn to be non-nil")
		}
	})
	t.Run("pingMonodb", func(t *testing.T) {
		err = pingMongoDB(ctx)

		if err != nil {
			t.Errorf("Failed to ping MongoDB server: %s", err)
		}

	})

	t.Run("shutdown", func(t *testing.T) {
		Shutdown(ctx)

		if Conn != nil {
			t.Errorf("Expected Conn to be nil after shutdown")
		}

	})
}
