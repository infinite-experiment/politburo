/**
 * Flight Map Integration with Gleo Library
 * Handles Mercator map initialization, flight path rendering, and interactive features
 */

class FlightMapManager {
  constructor(mapContainerId = 'map') {
    this.mapContainerId = mapContainerId;
    this.map = null;
    this.currentRoute = null;
    this.currentEndpoints = [];
    this.initialized = false;
  }

  /**
   * Initialize the map with gleo library
   */
  async initializeMap() {
    try {
      // Dynamically import gleo modules
      const MercatorMap = (await import('gleo/MercatorMap.mjs')).default;
      const MercatorTiles = (await import('gleo/loaders/MercatorTiles.mjs')).default;

      const mapContainer = document.getElementById(this.mapContainerId);
      if (!mapContainer) {
        console.error(`Map container with id '${this.mapContainerId}' not found`);
        return false;
      }

      // Initialize map centered on the world
      this.map = new MercatorMap(this.mapContainerId, {
        center: [40.7128, -74.0060], // Default to New York
        span: 5e6, // Initial zoom level
        maxSpan: 14e6,
      });

      // Add OSM tiles layer
      const tiles = new MercatorTiles('https://tile.osm.org/{z}/{x}/{y}.png', {
        attribution: "<a href='http://osm.org/copyright'>© OSM contributors</a>"
      });
      tiles.addTo(this.map);

      this.initialized = true;
      console.log('Flight map initialized successfully');
      return true;
    } catch (error) {
      console.error('Failed to initialize flight map:', error);
      return false;
    }
  }

  /**
   * Draw a flight path on the map
   * @param {Array} pathCoords - Array of {lat, long} coordinate objects
   * @param {Object} flightMeta - Flight metadata for display
   */
  async drawFlightPath(pathCoords, flightMeta) {
    if (!this.initialized) {
      const initialized = await this.initializeMap();
      if (!initialized) return;
    }

    try {
      // Import drawing symbols
      const Chain = (await import('gleo/symbols/Chain.mjs')).default;
      const Circle = (await import('gleo/symbols/Circle.mjs')).default;
      const TextLabel = (await import('gleo/symbols/TextLabel.mjs')).default;

      this.resetMap();

      // Convert path coordinates to [lat, lng] format
      const coords = pathCoords.map(p => [
        parseFloat(p.lat || p.latitude),
        parseFloat(p.long || p.longitude)
      ]);

      if (coords.length === 0) {
        console.warn('No path coordinates provided');
        return;
      }

      // Draw the flight route as a green line
      this.currentRoute = new Chain(coords, {
        colour: '#14b8a6', // Teal accent color (matching theme)
        width: 3
      });
      this.currentRoute.addTo(this.map);

      // Add start and end point markers
      const start = coords[0];
      const end = coords[coords.length - 1];

      const startCircle = new Circle(start, {
        size: 10,
        colour: '#10b981' // Green for start
      });

      const endCircle = new Circle(end, {
        size: 10,
        colour: '#ef4444' // Red for end
      });

      const startLabel = new TextLabel(start, {
        text: 'START',
        offset: [10, -10],
        colour: '#10b981'
      });

      const endLabel = new TextLabel(end, {
        text: 'END',
        offset: [10, -10],
        colour: '#ef4444'
      });

      // Add all elements to map and track for cleanup
      [startCircle, endCircle, startLabel, endLabel].forEach(el => {
        el.addTo(this.map);
        this.currentEndpoints.push(el);
      });

      // Calculate bounding box and zoom to fit route
      this.zoomToFitRoute(coords);

      console.log(`Flight path drawn with ${coords.length} waypoints`);
    } catch (error) {
      console.error('Failed to draw flight path:', error);
    }
  }

