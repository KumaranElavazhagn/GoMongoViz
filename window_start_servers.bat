@echo off
echo GoMongoViz Server Startup
echo -------------------------------

REM Install frontend dependencies
cd /d "%~dp0FE"
echo Installing frontend dependencies...
call npm i

REM Start frontend in a new window
echo Starting frontend server...
start cmd /k "title GoMongoViz Frontend && npm start"

REM Navigate to the backend directory
cd /d "%~dp0BE"
echo Starting backend server...
start cmd /k "title GoMongoViz Backend && go run main.go"

echo Both servers are now running.
echo - Frontend: http://localhost:3000
echo - Backend: http://localhost:8080
echo -------------------------------
echo Press any key to close this window. The servers will continue running.
pause