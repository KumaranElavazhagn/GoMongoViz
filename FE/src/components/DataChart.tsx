import React from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
  Filler,
} from 'chart.js';
import 'chartjs-adapter-date-fns';
import zoomPlugin from 'chartjs-plugin-zoom';
import { Line, Bar, Scatter } from 'react-chartjs-2';
import { ChartData } from '../types';

// Register ChartJS components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
  Filler,
  zoomPlugin
);

type ChartType = "line" | "bar" | "area" | "scatter";

interface DataChartProps {
  data: ChartData;
  chartType?: ChartType;
  title?: string;
}

export const DataChart: React.FC<DataChartProps> = ({ data, chartType = "line", title }) => {
  // Prepare data for line, area, and scatter charts with appropriate styling
  const prepareChartData = () => {
    if (chartType === 'area') {
      return {
        datasets: data.datasets.map((dataset, index) => ({
          ...dataset,
          fill: true, // Fill area under the line
          backgroundColor: `hsla(${index * 120}, 70%, 50%, 0.15)`,
          borderWidth: 2
        }))
      };
    }
    
    if (chartType === 'scatter') {
      return {
        datasets: data.datasets.map((dataset, index) => ({
          ...dataset,
          pointBackgroundColor: `hsl(${index * 120}, 70%, 50%)`,
          pointBorderColor: `hsl(${index * 120}, 70%, 40%)`,
          pointRadius: 5,
          pointHoverRadius: 7
        }))
      };
    }
    
    // Default line chart
    return {
      datasets: data.datasets.map(dataset => ({
        ...dataset,
        borderWidth: 2,
        pointRadius: 3
      }))
    };
  };

  // Special data preparation for bar charts which require a different format
  const prepareBarChartData = () => {
    // Extract labels (timestamps) and datasets in the format expected by Chart.js for bar charts
    const timestamps = data.datasets[0]?.data.map(point => 
      point.x.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})
    ) || [];
    
    return {
      labels: timestamps,
      datasets: data.datasets.map((dataset, index) => ({
        label: dataset.label,
        data: dataset.data.map(point => point.y), // Bar charts need plain y values
        backgroundColor: `hsla(${index * 120}, 70%, 50%, 0.7)`,
        borderColor: `hsl(${index * 120}, 70%, 40%)`,
        borderWidth: 1,
        borderRadius: 3,
      }))
    };
  };

  // Charts options for line, area, and scatter with time-based x-axis
  const options = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      x: {
        type: 'time' as const,
        time: {
          unit: 'minute' as const,
          displayFormats: {
            minute: 'HH:mm',
            hour: 'HH:mm',
            day: 'MMM d',
            week: 'MMM d',
            month: 'MMM yyyy'
          }
        },
        grid: {
          color: 'rgba(0, 0, 0, 0.05)',
        },
      },
      y: {
        grid: {
          color: 'rgba(0, 0, 0, 0.05)',
        },
        beginAtZero: true
      }
    },
    plugins: {
      legend: {
        position: 'top' as const,
      },
      title: {
        display: !!title,
        text: title || (chartType === 'area' ? 'Area Chart' : chartType === 'scatter' ? 'Scatter Plot' : 'Line Chart')
      },
      zoom: {
        pan: {
          enabled: true,
          mode: 'xy' as const,
        },
        zoom: {
          wheel: {
            enabled: true,
          },
          mode: 'xy' as const,
        }
      }
    }
  };

  // Bar chart options with category axis instead of time axis
  const barOptions = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      x: {
        grid: {
          color: 'rgba(0, 0, 0, 0.05)',
        },
      },
      y: {
        beginAtZero: true,
        grid: {
          color: 'rgba(0, 0, 0, 0.05)',
        }
      }
    },
    plugins: {
      legend: {
        position: 'top' as const,
      },
      title: {
        display: !!title,
        text: title || 'Bar Chart'
      },
      zoom: {
        pan: {
          enabled: true,
          mode: 'xy' as const,
        },
        zoom: {
          wheel: {
            enabled: true,
          },
          mode: 'xy' as const,
        }
      }
    }
  };

  // Render the appropriate chart component based on chart type
  const renderChart = () => {
    switch (chartType) {
      case 'bar':
        // For bar charts, use the special data format
        return <Bar data={prepareBarChartData()} options={barOptions} />;
      case 'scatter':
        return <Scatter data={prepareChartData()} options={options} />;
      case 'area':
      case 'line':
      default:
        return <Line data={prepareChartData()} options={options} />;
    }
  };

  return (
    <div className="chart-wrapper">
      <div className="chart-container">
        {renderChart()}
      </div>
      <div className="chart-instructions">
        <p>Hold <kbd>Ctrl</kbd> + scroll to zoom • <kbd>Ctrl</kbd> + drag to pan</p>
      </div>
    </div>
  );
};