  /**
   * Calculate and zoom to fit the entire flight route
   * @param {Array} coords - Array of [lat, lng] coordinates
   */
  zoomToFitRoute(coords) {
    if (!coords || coords.length < 2 || !this.map) return;

    // Calculate bounding box
    const lats = coords.map(c => c[0]);
    const lngs = coords.map(c => c[1]);

    const minLat = Math.min(...lats);
    const maxLat = Math.max(...lats);
    const minLng = Math.min(...lngs);
    const maxLng = Math.max(...lngs);

    // Calculate center and span
    const centerLat = (minLat + maxLat) / 2;
    const centerLng = (minLng + maxLng) / 2;

    // Calculate appropriate zoom level based on bounding box
    const latDiff = maxLat - minLat;
    const lngDiff = maxLng - minLng;
    const maxDiff = Math.max(latDiff, lngDiff);

    // Rough conversion from degrees to map units (meters)
    // 1 degree ≈ 111,000 meters at the equator
    const span = Math.max(maxDiff * 111000 * 1.5, 1e6); // Add 50% padding

    try {
      this.map.setCenter([centerLat, centerLng]);
      // Note: If your map library has a method to set span/zoom
      // Update this accordingly. The exact method depends on gleo API
    } catch (error) {
      console.warn('Could not set optimal zoom level:', error);
    }
  }

  /**
   * Clear the map of all drawings
   */
  resetMap() {
    if (this.currentRoute && this.map) {
      try {
        this.map.remove(this.currentRoute);
      } catch (error) {
        console.warn('Error removing route:', error);
      }
      this.currentRoute = null;
    }

    this.currentEndpoints.forEach(endpoint => {
      try {
        if (this.map) {
          this.map.remove(endpoint);
        }
      } catch (error) {
        console.warn('Error removing endpoint:', error);
      }
    });
    this.currentEndpoints = [];
  }

  /**
   * Clean up and destroy the map
   */
  destroy() {
    this.resetMap();
    this.map = null;
    this.initialized = false;
  }
}

// Global instance for easy access
let flightMapManager = null;

/**
 * Initialize the flight map manager on page load
 */
function initializeFlightMap() {
  if (!flightMapManager) {
    flightMapManager = new FlightMapManager('map');
  }
  return flightMapManager;
}

/**
 * Load and display flight details
 * @param {string} flightId - The flight ID to load
 */
async function loadFlightDetails(flightId) {
  try {
    const response = await fetch(`/flights/api/details/${flightId}`);
    if (!response.ok) {
      throw new Error(`Failed to load flight details: ${response.statusText}`);
    }

    const data = await response.json();

    // Update flight details panel
    updateFlightDetailsPanel(data.meta);

    // Initialize map if needed
    if (!flightMapManager || !flightMapManager.initialized) {
      flightMapManager = initializeFlightMap();
      await flightMapManager.initializeMap();
    }

    // Draw the flight path on the map
    if (data.path && data.path.length > 0) {
      await flightMapManager.drawFlightPath(data.path, data.meta);
    }

    // Highlight the selected flight in the list
    document.querySelectorAll('[data-flight-id]').forEach(el => {
      el.classList.remove('ring-2', 'ring-flight-primary');
      if (el.getAttribute('data-flight-id') === flightId) {
        el.classList.add('ring-2', 'ring-flight-primary');
      }
    });

  } catch (error) {
    console.error('Error loading flight details:', error);
    const panel = document.getElementById('flight-details-panel');
    if (panel) {
      panel.innerHTML = `<div class="text-red-500 dark:text-red-400">Error loading flight details: ${error.message}</div>`;
    }
  }
}

/**
 * Update the flight details panel with flight information
 * @param {Object} meta - Flight metadata
 */
