package configs

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB                    *mongo.Client
	BatteryCollection     *mongo.Collection
	StationCollection     *mongo.Collection
	SwapEventCollection   *mongo.Collection
	TimeSeriesCollection  *mongo.Collection
	UncertaintyCollection *mongo.Collection
	AlertCollection       *mongo.Collection
	once                  sync.Once
)

func ConnectDB() {
	once.Do(func() {
		clientOptions := options.Client().ApplyURI(EnvMongoURI())

		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal("❌ Failed to connect to MongoDB:", err)
		}

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal("❌ Failed to ping MongoDB:", err)
		}

		fmt.Println("✅ Connected to MongoDB")

		DB = client
		db := client.Database("greenkwh")

		BatteryCollection = db.Collection("batteries")
		StationCollection = db.Collection("swap_stations")
		SwapEventCollection = db.Collection("swap_events")
		TimeSeriesCollection = db.Collection("battery_timeseries")
		UncertaintyCollection = db.Collection("uncertainties")
		AlertCollection = db.Collection("alerts")

		createIndexes(db)
	})
}

func createIndexes(db *mongo.Database) {
	ctx := context.TODO()

	// Battery: index on IMEI
	BatteryCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "imei", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	// TimeSeries: compound index for range queries
	TimeSeriesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "imei", Value: 1},
			{Key: "metric", Value: 1},
			{Key: "timestamp", Value: -1},
		},
	})

	// Uncertainty: index for unresolved lookups
	UncertaintyCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "imei", Value: 1},
			{Key: "resolved", Value: 1},
		},
	})

	// Alerts: index for filtering
	AlertCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "acknowledged", Value: 1},
			{Key: "severity", Value: 1},
			{Key: "created_at", Value: -1},
		},
	})

	// Stations: index on station_id
	StationCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "station_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	fmt.Println("✅ Database indexes created")
}
