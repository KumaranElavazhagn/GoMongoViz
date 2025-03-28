package model

import (
	"time"
)

// SensorData represents the time series data collected from sensors
// It contains various metrics and measurements from the monitoring system
type SensorData struct {
	ID              string    `bson:"_id,omitempty" json:"id"`                    // MongoDB document ID
	Timestamp       time.Time `bson:"timestamp" json:"timestamp"`                 // Time when the data was recorded
	ObjectID        float64   `bson:"object_id" json:"object_id"`                 // Identifier for the monitored object
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`               // Time when the record was created
	State           float64   `bson:"state" json:"state"`                         // Current state of the object
	ControllerError float64   `bson:"controller_error" json:"controller_error"`   // Error value from the controller
	AI3             float64   `bson:"ai3" json:"ai3"`                             // Analog input 3 reading
	AI5             float64   `bson:"ai5" json:"ai5"`                             // Analog input 5 reading
	FWVersion       string    `bson:"fw_version" json:"fw_version"`               // Firmware version
	PortNum         float64   `bson:"port_num" json:"port_num"`                   // Port number
	QCharge         float64   `bson:"q_charge" json:"q_charge"`                   // Charge value
	Voltage         float64   `bson:"voltage" json:"voltage"`                     // Current voltage
	VoltageSetPoint float64   `bson:"voltage_set_point" json:"voltage_set_point"` // Target voltage value
	Command         float64   `bson:"command" json:"command"`                     // Command value
	AI1             float64   `bson:"ai1" json:"ai1"`                             // Analog input 1 reading
	TargetQ         float64   `bson:"target_q" json:"target_q"`                   // Target charge value
	Current         float64   `bson:"current" json:"current"`                     // Current amperage
	AI4             float64   `bson:"ai4" json:"ai4"`                             // Analog input 4 reading
	VendorID        string    `bson:"vendor_id" json:"vendor_id"`                 // Vendor identifier
	AI2             float64   `bson:"ai2" json:"ai2"`                             // Analog input 2 reading
	SupplyCurrent   float64   `bson:"supply_current" json:"supply_current"`       // Supply current value
	StepNumber      float64   `bson:"step_number" json:"step_number"`             // Step number in the sequence
	VoltageDrop     float64   `bson:"voltage_drop" json:"voltage_drop"`           // Voltage drop measurement
	LiteID          string    `bson:"lite_id" json:"lite_id"`                     // Lite identifier
	VOCMode         float64   `bson:"voc_mode" json:"voc_mode"`                   // VOC (Voltage Open Circuit) mode
	VOC             float64   `bson:"voc" json:"voc"`                             // Voltage Open Circuit value
	ReadError       bool      `bson:"read_error" json:"read_error"`               // Indicates if there was an error during reading
	TargetVOC       float64   `bson:"target_voc" json:"target_voc"`               // Target Voltage Open Circuit value
	SupplyVolt      float64   `bson:"supply_volt" json:"supply_volt"`             // Supply voltage
	VOCState        float64   `bson:"voc_state" json:"voc_state"`                 // State of Voltage Open Circuit
	VOCExit         float64   `bson:"voc_exit" json:"voc_exit"`                   // VOC exit condition
}

// SensorDataRes is the response structure for sensor data queries
// It includes both the data and the total count for pagination purposes
type SensorDataRes struct {
	SensorData []SensorData // Array of sensor data records
	Total      int64        // Total number of records matching the query
}

// PortInfo represents information about a port associated with an object
type PortInfo struct {
	PortNum float64 `bson:"port_num" json:"portNum"` // Port number identifier
}

// ObjectInfo represents information about a monitored object
type ObjectInfo struct {
	ObjectID float64 `bson:"object_id" json:"objectId"` // Object identifier
}