function updateFlightDetailsPanel(meta) {
  const panel = document.getElementById('flight-details-panel');
  const contentDiv = panel.querySelector('.flex-1');
  if (!panel || !contentDiv) return;

  const durationMinutes = Math.floor((meta.duration_seconds || 0) / 60);
  const startDate = new Date(meta.started_at || '');
  const formattedDate = startDate.toLocaleString();

  contentDiv.innerHTML = `
    <div class="space-y-4">
      <!-- Header Card -->
      <div class="bg-gradient-to-r from-flight-primary/10 to-flight-accent/10 border-l-4 border-flight-primary rounded-lg p-3">
        <div class="flex items-start justify-between">
          <div>
            <h3 class="text-sm font-bold text-gray-900 dark:text-white">${escapeHtml(meta.call_sign || 'N/A')}</h3>
            <p class="text-xs text-gray-600 dark:text-gray-400 mt-1">${escapeHtml(meta.aircraft || 'N/A')}</p>
          </div>
          <span class="inline-flex items-center rounded-full bg-flight-primary/20 px-3 py-1 text-xs font-medium text-flight-primary">✈️ Active</span>
        </div>
      </div>

      <!-- Route Card -->
      <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 bg-white dark:bg-gray-800">
        <p class="text-xs text-gray-500 dark:text-gray-400 uppercase font-semibold mb-1">Route</p>
        <p class="text-sm font-semibold text-gray-900 dark:text-white">${escapeHtml(meta.route || 'N/A')}</p>
      </div>

      <!-- Stats Grid -->
      <div class="grid grid-cols-2 gap-2">
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-2.5 bg-white dark:bg-gray-800">
          <p class="text-xs text-gray-500 dark:text-gray-400 uppercase font-semibold">Speed</p>
          <p class="text-sm font-bold text-flight-primary mt-1">${(meta.max_speed || 0).toFixed(0)} kts</p>
        </div>
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-2.5 bg-white dark:bg-gray-800">
          <p class="text-xs text-gray-500 dark:text-gray-400 uppercase font-semibold">Altitude</p>
          <p class="text-sm font-bold text-flight-primary mt-1">${(meta.max_alt || 0).toFixed(0)} ft</p>
        </div>
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-2.5 bg-white dark:bg-gray-800">
          <p class="text-xs text-gray-500 dark:text-gray-400 uppercase font-semibold">Duration</p>
          <p class="text-sm font-bold text-flight-primary mt-1">${durationMinutes}m</p>
        </div>
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-2.5 bg-white dark:bg-gray-800">
          <p class="text-xs text-gray-500 dark:text-gray-400 uppercase font-semibold">Landings</p>
          <p class="text-sm font-bold text-flight-primary mt-1">${meta.landings || 0}</p>
        </div>
      </div>

      <!-- Status Badges -->
      <div class="flex flex-wrap gap-2">
        ${meta.violations > 0 ? `<span class="inline-flex items-center rounded-full bg-red-100 dark:bg-red-900/30 px-3 py-1 text-xs font-medium text-red-700 dark:text-red-400">⚠️ ${meta.violations} Violations</span>` : ''}
        ${meta.livery ? `<span class="inline-flex items-center rounded-full bg-blue-100 dark:bg-blue-900/30 px-3 py-1 text-xs font-medium text-blue-700 dark:text-blue-400">🎨 ${escapeHtml(meta.livery)}</span>` : ''}
      </div>
    </div>
  `;
}

/**
 * Escape HTML characters to prevent XSS
 * @param {string} text - Text to escape
 * @returns {string} Escaped text
 */
function escapeHtml(text) {
  if (!text) return '';
  const map = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  };
  return String(text).replace(/[&<>"']/g, m => map[m]);
}

// Expose functions globally for inline script access
window.FlightMapManager = FlightMapManager;
window.flightMapManager = flightMapManager;
window.initializeFlightMap = initializeFlightMap;
window.loadFlightDetails = loadFlightDetails;
window.updateFlightDetailsPanel = updateFlightDetailsPanel;
window.escapeHtml = escapeHtml;

// Auto-initialize on page load
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initializeFlightMap);
} else {
  initializeFlightMap();
}
