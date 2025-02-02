class LeafletMap {
  constructor() {
    this.map = null
  }
  displayMap() {
    var map = L.map("map", { zoomControl: false }).setView([1.300532, 103.780721], 18);
    this.map = map
    L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution:
        'Map data &copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
      maxZoom: 18,
    }).addTo(map);

    L.control.zoom({
      position: 'bottomright'
    }).addTo(map);
  }
  displayShade() {
    /* ShadeMap setup */
    map = this.map
    const loaderEl = document.getElementById('loader');
    let now = new Date();
    const shadeMap = L.shadeMap({
      apiKey: "eyJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6Im9uZ2pvcmRhbjk4QGdtYWlsLmNvbSIsImNyZWF0ZWQiOjE3MzgxNDk4OTYyOTQsImlhdCI6MTczODE0OTg5Nn0.w38etgd8lPutKfmVRAjbkknPsBwfGGbtdnPd5AQPRyM",
      date: now,
      color: '#01112f',
      opacity: 0.7,
      terrainSource: {
        maxZoom: 15,
        tileSize: 256,
        getSourceUrl: ({ x, y, z }) => `https://s3.amazonaws.com/elevation-tiles-prod/terrarium/${z}/${x}/${y}.png`,
        getElevation: ({ r, g, b, a }) => (r * 256 + g + b / 256) - 32768,
        _overzoom: 18,
      },
      getFeatures: async () => {
        try {
          if (map.getZoom() > 15) {
            const bounds = map.getBounds();
            const north = bounds.getNorth();
            const south = bounds.getSouth();
            const east = bounds.getEast();
            const west = bounds.getWest();
            const query = `https://overpass-api.de/api/interpreter?data=%2F*%0AThis%20has%20been%20generated%20by%20the%20overpass-turbo%20wizard.%0AThe%20original%20search%20was%3A%0A%E2%80%9Cbuilding%E2%80%9D%0A*%2F%0A%5Bout%3Ajson%5D%5Btimeout%3A25%5D%3B%0A%2F%2F%20gather%20results%0A%28%0A%20%20%2F%2F%20query%20part%20for%3A%20%E2%80%9Cbuilding%E2%80%9D%0A%20%20way%5B%22building%22%5D%28${south}%2C${west}%2C${north}%2C${east}%29%3B%0A%29%3B%0A%2F%2F%20print%20results%0Aout%20body%3B%0A%3E%3B%0Aout%20skel%20qt%3B`;
            const response = await fetch(query)
            const json = await response.json();
            const geojson = osmtogeojson(json);
            // If no building height, default to one storey of 3 meters
            geojson.features.forEach(feature => {
              if (!feature.properties) {
                feature.properties = {};
              }
              if (!feature.properties.height) {
                feature.properties.height = 3;
              }
            });
            return geojson.features;
          }
        } catch (e) {
          console.error(e);
        }
        return [];
      },
      debug: (msg) => { console.log(new Date().toISOString(), msg) }
    }).addTo(map);

    shadeMap.on('tileloaded', (loadedTiles, totalTiles) => {
      loaderEl.innerText = `Loading: ${(loadedTiles / totalTiles * 100).toFixed(0)}%`;
    });
    /* End ShadeMap setup */


  }
}

let LM = new LeafletMap();

function initLeaflet() {
  LM.displayMap()
  LM.displayShade()
}

function addMarker() {

}

function initLeafletControls(map) {
  /* Controls setup */
  const exposure = document.getElementById('exposure');
  const exposureGradientContainer = document.getElementById('exposure-gradient-container');
  const exposureGradient = document.getElementById('exposure-gradient');

  exposure.addEventListener('click', (e) => {
    const target = e.target;
    if (!target.checked) {
      shadeMap && shadeMap.setSunExposure(false);
      exposureGradientContainer.style.display = 'none';
    } else {
      const { lat, lng } = map.getCenter();
      const { sunrise, sunset } = SunCalc.getTimes(now, lat, lng);
      shadeMap && shadeMap.setSunExposure(true, {
        startDate: sunrise,
        endDate: sunset
      });

      const hours = (sunset - sunrise) / 1000 / 3600;
      const partial = hours - Math.floor(hours);
      const html = [];
      for (let i = 0; i < hours; i++) {
        html.push(`<div>${i + 1}</div>`)
      }
      html.push(`<div style="flex: ${partial}"></div>`);
      exposureGradientContainer.style.display = 'block';
      exposureGradient.innerHTML = html.join('');
    }
  })
  /* End controls setup */
}
