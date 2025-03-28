package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gomongoviz/model"
	"gomongoviz/service"

	"github.com/gorilla/mux"
)

// Handler implements the HTTP handlers for the API endpoints
// It uses the service layer to process business logic
type Handler struct {
	service *service.Service // Service for business logic operations
}

// NewHandler creates a new handler with the provided service
// This follows the dependency injection pattern
func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		service: svc,
	}
}

// GetPorts handles HTTP requests to retrieve ports for a specific object ID
// URL pattern: /api/ports/{objectId}
func (h *Handler) GetPorts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectID := vars["objectId"]
	if objectID == "" {
		http.Error(w, "objectId is required", http.StatusBadRequest)
		return
	}

	// Convert objectID string to int
	id, err := strconv.Atoi(objectID)
	if err != nil {
		http.Error(w, "invalid objectId: must be a number", http.StatusBadRequest)
		return
	}

	ports, err := h.service.GetPorts(id)
	if err != nil {
		var errMsg string
		if e, ok := err.(error); ok {
			errMsg = e.Error()
		} else {
			errMsg = "Internal server error"
		}
		writeResponse(w, http.StatusInternalServerError, errMsg)
		return
	}

	writeResponse(w, http.StatusOK, ports)
}

// GetUniqueObjectIDs handles HTTP requests to get all unique object IDs
// URL pattern: /api/objects
func (h *Handler) GetUniqueObjectIDs(w http.ResponseWriter, r *http.Request) {
	objectIDs, err := h.service.GetUniqueObjectIDs()
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	writeResponse(w, http.StatusOK, objectIDs)
}

// GetDataByObjectID handles HTTP requests to get sensor data for a specific object ID
// Optionally filtered by port number via query parameter
// URL pattern: /api/data/{objectId}?port_num=X
func (h *Handler) GetDataByObjectID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectID := vars["objectId"]
	portNum := r.URL.Query().Get("port_num")

	if objectID == "" {
		writeResponse(w, http.StatusBadRequest, "objectId is required")
		return
	}

	data, err := h.service.GetDataByObjectID(objectID, portNum)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, data)
}

