import axios from 'axios';

interface RequestConfig {
  method: string;
  headers: Record<string, string>;
  data?: any;
  [key: string]: any;
}

/**
 * Generic HTTP client function for making API requests
 * @param endpoint The API endpoint to call (without base URL)
 * @param options Request configuration options including method, body, and headers
 * @returns Promise with the API response or error
 */
export async function client(
  endpoint: string,
  { method = 'GET', body, ...customConfig }: { method?: string; body?: any; [key: string]: any } = {}
) {
  // Default headers for all requests
  const headers: Record<string, string> = {
    'Content-Type': "application/json"
  };

  // Clean endpoint to prevent double slashes
  const cleanEndpoint = endpoint.startsWith('/') ? endpoint.substring(1) : endpoint;
  
  // For JSON API calls, we need the full URL (direct fetch endpoints like uploads are handled separately)
  // NOTE: File uploads should NOT use this client - they require special handling
  // File uploads use fetch API directly with FormData and no explicit Content-Type header
  const baseUrl = 'http://localhost:8080';
  const url = `${baseUrl}/${cleanEndpoint}`;

  // Prepare the request configuration
  const requestConfig: RequestConfig = {
    method,
    ...customConfig,
    headers: {
      ...headers,
      ...customConfig.headers
    }
  };
  
  // Add request body if provided
  // This is for JSON data only - file uploads should use FormData instead
  if (body) {
    requestConfig.data = JSON.stringify(body);
  }

  try {
    console.log(`API Request: ${method} ${url}`);
    // Axios automatically handles JSON requests but is not ideal for file uploads
    // For file uploads, use fetch API directly with FormData (see handleCSVUpload in App.tsx)
    const response = await axios(url, requestConfig);
    return response;
  } catch (error: any) {
    console.error(`API Error:`, error);
    // Return error response if available
    if (error?.response) {
      return error.response;
    }
    throw error;
  }
}