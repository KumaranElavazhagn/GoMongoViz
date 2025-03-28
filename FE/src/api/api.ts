import { client } from "./Client";
import { ApiResponse } from "../types";

/**
 * Fetches available devices from the backend
 * @returns List of device objects with their IDs
 */
export const GetObject = async () => {
  try {
    const response = await client(`api/objects`, {
      method: "GET",
    });
    
    if (response && response.status === 200 && response.data) {
      return response;
    } else {
      console.error("Failed to get objects:", response);
      return { data: [] };
    }
  } catch (error) {
    console.error("Error in GetObject:", error);
    return { data: [] };
  }
};

/**
 * Fetches ports available for a specific device
 * @param objectId The ID of the selected device
 * @returns List of ports available for the device
 */
export const GetPort = async (objectId: string) => {
  try {
    const response = await client(`api/ports/${objectId}`, {
      method: "GET",
    });
    
    if (response && response.status === 200 && response.data) {
      return response;
    } else {
      console.error("Failed to get ports:", response);
      return { data: [] };
    }
  } catch (error) {
    console.error("Error in GetPort:", error);
    return { data: [] };
  }
};

/**
 * Fetches sensor data for a specific device and port
 * @param objectId The ID of the selected device
 * @param portNum The port number or "all" for all ports
 * @param timeRange Time range filter (unused currently)
 * @returns Sensor data for the specified device and port
 */
export const GetData = async (objectId: string, portNum: string) => {
  try {
    // If portNum is "all", get data from all ports
    const endpoint = portNum === "all" 
      ? `api/data/${objectId}`
      : `api/data/${objectId}?port_num=${portNum}`;
    
    const response = await client(endpoint, {
      method: "GET",
    });
    
    if (response && response.status === 200 && response.data) {
      return response;
    } else {
      console.error("Failed to get data:", response);
      return { data: { SensorData: [] } };
    }
  } catch (error) {
    console.error("Error in GetData:", error);
    return { data: { SensorData: [] } };
  }
};