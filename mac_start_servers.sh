#!/bin/bash

echo "GoMongoViz Server Startup"
echo "-------------------------------"

# Install frontend dependencies
cd "$(dirname "$0")/FE"
echo "Installing frontend dependencies..."
npm i

# Start frontend in a new terminal window
echo "Starting frontend server..."
osascript -e 'tell app "Terminal" to do script "cd '"$(pwd)"' && npm start"'

# Navigate to the backend directory
cd "../BE"
echo "Starting backend server..."
osascript -e 'tell app "Terminal" to do script "cd '"$(pwd)"' && go run main.go"'

echo "Both servers are now running."
echo "- Frontend: http://localhost:3000"
echo "- Backend: http://localhost:8080"
echo "-------------------------------"
echo "You can close this terminal. The servers will continue running."