// UploadCSV handles HTTP requests to upload CSV files containing sensor data
// URL pattern: /api/upload
// The CSV file should contain properly formatted sensor data with all required fields
func (h *Handler) UploadCSV(w http.ResponseWriter, r *http.Request) {
	// Log request content type and method for debugging
	// This is crucial for diagnosing Content-Type issues with file uploads
	// The Content-Type should contain 'multipart/form-data' with a boundary parameter
	log.Printf("Request Content-Type: %s, Method: %s", r.Header.Get("Content-Type"), r.Method)

	// Handle OPTIONS preflight request
	// Browsers send this before the actual POST request for CORS validation
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Maximum upload size of 10MB
	// ParseMultipartForm parses the multipart form including file uploads
	// This fails if the Content-Type header isn't set correctly by the client
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		// Common issue: Content-Type header is missing or incorrect
		// This happens when the frontend manually sets the Content-Type header
		// or when using axios without proper configuration for file uploads
		if r.Header.Get("Content-Type") == "" || !strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "multipart/form-data") {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error":   "Failed to parse form data",
				"message": "request Content-Type isn't multipart/form-data",
			})
		} else {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error":   "Failed to parse form data",
				"message": err.Error(),
			})
		}
		return
	}

	// Get file from request using the field name 'file'
	// The frontend must use this same field name when appending to FormData
	file, handler, err := r.FormFile("file")
	if err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]string{
			"error":   "Could not get file from request",
			"message": err.Error(),
		})
		return
	}
	defer file.Close()

	// Log file info for debugging
	log.Printf("Received file: %s, size: %d, content type: %s",
		handler.Filename, handler.Size, handler.Header.Get("Content-Type"))

	// Validate file type - more permissive check for different CSV mime types
	contentType := handler.Header.Get("Content-Type")
	isCSV := contentType == "text/csv" ||
		contentType == "application/csv" ||
		contentType == "application/vnd.ms-excel" ||
		contentType == "text/plain" ||
		contentType == "application/octet-stream" ||
		handler.Filename[len(handler.Filename)-4:] == ".csv"

	if !isCSV {
		writeResponse(w, http.StatusBadRequest, map[string]string{
			"error":   "Invalid file type",
			"message": fmt.Sprintf("Only CSV files are allowed. Received: %s", contentType),
		})
		return
	}

	// Parse CSV
	csvReader := csv.NewReader(file)

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]string{
			"error":   "Failed to read CSV header",
			"message": err.Error(),
		})
		return
	}

	// Log header for debugging
	log.Printf("CSV header: %v", header)

	// Map column indices
	headerMap := make(map[string]int)
	for i, column := range header {
		headerMap[column] = i
	}

	// Validate required fields
	requiredFields := []string{
		"timestamp", "object_id", "port_num", "voltage", "current",
		"supply_current", "supply_volt", "voltage_drop", "voc",
	}

	missingFields := []string{}
	for _, field := range requiredFields {
		if _, exists := headerMap[field]; !exists {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		writeResponse(w, http.StatusBadRequest, map[string]string{
			"error":   "Missing required fields in CSV",
			"message": fmt.Sprintf("The following required fields are missing: %v", missingFields),
		})
		return
	}

	// Process all rows
	var sensorData []model.SensorData
	lineNum := 1 // Start after header

	for {
		lineNum++
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error":   "Failed to read CSV row",
				"message": fmt.Sprintf("Error at line %d: %s", lineNum, err.Error()),
			})
			return
		}

		// Parse data from CSV row
		timestamp, err := time.Parse(time.RFC3339, record[headerMap["timestamp"]])
		if err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error":   "Invalid timestamp format",
				"message": fmt.Sprintf("Error at line %d: timestamp should be in RFC3339 format", lineNum),
			})
			return
		}

		objectID, err := strconv.ParseFloat(record[headerMap["object_id"]], 64)
		if err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error":   "Invalid object_id",
				"message": fmt.Sprintf("Error at line %d: object_id should be a number", lineNum),
			})
			return
		}

		portNum, err := strconv.ParseFloat(record[headerMap["port_num"]], 64)
		if err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error":   "Invalid port_num",
				"message": fmt.Sprintf("Error at line %d: port_num should be a number", lineNum),
			})
			return
		}

		// Parse numeric fields with error handling
		parseFloatField := func(field string) (float64, error) {
			value := record[headerMap[field]]
			if value == "" {
				return 0, nil
			}
			return strconv.ParseFloat(value, 64)
		}

		voltage, err := parseFloatField("voltage")
		if err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error":   "Invalid voltage",
				"message": fmt.Sprintf("Error at line %d: voltage should be a number", lineNum),
			})
			return
		}

		current, err := parseFloatField("current")
		if err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error":   "Invalid current",
				"message": fmt.Sprintf("Error at line %d: current should be a number", lineNum),
			})
			return
		}

		supplyCurrent, err := parseFloatField("supply_current")
		if err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error":   "Invalid supply_current",
				"message": fmt.Sprintf("Error at line %d: supply_current should be a number", lineNum),
			})
			return
		}

		supplyVolt, err := parseFloatField("supply_volt")
		if err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error":   "Invalid supply_volt",
				"message": fmt.Sprintf("Error at line %d: supply_volt should be a number", lineNum),
			})
			return
		}

		voltageDrop, err := parseFloatField("voltage_drop")
		if err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error":   "Invalid voltage_drop",
				"message": fmt.Sprintf("Error at line %d: voltage_drop should be a number", lineNum),
			})
			return
		}

		voc, err := parseFloatField("voc")
		if err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error":   "Invalid voc",
				"message": fmt.Sprintf("Error at line %d: voc should be a number", lineNum),
			})
			return
		}

		// Create new SensorData object
		data := model.SensorData{
			Timestamp:     timestamp,
			ObjectID:      objectID,
			PortNum:       portNum,
			Voltage:       voltage,
			Current:       current,
			SupplyCurrent: supplyCurrent,
			SupplyVolt:    supplyVolt,
			VoltageDrop:   voltageDrop,
			VOC:           voc,
			CreatedAt:     time.Now(),
		}

		// Add optional fields if present
		for _, field := range []string{
			"state", "controller_error", "ai1", "ai2", "ai3", "ai4", "ai5",
			"fw_version", "q_charge", "voltage_set_point", "command", "target_q",
			"vendor_id", "step_number", "lite_id", "voc_mode", "target_voc",
			"voc_state", "voc_exit",
		} {
			if index, exists := headerMap[field]; exists && index < len(record) {
				value := record[index]
				switch field {
				case "fw_version", "vendor_id", "lite_id":
					// String fields
					if field == "fw_version" {
						data.FWVersion = value
					} else if field == "vendor_id" {
						data.VendorID = value
					} else if field == "lite_id" {
						data.LiteID = value
					}
				default:
					// Numeric fields
					if value != "" {
						if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
							switch field {
							case "state":
								data.State = floatVal
							case "controller_error":
								data.ControllerError = floatVal
							case "ai1":
								data.AI1 = floatVal
							case "ai2":
								data.AI2 = floatVal
							case "ai3":
								data.AI3 = floatVal
							case "ai4":
								data.AI4 = floatVal
							case "ai5":
								data.AI5 = floatVal
							case "q_charge":
								data.QCharge = floatVal
							case "voltage_set_point":
								data.VoltageSetPoint = floatVal
							case "command":
								data.Command = floatVal
							case "target_q":
								data.TargetQ = floatVal
							case "step_number":
								data.StepNumber = floatVal
							case "voc_mode":
								data.VOCMode = floatVal
							case "target_voc":
								data.TargetVOC = floatVal
							case "voc_state":
								data.VOCState = floatVal
							case "voc_exit":
								data.VOCExit = floatVal
							}
						}
					}
				}
			}
		}

		// Handle read_error boolean field if present
		if index, exists := headerMap["read_error"]; exists && index < len(record) {
			value := record[index]
			if value == "true" {
				data.ReadError = true
			} else {
				data.ReadError = false
			}
		}

		sensorData = append(sensorData, data)
	}

	// Check if we have any data to save
	if len(sensorData) == 0 {
		writeResponse(w, http.StatusBadRequest, map[string]string{
			"error":   "No valid data found",
			"message": "The CSV file contains a header but no valid data rows",
		})
		return
	}

	// Send data to service layer for processing and storage
	err = h.service.SaveSensorData(sensorData)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]string{
			"error":   "Failed to save sensor data",
			"message": err.Error(),
		})
		return
	}

	writeResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Successfully uploaded %d sensor data records", len(sensorData)),
		"count":   len(sensorData),
	})
}

// writeResponse is a helper function to write JSON responses
// It sets appropriate headers, status code, and serializes the data to JSON
func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	// Set CORS headers to allow cross-origin requests
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

	w.WriteHeader(code)

	// If data is nil, just return with status code
	if data == nil {
		return
	}

	// Marshal data to JSON and write to response
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
