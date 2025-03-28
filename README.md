# GoMongoViz - Sensor Data Visualization Dashboard

GoMongoViz is a full-stack sensor data visualization dashboard that allows users to view, analyze, and upload sensor data stored in MongoDB. This application provides real-time monitoring of sensor metrics such as voltage, current, and other measurements with interactive charts and filtering capabilities.

![Dashboard Screenshot](path-to-screenshot.png)

## Features

- **Data Visualization**: Interactive charts with zoom and pan functionality
- **Real-time Monitoring**: View and analyze sensor data in real-time
- **Multiple Upload Formats**: Upload sensor data via CSV or JSON files
- **Filtering Capabilities**: Filter data by device, port, and time range
- **Multi-metric Analysis**: Compare up to 3 metrics simultaneously
- **Responsive Design**: Works on desktop and mobile devices

## Tech Stack

### Frontend
- React (TypeScript)
- Chart.js for data visualization
- CSS for styling
- Axios for API communication

### Backend
- Go (Golang)
- Gorilla Mux for routing
- MongoDB for data storage
- Clean Architecture pattern

## Project Structure

The project follows a clean architecture pattern with clear separation of concerns:

```
GoMongoViz/
├── BE/                  # Backend (Go)
│   ├── database/        # Database connection
│   ├── handlers/        # HTTP request handlers
│   ├── model/           # Data models
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic
│   └── main.go          # Entry point
├── FE/                  # Frontend (React)
│   ├── public/          # Static files
│   ├── src/             # Source code
│   │   ├── api/         # API clients
│   │   ├── components/  # React components
│   │   ├── types/       # TypeScript types
│   │   └── utils/       # Utilities
│   └── package.json     # Dependencies
└── start_servers.bat    # Startup script
```

## Getting Started

### Prerequisites

- Node.js (v14+)
- Go (v1.16+)
- MongoDB (local or Atlas)

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/KumaranElavazhagn/GoMongoViz.git
   cd GoMongoViz
   ```

2. Using the startup script (Windows):
   - Simply run `start_servers.bat` to install dependencies and start both servers
   - This will open two command windows: one for the frontend and one for the backend

3. Manual installation:

   **Frontend:**
   ```
   cd FE
   npm install
   npm start
   ```

   **Backend:**
   ```
   cd BE
   go run main.go
   ```

4. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## API Endpoints

- `GET /api/objects` - Get all unique object IDs
- `GET /api/ports/{objectId}` - Get ports for a specific object
- `GET /api/data/{objectId}?port_num={portNum}` - Get data for a specific object and port
- `POST /api/upload` - Upload and process CSV data
- `POST /api/upload-json` - Upload and process JSON data

## Data Upload Formats

### CSV Upload Format

When uploading CSV files, ensure they follow this format:

**Required fields:**
- `timestamp` - Format: YYYY-MM-DDThh:mm:ssZ (RFC3339)
- `object_id` - Sensor object identifier
- `port_num` - Port number
- `voltage` - Voltage reading
- `current` - Current reading
- `supply_current` - Supply current value
- `supply_volt` - Supply voltage
- `voltage_drop` - Voltage drop measurement
- `voc` - Voltage Open Circuit value

**Optional fields include:**
- `state`, `controller_error`, `ai1`-`ai5`, etc.

You can download a sample CSV template from the upload modal.

### JSON Upload Format

When uploading JSON files, they should contain an array of objects with the following structure:

**Required fields:**
- `timestamp` - Format: "2023-09-01T10:00:00Z" (RFC3339 string)
- `object_id` - Numeric sensor object identifier
- `port_num` - Numeric port number

**Optional fields include:**
- `voltage`, `current`, `supply_current`, `supply_volt`, `voltage_drop`, `voc`, etc.

Example JSON format:
```json
[
  {
    "timestamp": "2023-09-01T10:00:00Z",
    "object_id": 1,
    "port_num": 1,
    "voltage": 12.5,
    "current": 2.3,
    "supply_current": 1.8,
    "supply_volt": 24.0
  },
  {
    "timestamp": "2023-09-01T10:05:00Z",
    "object_id": 1,
    "port_num": 1,
    "voltage": 12.4,
    "current": 2.4
  }
]
```

You can download a sample JSON template from the upload modal.

## File Upload Implementation

The application supports two types of file uploads:

1. **CSV Upload**: Implemented using `FormData` and multipart/form-data encoding
2. **JSON Upload**: Implemented using direct JSON posting with application/json content type

### Key File Upload Implementation Details

1. **Frontend Implementation**
   - CSV uploads use native `FormData` API for building multipart requests
   - JSON uploads use `FileReader` to read file contents and send as JSON body
   - Both methods include proper error handling and user feedback
   - The interface provides a split-button approach for selecting upload type

2. **Backend Implementation**
   - CSV uploads are handled with Go's `multipart` package and CSV parser
   - JSON uploads are processed using the standard JSON decoder
   - Both methods validate data and provide detailed error messages
   - Common data storage logic is reused between both upload types

3. **Common Upload Issues**
   - **Content-Type issues**: Each upload method requires different Content-Type headers
   - **CORS issues**: Preflight requests must be properly handled for file uploads
   - **File size limits**: The backend limits uploads to 10MB to prevent abuse

4. **Debugging File Uploads**
   - Backend logs include detailed information about incoming requests
   - Frontend console logs show file information before upload
   - Error responses include detailed information about what went wrong

## Development Notes

- Frontend: The React application uses a proxy configuration to simplify API calls
- Backend: CORS is configured to allow cross-origin requests for development

## MongoDB Configuration

The MongoDB connection is configured in `BE/database/db.go`. Update the credentials to match your MongoDB instance.

## Acknowledgements

- [Chart.js](https://www.chartjs.org/)
- [React](https://reactjs.org/)
- [Go](https://golang.org/)
- [MongoDB](https://www.mongodb.com/)