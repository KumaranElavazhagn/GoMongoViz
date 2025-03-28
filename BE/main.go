package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"gomongoviz/database"
	"gomongoviz/handlers"
	domain "gomongoviz/repository"
	"gomongoviz/service"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize MongoDB connection
	// This establishes a connection to the MongoDB Atlas cluster
	mongoClient := database.ConnectMongo()

	// Ensure the MongoDB connection is properly closed when the application shuts down
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Initialize the application layers following the Clean Architecture pattern:
	// Database -> Repository -> Service -> Handler -> Router

	// Set up the repository layer with database connection
	repositoryDb := domain.NewRepositoryDefault(mongoClient)

	// Set up the service layer with repository
	svc := service.NewService(repositoryDb)

	// Set up the handler layer with service
	h := handlers.NewHandler(svc)

	// Initialize router with Gorilla Mux
	router := mux.NewRouter()

	// Configure CORS middleware to allow cross-origin requests
	// This is essential for the frontend to communicate with the API
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                                                                     // Allow all origins - in production, restrict to your frontend domain
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},                       // Must include OPTIONS for preflight requests
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}, // Content-Type is crucial for file uploads
		ExposedHeaders:   []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,  // Allow credentials such as cookies
		MaxAge:           86400, // 24 hours for preflight cache - reduces OPTIONS requests
	})

	// Set up logging middleware first to log all requests
	// This helps troubleshoot issues with API calls, especially file uploads
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("REQUEST: %s %s %s - Headers: %v", r.RemoteAddr, r.Method, r.URL, r.Header)
			next.ServeHTTP(w, r)
		})
	})

	// Define API routes with their corresponding handlers
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/objects", h.GetUniqueObjectIDs).Methods("GET")        // Get all unique object IDs
	api.HandleFunc("/ports/{objectId}", h.GetPorts).Methods("GET")         // Get ports for a specific object
	api.HandleFunc("/data/{objectId}", h.GetDataByObjectID).Methods("GET") // Get data for a specific object

	// Special handling for the upload endpoint
	// For file uploads, we need to handle both POST and OPTIONS methods
	// OPTIONS is used for CORS preflight requests from the browser
	api.HandleFunc("/upload", h.UploadCSV).Methods("POST", "OPTIONS") // Upload and process CSV data

	// Test route to check if the API is working
	api.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","message":"API is running"}`))
	}).Methods("GET")

	// Apply CORS middleware to the router
	handler := c.Handler(router)

	// Start the HTTP server on port 8080
	log.Printf("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
