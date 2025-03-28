/**
 * Represents a selectable data field with display label and internal value
 */
export interface Field {
  label: string;
  value: string;
}

/**
 * Represents a single data point on a chart with x-y coordinates
 * x is typically a timestamp, y is the measurement value
 */
export interface DataPoint {
  x: Date;
  y: number;
}

/**
 * Represents styling and data for a single chart dataset
 */
export interface ChartDataset {
  label: string;
  data: DataPoint[];
  borderColor: string;
  backgroundColor?: string;
  tension: number;
  fill?: boolean;
}

/**
 * Collection of datasets to be displayed on a chart
 */
export interface ChartData {
  datasets: ChartDataset[];
}

/**
 * Generic API response type with timestamp and additional properties
 */
export interface ApiResponse {
  timestamp: string;
  [key: string]: any;
} 