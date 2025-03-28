package repository

import (
	"context"
	"encoding/json"
	"gomongoviz/model"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ConnectedMongoClient is a global variable to hold the MongoDB client connection
var ConnectedMongoClient *mongo.Client

// RepositoryDefault is the concrete implementation of the Repository interface
// It handles database operations for the application
type RepositoryDefault struct {
	Client *mongo.Client // MongoDB client connection
}

// Repository defines the interface for data access operations
// It abstracts the database layer, making it easier to test and replace implementations
type Repository interface {
	GetUniqueObjectIDs() ([]model.ObjectInfo, error)
	GetPorts(objectID int) ([]model.PortInfo, error)
	GetDataByObjectID(objectID string, portNum string) (*model.SensorDataRes, error)
	SaveSensorData(data []model.SensorData) error
}

// GetUniqueObjectIDs retrieves a list of all unique object IDs in the database
// It uses MongoDB aggregation to group and sort the results
func (r RepositoryDefault) GetUniqueObjectIDs() ([]model.ObjectInfo, error) {
	collection := r.Client.Database("gomongoviz").Collection("sensor_data")

	// Create an aggregation pipeline to get distinct object_ids
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":       "$object_id",
			"object_id": bson.M{"$first": "$object_id"},
		}},
		{"$project": bson.M{
			"_id": 0,
		}},
		{"$sort": bson.M{
			"object_id": 1,
		}},
	}

	// Print the pipeline for debugging
	pipelineJSON, err := json.MarshalIndent(pipeline, "", "  ")
	if err != nil {
		log.Printf("Error marshaling pipeline: %v", err)
	} else {
		log.Printf("MongoDB Aggregation Pipeline:\n%s", string(pipelineJSON))
	}

	// Execute the aggregation pipeline
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Parse results into model objects
	var results []model.ObjectInfo
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	// Log the results for debugging
	log.Printf("Found %d unique object IDs", len(results))
	for _, obj := range results {
		log.Printf("Object ID: %v", obj.ObjectID)
	}

	return results, nil
}

// GetPorts retrieves a list of all ports for a specific object ID
// It uses MongoDB aggregation to find distinct port numbers for the given object
func (r RepositoryDefault) GetPorts(objectID int) ([]model.PortInfo, error) {
	collection := r.Client.Database("gomongoviz").Collection("sensor_data")

	// Create an aggregation pipeline to get distinct port numbers for the given object ID
	pipeline := []bson.M{
		{"$match": bson.M{"object_id": float64(objectID)}},
		{"$group": bson.M{
			"_id":      "$port_num",
			"port_num": bson.M{"$first": "$port_num"},
		}},
		{"$project": bson.M{
			"_id": 0,
		}},
	}

	// Print the pipeline for debugging
	pipelineJSON, err := json.MarshalIndent(pipeline, "", "  ")
	if err != nil {
		log.Printf("Error marshaling pipeline: %v", err)
	} else {
		log.Printf("MongoDB Aggregation Pipeline:\n%s", string(pipelineJSON))
	}

	// Execute the aggregation pipeline
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Parse results into model objects
	var results []model.PortInfo
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	// Log the results for debugging
	log.Printf("Found %d unique ports", len(results))
	for _, port := range results {
		log.Printf("Port: %v", port.PortNum)
	}

	return results, nil
}

// GetDataByObjectID retrieves sensor data for a specific object ID and optionally filtered by port number
// It returns both the matching data and the total count of matching documents
func (r RepositoryDefault) GetDataByObjectID(objectID string, portNum string) (*model.SensorDataRes, error) {
	collection := r.Client.Database("gomongoviz").Collection("sensor_data")

	// Convert objectID string to float64
	objectIDFloat, err := strconv.ParseFloat(objectID, 64)
	if err != nil {
		return nil, err
	}

	// Build the filter based on objectID and optional portNum
	filter := bson.M{"object_id": objectIDFloat}

	if portNum != "" {
		portNumFloat, err := strconv.ParseFloat(portNum, 64)
		if err != nil {
			return nil, err
		}
		filter["port_num"] = portNumFloat
	}

	// Log the filter being used
	log.Printf("MongoDB filter: %+v", filter)

	// First, get count of matching documents
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	log.Printf("Found %d matching documents", count)

	// Execute the query
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Parse results into model objects
	results := make([]model.SensorData, 0)
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	// Log the number of results
	log.Printf("Retrieved %d documents", len(results))

	return &model.SensorDataRes{
		SensorData: results,
		Total:      count,
	}, nil
}

// SaveSensorData saves a batch of sensor data records to the database
// Used for CSV file uploads to add new sensor data to the collection
func (r RepositoryDefault) SaveSensorData(data []model.SensorData) error {
	if len(data) == 0 {
		return nil // No data to save
	}

	collection := r.Client.Database("gomongoviz").Collection("sensor_data")
	
	// Convert data to interface slice for bulk insert
	var documents []interface{}
	for _, item := range data {
		documents = append(documents, item)
	}
	
	// Insert all documents in a single operation
	log.Printf("Inserting %d documents to MongoDB", len(documents))
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Printf("Error inserting documents: %v", err)
		return err
	}
	
	log.Printf("Successfully inserted %d documents", len(documents))
	return nil
}

// NewRepositoryDefault creates a new instance of the default repository implementation
// It takes a MongoDB client as input and returns a Repository interface
func NewRepositoryDefault(client *mongo.Client) Repository {
	return RepositoryDefault{
		Client: client,
	}
}
