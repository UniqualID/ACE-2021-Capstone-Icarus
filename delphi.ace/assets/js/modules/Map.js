var waypointsList = [];
var distanceList = [];

export default class Map{

    #map;
    #vehicleToggle;
    #infrastructureToggle;
    #iadsToggle;
    #bounds;
    redMarker;
    greenMarker;
    markerLayer;

    /**
     * 
     * @param {float} minZoom - The minimum zoom level for the map
     * @param {float} maxZoom - The maximum zoom level for the map
     * @param {L.latLngBounds} bounds - The bounds of the training area
     */
    constructor(minZoom, maxZoom, bounds){
        
        this.redMarker = L.icon({
            iconUrl: 'assets/img/red.png',
            iconSize:     [28.5, 48.75], // size of the icon
            iconAnchor:   [14.25, 48], // point of the icon which will correspond to marker's location
            popupAnchor:  [0,0] // point from which the popup should open relative to the iconAnchor
        });
        this.greenMarker = L.icon({
            iconUrl: '/assets/img/green.png',
            iconSize:     [28.5, 48.75], // size of the icon
            iconAnchor:   [14.25, 48], // point of the icon which will correspond to marker's location
            popupAnchor:  [0,0] // point from which the popup should open relative to the iconAnchor
        });
        this.blackMarker = L.icon({
            iconUrl: 'assets/svgs/pin.svg',
            iconSize:     [20, 20], // size of the icon
            iconAnchor:   [10, 20], // point of the icon which will correspond to marker's location
        });

        this.targetMarker = null

        console.log("Initialising Leaflet Map for Training Area");
        this.#map = L.map('mapid', {zoomSnap: 0.25}).fitBounds(bounds);;
        L.tileLayer('http://maps.ace/tile/{z}/{x}/{y}.png', {minZoom: minZoom, maxZoom: maxZoom}).addTo(this.#map);

        this.vehicleLayer = L.layerGroup().addTo(this.#map);
        this.infrastructureLayer = L.layerGroup().addTo(this.#map);
        this.iadsLayer = L.layerGroup().addTo(this.#map);
        this.rectangleLayer = L.layerGroup().addTo(this.#map);
        this.bordersLayer = L.layerGroup().addTo(this.#map);

        // this.markerLayer =  new L.LayerGroup();
        this.markerLayer = new L.layerGroup();
        this.markerLayer.addTo(this.#map);
        

        this.boxStart = null;
        this.boxEnd = null;
        this.dragtangle = null;
        this.#bounds = bounds;




        //Terratory Labels
        {
            let valinar = L.marker([50.6807,-68.4085], {icon: L.divIcon({className: 'leaflet-map-labels', html: "VALINAR"}), zIndexOffset: -80}).addTo(this.#map);
            let gallifrey = L.marker([47.8635,-67.0719], {icon: L.divIcon({className: 'leaflet-map-labels', html: "GALLIFREY"}), zIndexOffset: -80}).addTo(this.#map);
            let halcyon = L.marker([49.6154,-62.8057], {icon: L.divIcon({className: 'leaflet-map-labels', html: "HALCYON"}), zIndexOffset: -80}).addTo(this.#map);
            let malazan = L.marker([45.6445,-64.2123], {icon: L.divIcon({className: 'leaflet-map-labels', html: "MALAZAN"}), zIndexOffset: -80}).addTo(this.#map);
        }

        // toggle variables
        this.#vehicleToggle = false;
        this.#infrastructureToggle = false;
        this.#iadsToggle = false;


        $.getJSON( "assets/borders/borders.geojson", function( borders ) {
        if (borders!=""){
            L.geoJson(borders, {
                style: function(feature) {
                    switch (feature.properties.name) {
                        case 'Malazan': return {color: "#385D8A", fillOpacity: 0};
                        case 'Gallifrey':   return {color: "#E46C0A", fillOpacity: 0};
                        case 'Halcyon': return {color: "#FCD016", fillOpacity: 0};
                        case 'Valinar': return {color: "#BE4B48", fillOpacity: 0};                  
                    }
                }
            }).addTo(this.bordersLayer);
            {
                L.rectangle(this.#bounds, {color: "#FF0000", fill:false, weight: 2}).addTo(this.bordersLayer);
            }
        }}.bind(this)).fail(function(){
            window.ui.alert("Delphi cannot get nation borders.")
        });

        /**
         * The Dragtangle
         */
        this.dragtangle = null;

        this.#map.on('mousedown', function(e) {
            // if (event.ctrlKey) {
            //     this.#map.dragging.disable();
            //     this.boxStart = e.latlng;
            //     this.boxEnd = e.latlng;
            //     this.#map.removeLayer(this.rectangleLayer);
            //     this.rectangleLayer = L.layerGroup().addTo(this.#map);
            //     this.dragtangle = L.rectangle([this.boxStart, this.boxEnd], {color: "#ff7800", weight: 1}).addTo(this.rectangleLayer);
            // }
            // if (event.altKey){
            //     this.boxStart = null;
            //     this.boxEnd = null;
            //     this.#map.removeLayer(this.rectangleLayer);
            //     //hideSidebar
            // }
        }.bind(this));
     
        this.#map.on('mouseup', function(e) {
        //     if (event.ctrlKey) {
        //         this.boxEnd = e.latlng;
        //         this.#map.dragging.enable();
        //         //this.#map.removeLayer(this.rectangleLayer);
        //         window.ui.showInBounds();
        //    }
         }.bind(this));

        //Mouse hover coords
        this.#map.on('mousemove', function(e){
            $("#mousePos").text(e.latlng.lat.toFixed(4) + "," + e.latlng.lng.toFixed(4));
            if ((event.ctrlKey) && (this.dragtangle !== null)) {
                this.boxEnd = e.latlng;
                this.dragtangle.setBounds([this.boxStart, this.boxEnd]);
              }
        }.bind(this));


        /*
         *
         *
         * 
         * 
         *                                                                       
                                           ,----,                                                                                                 
                     ,--,                ,/   .`|                         ,----..                                                                 
  .--.--.          ,--.'|    ,---,     ,`   .'  :           ,----..      /   /   \       ,---,         ,---,.                            ,----,   
 /  /    '.     ,--,  | : ,`--.' |   ;    ;     /          /   /   \    /   .     :    .'  .' `\     ,'  .' |                ,---.     .'   .' \  
|  :  /`. /  ,---.'|  : ' |   :  : .'___,/    ,'          |   :     :  .   /   ;.  \ ,---.'     \  ,---.'   |               /__./|   ,----,'    | 
;  |  |--`   |   | : _' | :   |  ' |    :     |           .   |  ;. / .   ;   /  ` ; |   |  .`\  | |   |   .'          ,---.;  ; |   |    :  .  ; 
|  :  ;_     :   : |.'  | |   :  | ;    |.';  ;           .   ; /--`  ;   |  ; \ ; | :   : |  '  | :   :  |-,         /___/ \  | |   ;    |.'  /  
 \  \    `.  |   ' '  ; : '   '  ; `----'  |  |           ;   | ;     |   :  | ; | ' |   ' '  ;  : :   |  ;/|         \   ;  \ ' |   `----'/  ;   
  `----.   \ '   |  .'. | |   |  |     '   :  ;           |   : |     .   |  ' ' ' : '   | ;  .  | |   :   .'          \   \  \: |     /  ;  /    
  __ \  \  | |   | :  | ' '   :  ;     |   |  '           .   | '___  '   ;  \; /  | |   | :  |  ' |   |  |-,           ;   \  ' .    ;  /  /-,   
 /  /`--'  / '   : |  : ; |   |  '     '   :  |           '   ; : .'|  \   \  ',  /  '   : | /  ;  '   :  ;/|            \   \   '   /  /  /.`|   
'--'.     /  |   | '  ,/  '   :  |     ;   |.'            '   | '/  :   ;   :    /   |   | '` ,/   |   |    \             \   `  ; ./__;      :   
  `--'---'   ;   : ;--'   ;   |.'      '---'              |   :    /     \   \ .'    ;   :  .'     |   :   .'              :   \ | |   :    .'    
             |   ,/       '---'                            \   \ .'       `---`      |   ,.'       |   | ,'                 '---"  ;   | .'       
             '---'                                          `---`                    '---'         `----'                          `---'          
                                                                                                                                                           * 
         * 
         * 
         * 
         * 
         */

        this.linesToDelete = []
        this.lastMarker = null


        // var self = this
        // $("#submit").click(function() {
        //   var input = $('#textbox').val().split(/\n/);
        //   console.log(typeof(input))
        //   var coordinates = []
        //   $.each( input, function( index, value ) {
        //     /* console.log(index,value) */
        //     var temp = value.split(/[ ,]+/).filter(Boolean)
        //     coordinates.push([parseFloat(temp[0]), parseFloat(temp[1])])
            
        //   });
        //   // console.log(coordinates)
          
        //   $.each( coordinates, function(index,value){
        //     console.log(value, self.markerLayer)
        //     var m = new L.marker(value, {icon:greenMarker})
        //     m.bindPopup("" + value[0].toFixed(4) + "\t" + value[1].toFixed(4));
        //     m.addTo(self.markerLayer)
        //   })
        // });
      

        this.#map.on('click', function(e) {

            let plotToggle = document.getElementById("plot-active");
            let toggled = plotToggle.classList.contains("active");

            if (toggled) {
                console.log("plot toggle active");
      
                // targeting marker
                if (e.originalEvent.ctrlKey){
                    var lat = e.latlng.lat;
                    var lon = e.latlng.lng;    
                    if(this.targetMarker != null){
                        console.log(this.targetMarker)
                        this.targetMarker.remove();
                    }
                    this.targetMarker = new L.marker([lat,lon], {icon:this.redMarker})
                    this.targetMarker.addTo(this.markerLayer);
                }
            
            
                //ALL OTHER MARKERS ()
                else{
            
                    var mar = new L.marker(e.latlng, {icon:this.blackMarker});
                    mar.bindPopup("" + e.latlng.lat.toFixed(4) + "\t" + e.latlng.lng.toFixed(4));      
                    mar.addTo(this.markerLayer);
            
                    
                    if(this.lastMarker !== null){
                        var pointList = [mar.getLatLng(), this.lastMarker];
                        var firstpolyline = new L.Polyline(pointList, {
                            color: 'red',
                            weight: 3,
                            opacity: 0.5,
                            smoothFactor: 1
                        });  
                    
                        /* For UI waylist function */
                        var markerLatLng = mar.getLatLng();
                        var distance = this.lastMarker.distanceTo(markerLatLng).toFixed(0);
                        console.log(this.lastMarker, mar.getLatLng());
                        console.log(distance);
                        waypointsList.push(markerLatLng);
                        distanceList.push(distance);
                        console.log(waypointsList);

                        ui.listWaypoints(waypointsList, distanceList);
                        
                
                        firstpolyline.bindPopup('Distance: ' + this.lastMarker.distanceTo(mar.getLatLng()).toFixed(4));
                        firstpolyline.on('mouseover', function (e) {
                            this.openPopup();
                        });
                        firstpolyline.on('mouseout', function (e) {
                            this.closePopup();
                        });
                
                        firstpolyline.addTo(this.markerLayer);
                        this.linesToDelete.push(firstpolyline)
                        this.lastMarker = mar.getLatLng()
                    } else {
                    // console.log('lines')
                    this.lastMarker = mar.getLatLng()
                    }
                }
            }
        }.bind(this));
      
        // this.addEnemyInfrastructure()
    }

    /**
     * 
     * @returns Leaflet Map Object
     */
    object(){
        return this.#map;
    }

    val_plotPoints(raw){
        // console.log(raw)
        var input = raw.split(/\n/)
        // console.log(input)        
        
        var self = this;
        var coordinates = []
        $.each( input, function( index, value ) {
        /* console.log(index,value) */
            if(value != ""){
                var temp = value.split(/[ ,]+/).filter(Boolean)
                coordinates.push([parseFloat(temp[0]), parseFloat(temp[1])])
            }
        });
        // console.log(coordinates)
        
        $.each( coordinates, function(index,value){
            // console.log(value, self.markerLayer)
            var m = new L.marker(value, {icon:self.greenMarker})
            m.bindPopup("" + value[0].toFixed(4) + "\t" + value[1].toFixed(4));
            // console.log(self.markerLayer, "is this undefined?")
            m.addTo(self.markerLayer)
        })
    }
    /**
     * Centers the map to the orignal Map bounds
     */
    center(){
        this.#map.fitBounds(this.#bounds);
    }

    toggleVehicles() {
        if(!this.#vehicleToggle) {
            this.#map.removeLayer(this.vehicleLayer);
        } else {
            this.#map.addLayer(this.vehicleLayer);
            for(let vehicle in window.delphiMap.vehicles){
                window.delphiMap.vehicles[vehicle].teamifyIcon();
            }
        }
        this.#vehicleToggle = !this.#vehicleToggle;
    }

    toggleInfrastructure() {
        if(!this.#infrastructureToggle) {
            this.#map.removeLayer(this.infrastructureLayer);
        } else {
            this.#map.addLayer(this.infrastructureLayer);
            for(let infrastructure in window.delphiMap.infrastructure){
                window.delphiMap.infrastructure[infrastructure].teamifyIcon();
            }
        }
        this.#infrastructureToggle = !this.#infrastructureToggle;
    }

    toggleIads() {
        if(!this.#iadsToggle) {
            this.#map.removeLayer(this.iadsLayer);
        } else {
            this.#map.addLayer(this.iadsLayer);
            for(let iad in window.delphiMap.iads){
                window.delphiMap.iads[iad].teamifyIcon();
            }
        }
        this.#iadsToggle = !this.#iadsToggle;
    }

    removeWaypoints() {
        waypointsList = [];
        distanceList = [];
    }

    addEnemyInfrastructure(){

        // aribase 1

        var mar1 = L.marker([49.8368, -64.2868], {icon:this.greenMarker});

        mar1.bindPopup("" + 49.8368 + "\t" + -64.2868);      

        mar1.addTo(this.markerLayer);



        // airbase 2

        var mar2 = L.marker([49.3325, -62.902], {icon:this.greenMarker});

        mar2.bindPopup("" + 49.3325 + "\t" + -62.902);      

        mar2.addTo(this.markerLayer); 



        // depot 1

        var mar3 = L.marker([49.6605, -63.6695], {icon:this.greenMarker});

        mar3.bindPopup("" + 49.6605 + "\t" + -63.6695);      

        mar3.addTo(this.markerLayer); 



        // depot 2

        var mar4 = L.marker([49.3175, -62.9425], {icon:this.greenMarker});

        mar4.bindPopup("" + 49.3175 + "\t" + -62.9425);      

        mar4.addTo(this.markerLayer); 



        // depot 3

        var mar5 = L.marker([49.1418, -62.2238], {icon:this.greenMarker});

        mar5.bindPopup("" + 49.1418 + "\t" + -62.2238);      

        mar5.addTo(this.markerLayer); 



        // datacentre

        var mar6 = L.marker([49.476, -63.5111], {icon:this.greenMarker});

        mar6.bindPopup("" + 49.476 + "\t" + -63.5111);      

        mar6.addTo(this.markerLayer); 



        // fort 1

        var mar7 = L.marker([49.7725, -63.8525], {icon:this.greenMarker});

        mar7.bindPopup("" + 49.7725 + "\t" + -63.8525);      

        mar7.addTo(this.markerLayer); 



        // fort 2

        var mar8 = L.marker([49.559, -63.104], {icon:this.greenMarker});

        mar8.bindPopup("" + 49.559 + "\t" + -63.104);      

        mar8.addTo(this.markerLayer); 



        // fort 3

        var mar9 = L.marker([49.406, -64.2868], {icon:this.greenMarker});

        mar9.bindPopup("" + 49.406 + "\t" + -64.2868);      

        mar9.addTo(this.markerLayer); 



        // net node 1

        var mar10 = L.marker([49.6, -63.4142], {icon:this.greenMarker});

        mar10.bindPopup("" + 49.6 + "\t" + -63.4142);      

        mar10.addTo(this.markerLayer); 



        // net node 2

        var mar11 = L.marker([49.2675, -62.4995], {icon:this.greenMarker});

        mar11.bindPopup("" + 49.2675 + "\t" + -62.4995);      

        mar11.addTo(this.markerLayer);

        

        // installation 1

        var mar12 = L.marker([49.73, -64.037], {icon:this.greenMarker});

        mar12.bindPopup("" + 49.73 + "\t" + -64.037);      

        mar12.addTo(this.markerLayer); 



        // installation 2

        var mar13 = L.marker([49.66, -63.555], {icon:this.greenMarker});

        mar13.bindPopup("" + 49.66 + "\t" + -63.555);      

        mar13.addTo(this.markerLayer); 

    }
}