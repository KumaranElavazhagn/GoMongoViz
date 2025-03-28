import React, { useState, useEffect } from "react";
import { DataChart } from "./components/DataChart";
import { GetObject, GetPort, GetData } from "./api/api";
import { Field, ChartData } from "./types";
import "./App.css";

interface ObjectOption {
  value: string;
  label: string;
}

type ChartType = "line" | "bar" | "area" | "scatter";

interface SensorDataItem {
  timestamp: string;
  [key: string]: any;
}

// Interface for the upload modal state
interface UploadModalState {
  isOpen: boolean;
  isLoading: boolean;
  error: string | null;
  success: boolean;
}

function App() {
  // State management
  const [objectOptions, setObjectOptions] = useState<ObjectOption[]>([]);
  const [selectedObjectId, setSelectedObjectId] = useState<string>("");
  const [portNums, setPortNums] = useState<string[]>([]);
  const [selectedPortNum, setSelectedPortNum] = useState<string>("");
  const [startDate, setStartDate] = useState<Date | null>(null);
  const [endDate, setEndDate] = useState<Date | null>(null);
  const [selectedFields, setSelectedFields] = useState<string[]>([]);
  const [chartData, setChartData] = useState<ChartData | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [chartType, setChartType] = useState<ChartType>("line");
  const [activeTimeFilter, setActiveTimeFilter] = useState<number | null>(null);
  const [allSensorData, setAllSensorData] = useState<SensorDataItem[]>([]);
  const [uploadModal, setUploadModal] = useState<UploadModalState>({
    isOpen: false,
    isLoading: false,
    error: null,
    success: false
  });

  // Available sensor data fields that can be displayed
  const availableFields: Field[] = [
    { label: "Current", value: "current" },
    { label: "Voltage", value: "voltage" },
    { label: "Supply Current", value: "supply_current" },
    { label: "Supply Voltage", value: "supply_volt" },
    { label: "Voltage Drop", value: "voltage_drop" },
    { label: "VOC", value: "voc" },
  ];

  // Chart type options for visualization
  const chartTypeOptions = [
    { label: "Line Chart", value: "line" },
    { label: "Bar Chart", value: "bar" },
    { label: "Area Chart", value: "area" },
    { label: "Scatter Plot", value: "scatter" },
  ];

  // Predefined time filter options
  const timeFilters = [
    { label: "1 Hour", value: 1 },
    { label: "3 Hours", value: 3 },
    { label: "12 Hours", value: 12 },
    { label: "24 Hours", value: 24 },
  ];

  // Fetch devices on component mount
  useEffect(() => {
    fetchObjectIds();
  }, []);

  // Fetch ports and all data when device changes
  useEffect(() => {
    if (selectedObjectId) {
      fetchPortNums();
      fetchAllData("all"); // Fetch data for all ports initially
      setSelectedPortNum("");
    }
  }, [selectedObjectId]);

  // Fetch data for specific port when port changes
  useEffect(() => {
    if (selectedObjectId && selectedPortNum) {
      fetchAllData(selectedPortNum); // Fetch data for specific port when selected
    }
  }, [selectedPortNum]);

  // Update chart data when data or filters change
  useEffect(() => {
    if (allSensorData.length > 0 && selectedFields.length > 0) {
      filterAndPrepareData();
    }
  }, [allSensorData, selectedFields, startDate, endDate]);

  // Fetch available devices from API
  const fetchObjectIds = async () => {
    try {
      setLoading(true);
      const response = await GetObject();
      
      if (response?.data && Array.isArray(response.data)) {
        const options = response.data.map((item: { objectId: number }) => {
          const optionValue = item.objectId.toString();
          return {
            value: optionValue,
            label: `Device ${optionValue}`
          };
        });
        
        setObjectOptions(options);
      }
    } catch (error) {
      console.error("Error fetching object IDs:", error);
    } finally {
      setLoading(false);
    }
  };

  // Fetch available ports for the selected device
  const fetchPortNums = async () => {
    try {
      setLoading(true);
      const response = await GetPort(selectedObjectId);
      
      if (response?.data && Array.isArray(response.data)) {
        const portNumbers = response.data.map((item: { portNum: number }) =>
          item.portNum.toString()
        );
        setPortNums(portNumbers);
      } else {
        setPortNums([]);
      }
    } catch (error) {
      console.error("Error fetching port numbers:", error);
      setPortNums([]);
    } finally {
      setLoading(false);
    }
  };

  // Fetch sensor data for selected device and port
  const fetchAllData = async (portNumber: string) => {
    if (!selectedObjectId) {
      return;
    }
    
    setLoading(true);
    try {
      const response = await GetData(
        selectedObjectId,
        portNumber,
      );

      if (response?.data?.SensorData && Array.isArray(response.data.SensorData)) {
        setAllSensorData(response.data.SensorData);
        
        // Auto-select first field if none selected
        if (selectedFields.length === 0 && availableFields.length > 0) {
          setSelectedFields([availableFields[0].value]);
        }
      } else {
        setAllSensorData([]);
      }
    } catch (error) {
      console.error("Error fetching data:", error);
      setAllSensorData([]);
    } finally {
      setLoading(false);
    }
  };

  // Handle file input change
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    // Reset any previous errors when a new file is selected
    setUploadModal(prev => ({
      ...prev,
      error: null,
      success: false
    }));
  };

  // Handle CSV file upload
  // This function handles the upload of CSV files to the server
  // It manually creates a FormData object and uses the fetch API directly
  // instead of axios to avoid Content-Type header issues
  const handleCSVUpload = async (e: React.FormEvent<HTMLFormElement>) => {
    // Prevent default form submission behavior
    e.preventDefault();
    
    // Get the file input element and selected file
    const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement;
    const file = fileInput.files?.[0];
    
    // Validate that a file has been selected
    if (!file) {
      setUploadModal(prev => ({
        ...prev,
        error: "Please select a CSV file to upload"
      }));
      return;
    }

    // Reset modal state and show loading indicator
    setUploadModal(prev => ({
      ...prev,
      isLoading: true,
      error: null,
      success: false
    }));

    // Create proper FormData object to handle file uploads
    // FormData automatically sets the correct Content-Type with boundary
    const formData = new FormData();
    formData.append('file', file); // 'file' is the field name expected by the backend

    try {
      console.log('Uploading file:', file.name, 'size:', file.size, 'type:', file.type);
      
      // IMPORTANT: Using direct URL for upload rather than the Client.ts axios wrapper
      // The fetch API handles multipart/form-data uploads better than axios in this case
      // We must use the full URL (not relative) for direct API access
      const uploadUrl = 'http://localhost:8080/api/upload';
      
      // Log the file we're uploading for debugging
      console.log('FormData contains file:', file.name);
      
      // Use the fetch API for the upload
      // CRITICAL: DO NOT manually set Content-Type header for multipart/form-data uploads
      // The browser will automatically set the correct Content-Type with proper boundary
      // Setting it manually will cause the backend to reject the upload
      const response = await fetch(uploadUrl, {
        method: 'POST',
        body: formData,
        // No Content-Type header here - browser sets it automatically with boundary
      });

      // Get response text first to make debugging easier
      const responseText = await response.text();
      console.log('Response text:', responseText);
      
      // If not OK status, handle error
      if (!response.ok) {
        console.error('Server response:', responseText);
        
        // Try to parse as JSON if possible
        let errorMessage = 'Error uploading CSV';
        try {
          const errorData = JSON.parse(responseText);
          errorMessage = errorData.message || 'Unknown server error';
        } catch (parseError) {
          // If can't parse JSON, use the response text (truncated if too long)
          errorMessage = responseText.length > 100 
            ? `${responseText.substring(0, 100)}...` 
            : responseText;
        }
        
        throw new Error(errorMessage);
      }
      
      // Parse successful response as JSON (if it's JSON)
      let result;
      try {
        result = JSON.parse(responseText);
        console.log('Upload success:', result);
      } catch (e) {
        console.log('Response is not JSON:', responseText);
      }

      // Update modal state on success
      setUploadModal(prev => ({
        ...prev,
        isLoading: false,
        success: true
      }));

      // Refresh data after successful upload to see new records
      fetchObjectIds();
    } catch (error) {
      console.error('Upload error:', error);
      setUploadModal(prev => ({
        ...prev,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Unknown error occurred'
      }));
    }
  };

  // Transform sensor data into chart-compatible format
  const filterAndPrepareData = () => {
    if (selectedFields.length === 0 || allSensorData.length === 0) {
      setChartData(null);
      return;
    }

    // Determine the data to use based on date range
    let dataToUse = allSensorData;
    
    // Only filter by date if both dates are set
    if (startDate && endDate) {
      dataToUse = allSensorData.filter(item => {
        const itemDate = new Date(item.timestamp);
        return itemDate >= startDate && itemDate <= endDate;
      });
    }

    if (dataToUse.length === 0) {
      setChartData({ datasets: [] });
      return;
    }

    // Prepare datasets for chart
    const datasets = selectedFields.map((field, index) => ({
      label: availableFields.find((f) => f.value === field)?.label || field,
      data: dataToUse.map((item) => ({
        x: new Date(item.timestamp),
        y: item[field] || 0,
      })),
      borderColor: `hsl(${index * 120}, 70%, 50%)`,
      backgroundColor: `hsla(${index * 120}, 70%, 50%, 0.2)`,
      tension: 0.1,
    }));

    setChartData({ datasets });
  };

  // Toggle selection of a data field
  const handleFieldToggle = (field: string) => {
    setSelectedFields((prev) => {
      if (prev.includes(field)) {
        return prev.filter((f) => f !== field);
      }
      if (prev.length < 3) {
        return [...prev, field];
      }
      return prev;
    });
  };

  // Apply a predefined time filter (1h, 3h, etc.)
  const applyTimeFilter = (hours: number) => {
    const end = new Date();
    const start = new Date();
    start.setHours(end.getHours() - hours);
    
    setActiveTimeFilter(hours);
    setStartDate(start);
    setEndDate(end);
  };

  // Handle custom date/time selection
  const handleDateChange = (event: React.ChangeEvent<HTMLInputElement>, isStart: boolean) => {
    setActiveTimeFilter(null);
    
    let dateValue = null;
    if (event.target.value) {
      // Parse the input value directly - browser provides value in local time
      dateValue = new Date(event.target.value);
      
      // Validate that it's a valid date
      if (isNaN(dateValue.getTime())) {
        dateValue = null;
      }
    }
    
    if (isStart) {
      setStartDate(dateValue);
      // If end date is before new start date, update end date
      if (dateValue && endDate && dateValue > endDate) {
        setEndDate(dateValue);
      }
    } else {
      setEndDate(dateValue);
    }
  };

  // Format date for input element value
  const formatDateForInput = (date: Date | null): string => {
    if (!date) return '';
    
    // Adjust for timezone to ensure the displayed time matches what was selected
    const tzOffset = date.getTimezoneOffset() * 60000; // offset in milliseconds
    const localDate = new Date(date.getTime() - tzOffset);
    return localDate.toISOString().slice(0, 16);
  };

  // Render the device selection dropdown
  const renderDeviceSelect = () => {
    return (
      <select
        className="select-input"
        value={selectedObjectId}
        onChange={(e) => {
          setSelectedObjectId(e.target.value);
        }}
      >
        <option value="">Select a device...</option>
        {objectOptions.map(option => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    );
  };

  // Open CSV upload modal
  const openUploadModal = () => {
    setUploadModal({
      isOpen: true,
      isLoading: false,
      error: null,
      success: false
    });
  };

  // Close CSV upload modal
  const closeUploadModal = () => {
    setUploadModal({
      isOpen: false,
      isLoading: false,
      error: null,
      success: false
    });
  };

  // Render the upload modal
  const renderUploadModal = () => {
    if (!uploadModal.isOpen) return null;

    return (
      <div className="modal-overlay">
        <div className="modal-container">
          <div className="modal-header">
            <h3>Upload CSV Data</h3>
            <button 
              className="modal-close-button" 
              onClick={closeUploadModal}
              aria-label="Close"
            >
              ×
            </button>
          </div>
          <div className="modal-content">
            {uploadModal.error && (
              <div className="error-message">
                <p>Error: {uploadModal.error}</p>
              </div>
            )}
            
            {uploadModal.success && (
              <div className="success-message">
                <p>CSV file uploaded successfully!</p>
              </div>
            )}
            
            {uploadModal.isLoading ? (
              <div className="modal-loading">
                <div className="loading-spinner"></div>
                <p>Uploading and processing data...</p>
              </div>
            ) : (
              <div className="file-upload-container">
                <p>Select a CSV file containing sensor data to upload:</p>
                <form onSubmit={handleCSVUpload} encType="multipart/form-data">
                  <input 
                    type="file" 
                    accept=".csv" 
                    onChange={handleFileChange}
                    className="file-input"
                    required
                  />
                  <button type="submit" className="upload-submit-button">
                    Upload File
                  </button>
                </form>
                <div className="file-format-info">
                  <p>The CSV file should include these required fields:</p>
                  <ul className="csv-requirements">
                    <li><strong>timestamp</strong> - Format: YYYY-MM-DDThh:mm:ssZ</li>
                    <li><strong>object_id</strong> - Sensor object identifier</li>
                    <li><strong>port_num</strong> - Port number</li>
                    <li><strong>voltage</strong> - Voltage reading</li>
                    <li><strong>current</strong> - Current reading</li>
                  </ul>
                  <a 
                    href="#" 
                    className="sample-csv-link"
                    onClick={downloadSampleCSV}
                  >
                    Download sample CSV template
                  </a>
                </div>
              </div>
            )}
          </div>
          <div className="modal-footer">
            <button 
              className="button" 
              onClick={closeUploadModal}
              disabled={uploadModal.isLoading}
            >
              Close
            </button>
          </div>
        </div>
      </div>
    );
  };

  // Function to generate and download a sample CSV file
  const downloadSampleCSV = (e: React.MouseEvent<HTMLAnchorElement>) => {
    e.preventDefault();
    
    // Create sample data with required fields and example rows
    const csvContent = `timestamp,object_id,port_num,voltage,current,supply_current,supply_volt,voltage_drop,voc,state,controller_error,ai1,ai2,ai3,ai4,ai5,fw_version,q_charge,voltage_set_point,command,target_q,vendor_id,step_number,lite_id,voc_mode,read_error,target_voc,voc_state,voc_exit
2023-09-01T10:00:00Z,1,1,12.5,2.3,1.8,24.0,0.5,14.2,1,0.02,0.5,0.8,1.2,0.3,0.7,v1.2.3,85.6,12.8,1,90.0,VENDOR123,1,LITE001,1,false,14.5,1,0
2023-09-01T10:05:00Z,1,1,12.4,2.4,1.9,24.1,0.6,14.3,1,0.03,0.6,0.9,1.3,0.4,0.8,v1.2.3,86.2,12.8,1,90.0,VENDOR123,2,LITE001,1,false,14.5,1,0
2023-09-01T10:10:00Z,1,1,12.3,2.5,2.0,24.0,0.7,14.1,1,0.02,0.7,1.0,1.4,0.5,0.9,v1.2.3,87.0,12.8,1,90.0,VENDOR123,3,LITE001,1,false,14.5,1,0
2023-09-01T10:15:00Z,1,2,11.8,1.9,1.5,23.8,0.4,13.8,1,0.01,0.4,0.7,1.1,0.2,0.6,v1.2.3,80.2,12.0,1,85.0,VENDOR123,1,LITE002,1,false,14.0,1,0
2023-09-01T10:20:00Z,1,2,11.7,2.0,1.6,23.9,0.5,13.7,1,0.02,0.5,0.8,1.2,0.3,0.7,v1.2.3,81.5,12.0,1,85.0,VENDOR123,2,LITE002,1,false,14.0,1,0`;
    
    // Create a blob from the CSV string
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    
    // Create a temporary link and trigger download
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', 'sample-sensor-data.csv');
    document.body.appendChild(link);
    link.click();
    
    // Clean up
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  };

  // Main render function
  return (
    <div className="app-container">
      <header className="app-header">
        <h1>Sensor Data Analytics Dashboard</h1>
        <p className="subtitle">Monitor and analyze sensor data in real-time</p>
      </header>

      <main className="dashboard-container">
        <section className="control-panel">
          <div className="control-section">
            <h2>Device Selection</h2>
            <div className="form-group">
              <label>Device ID</label>
              <div className="select-wrapper">
                {renderDeviceSelect()}
              </div>
            </div>

            <div className="form-group">
              <label>Port Number</label>
              <div className="select-wrapper">
                <select
                  value={selectedPortNum}
                  onChange={(e) => {
                    setSelectedPortNum(e.target.value);
                  }}
                  className="select-input"
                  disabled={!selectedObjectId || portNums.length === 0}
                >
                  <option value="">Select a port</option>
                  {portNums.map((port) => (
                    <option key={port} value={port}>
                      Port {port}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <div className="upload-csv-container">
              <button 
                className="upload-button" 
                onClick={openUploadModal}
              >
                Upload CSV Data
              </button>
            </div>
          </div>

          <div className="control-section">
            <h2>Time Range</h2>
            <div className="time-filter-buttons">
              {timeFilters.map((filter) => (
                <button
                  key={filter.value}
                  onClick={() => applyTimeFilter(filter.value)}
                  className={`time-filter-button ${activeTimeFilter === filter.value ? 'selected' : ''}`}
                  disabled={!selectedObjectId}
                >
                  {filter.label}
                </button>
              ))}
            </div>
            <div className="date-range-picker">
              <div className="date-picker-container">
                <label>Start Time</label>
                <input
                  type="datetime-local"
                  className="date-picker-input"
                  value={formatDateForInput(startDate)}
                  onChange={(e) => handleDateChange(e, true)}
                  disabled={!selectedObjectId}
                />
                <div className="time-hint">Local time (24h)</div>
              </div>
              <div className="date-picker-container">
                <label>End Time</label>
                <input
                  type="datetime-local"
                  className="date-picker-input"
                  value={formatDateForInput(endDate)}
                  onChange={(e) => handleDateChange(e, false)}
                  min={formatDateForInput(startDate)}
                  disabled={!selectedObjectId}
                />
                <div className="time-hint">Local time (24h)</div>
              </div>
            </div>
          </div>

          <div className="control-section">
            <h2>Data Fields</h2>
            <p className="field-hint">Select up to 3 metrics to display</p>
            <div className="field-buttons">
              {availableFields.map((field) => (
                <button
                  key={field.value}
                  onClick={() => handleFieldToggle(field.value)}
                  className={`field-button ${
                    selectedFields.includes(field.value) ? "selected" : ""
                  }`}
                  disabled={(!selectedFields.includes(field.value) && selectedFields.length >= 3) || !selectedObjectId || allSensorData.length === 0}
                >
                  {field.label}
                </button>
              ))}
            </div>
          </div>

          <div className="control-section">
            <h2>Visualization</h2>
            <p className="field-hint">Select chart type</p>
            <div className="chart-type-selector">
              {chartTypeOptions.map((type) => (
                <button
                  key={type.value}
                  onClick={() => setChartType(type.value as ChartType)}
                  className={`chart-type-button ${
                    chartType === type.value ? "selected" : ""
                  }`}
                  disabled={!chartData || !chartData.datasets.length}
                >
                  {type.label}
                </button>
              ))}
            </div>
          </div>
        </section>

        <section className="visualization-panel">
          {loading ? (
            <div className="loading-container">
              <div className="loading-spinner"></div>
              <p>Loading data...</p>
            </div>
          ) : chartData && chartData.datasets.length > 0 && chartData.datasets[0]?.data.length > 0 ? (
            <DataChart data={chartData} chartType={chartType} />
          ) : chartData && chartData.datasets.length > 0 && chartData.datasets[0]?.data.length === 0 ? (
            <div className="empty-state">
              <p>No data found for the selected time range. Try adjusting the time filter or selecting different dates.</p>
            </div>
          ) : allSensorData.length === 0 && selectedObjectId ? (
            <div className="empty-state">
              <p>No data available for the selected device. The API may not have returned any sensor data.</p>
            </div>
          ) : (
            <div className="empty-state">
              <p>Select a device and data fields to view your sensor visualization. You can optionally filter by port and time range.</p>
            </div>
          )}
        </section>
      </main>
      
      {/* Render the upload modal */}
      {renderUploadModal()}
    </div>
  );
}

export default App;
