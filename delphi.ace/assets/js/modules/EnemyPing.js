import Location from "./Location.js";

export default class EnemyPing{

    #scale;
    #infrastructureSize;
    #depotSize;
    #vehicleSize;

    constructor(iff, data, map){
        var locs = data["Loc"].split(/,/)

        this.#scale = 1.0;
        this.#infrastructureSize = 20*this.#scale;
        this.#depotSize = 15*this.#scale;
        this.#vehicleSize = 15*this.#scale;

        this.timeDiff = {}
        this.map = map
        this.iff = iff
        this.type = data["TypeName"]
        this.loc = {
            lat:    parseFloat(locs[0]),
            long:   parseFloat(locs[1]),
            height: parseFloat(locs[2]),
        }
        this.timeStamp = data["TimeStamp"]
        // console.log("EnemyPing Constructor", map)

        this.cssIcon = L.divIcon({
            // Specify a class name we can refer to in CSS.
            className: 'css-icon',
            html: '<div class="gps_ring"></div>'
            // Set marker width and height
            ,iconSize: [22,22]
            // ,iconAnchor: [11,11]
          });
  
        this.pulser = L.marker([this.loc.lat, this.loc.long], {icon: this.cssIcon})
        this.pulser.addTo(this.map.vehicleLayer);

        this.icon = new L.marker([this.loc.lat,this.loc.long],{icon: L.icon({
            iconUrl: 'assets/img/red.png',
            iconSize:     [28.5, 48.75], // size of the icon
            iconAnchor:   [14.25, 48], // point of the icon which will correspond to marker's location
            popupAnchor:  [0,0] // point from which the popup should open relative to the iconAnchor
        })})
        this.icon.bindPopup("<b style=\"text-align:center\">" + "Enemy IFF: " + this.iff + "</b><br/>" + locs[0] + "," + locs[1] + "<br/>" + this.timeStamp)
        this.icon.addTo(this.map.vehicleLayer)

        this.icon.on('mouseover', function (e) {
            this.openPopup();
        });
        this.icon.on('mouseout', function (e) {
            this.closePopup();
        });

        this.classify()
    }

    update(data){
        var locs = data["Loc"].split(/,/)
        var lastPingTime = new Date(data["TimeStamp"])
        var timeNow = new Date()

        // console.log("EnemyPing update")
        this.loc = {
            lat:    parseFloat(locs[0]),
            long:   parseFloat(locs[1]),
            height: parseFloat(locs[2]),
        }

        
        this.timeDiff = timeNow - lastPingTime;
        this.pulser.setLatLng(L.latLng(this.loc.lat, this.loc.long))
        // this.pulser.bindPopup("<b style=\"text-align:center\">" + this.iff +"<br/>" + this.type + "</b><br/>" + locs[0] + "," + locs[1] + "<br/>Last Seen:\t" + this.msToTime(this.timeDiff))


        this.icon.setLatLng(L.latLng(this.loc.lat, this.loc.long))

        this.icon.bindPopup("<b style=\"text-align:center\">" + "Enemy " + this.type + " (" + this.iff +")" + "</b><br/>" + locs[0] + "," + locs[1] + "<br/>Last Seen:\t" + this.msToTime(this.timeDiff))
        this.colorize();
    }

    msToTime(s) {

        // Pad to 2 or 3 digits, default is 2
        function pad(n, z) {
          z = z || 2;
          return ('00' + n).slice(-z);
        }
      
        var ms = s % 1000;
        s = (s - ms) / 1000;
        var secs = s % 60;
        s = (s - secs) / 60;
        var mins = s % 60;
        var hrs = (s - mins) / 60;
      
        // return pad(hrs) + ':' + pad(mins) + ':' + pad(secs) + '.' + pad(ms, 3);
        return pad(hrs) + ':' + pad(mins) + ':' + pad(secs);

      }
      
    colorize(){
        // console.log("here")
        if (this.timeDiff < 30000){
            console.log("First")
            // 30 seconds
            // console.log("asdfasdf")
            // console.log(this.timeDiff)
            // L.DomUtil.removeClass(this.icon._icon, 'recent')
            // L.DomUtil.removeClass(this.icon._icon, 'old')
            // L.DomUtil.removeClass(this.icon._icon, 'longGone')
            // L.DomUtil.addClass(this.icon._icon, "now")
            // L.DomUtil.setClass(this.icon._icon,"leaflet-marker-icon leaflet-zoom-animated leaflet-interactive now")
            this.pulser._icon.style.filter = "invert(100%) sepia(39%) saturate(3146%) hue-rotate(312deg) brightness(101%) contrast(105%)";
        }
        else if (this.timeDiff < 60000){
            // 1 minute
            // L.DomUtil.removeClass(this.icon._icon, 'now')
            // L.DomUtil.removeClass(this.icon._icon, 'old')
            // L.DomUtil.removeClass(this.icon._icon, 'longGone')
            // L.DomUtil.addClass(this.icon._icon, "recent")    
            // console.log("Here")
            this.pulser._icon.style.filter = "invert(72%) sepia(14%) saturate(1146%) hue-rotate(6deg) brightness(98%) contrast(87%)";
        }
        else if (this.timeDiff < 180000){
        //     // 3 minutes
        //     L.DomUtil.removeClass(this.icon._icon, 'recent')
        //     L.DomUtil.removeClass(this.icon._icon, 'now')
        //     L.DomUtil.removeClass(this.icon._icon, 'longGone')
        //     L.DomUtil.addClass(this.icon._icon, "old")
            console.log("Third")
            this.pulser._icon.style.filter = "invert(36%) sepia(11%) saturate(1727%) hue-rotate(6deg) brightness(93%) contrast(85%)";

        }
        else{
        //     // 5
        //     L.DomUtil.removeClass(this.icon._icon, 'recent')
        //     L.DomUtil.removeClass(this.icon._icon, 'old')
        //     L.DomUtil.removeClass(this.icon._icon, 'now')
        //     L.DomUtil.addClass(this.icon._icon, "longGone")
            this.pulser._icon.style.filter = "invert(0%) sepia(0%) saturate(7430%) hue-rotate(269deg) brightness(97%) contrast(100%)";

        }
        return
    }

    classify(){
        switch(this.type){
            case "ISR":
                this.icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/drone.svg',
                    iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                    iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                    })
                );
            break;

            case "BOMBER":
            this.icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/bomber.svg',
                    iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                    iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                })
            );
            break;

            case "FIGHTER":
                this.icon.setIcon(
                    L.icon({
                        iconUrl: 'assets/svgs/fighter.svg',
                        iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                        iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                    })
                );
            break;

            case "LOGISTICS":
                this.icon.setIcon(
                    L.icon({
                        iconUrl: 'assets/svgs/logisitcs.svg',
                        iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                        iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                    })
                );
            break;

            case "MULTI":
                this.icon.setIcon(
                    L.icon({
                        iconUrl: 'assets/svgs/drone.svg',
                        iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                        iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                    })
                );
            break;

            case "WMD":
                this.icon.setIcon(
                    L.icon({
                        iconUrl: 'assets/svgs/nuclear.svg',
                        iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                        iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                    })
                );
            break;
            
            case "Airbase":
                this.icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/airbase.svg',
                    iconSize:     [this.#infrastructureSize, this.#infrastructureSize],
                    iconAnchor:   [this.#infrastructureSize/2, this.#infrastructureSize/2], 
                    })
                );
                // L.marker(this.latLng(), {icon: L.divIcon({className: 'leaflet-text-labels', html: this.#name}), iconsize: 100, zIndexOffset: 1100}).addTo(window.map.infrastructureLayer);
            break;

            case "Depot":
                this.icon.setIcon(
                    L.icon({
                        iconUrl: 'assets/svgs/supply.svg',
                        iconSize:     [this.#infrastructureSize, this.#infrastructureSize],
                        iconAnchor:   [this.#infrastructureSize/2, this.#infrastructureSize/2], 
                        })
                    );
            break;
    
            case "Fort":
            this.icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/fort.svg',
                    iconSize:     [this.#infrastructureSize, this.#infrastructureSize],
                    iconAnchor:   [this.#infrastructureSize/2, this.#infrastructureSize/2], 
                    })
                );
            break;

            case "ENCAMPMENT":
            this.icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/encampment.svg',
                    iconSize:     [this.#infrastructureSize, this.#infrastructureSize],
                    iconAnchor:   [this.#infrastructureSize/2, this.#infrastructureSize/2], 
                    })
                );
            break;

            case "WMD-PRODUCTION":
            this.icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/nuclear.svg',
                    iconSize:     [this.#depotSize, this.#depotSize],
                    iconAnchor:   [this.#depotSize/2, this.#depotSize/2], 
                    })
                );
            break;

            case "AnimalHospital":
            this.icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/animalhospital.svg',
                    iconSize:     [this.#depotSize, this.#depotSize],
                    iconAnchor:   [this.#depotSize/2, this.#depotSize/2], 
                    })
                );
            break;

            case "Datacenter":
            this.icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/datacenter.svg',
                    iconSize:     [this.#depotSize, this.#depotSize],
                    iconAnchor:   [this.#depotSize/2, this.#depotSize/2], 
                    })
                );
            break;

            case "NETWORK-NODE":
            this.icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/networknode.svg',
                    iconSize:     [this.#depotSize, this.#depotSize],
                    iconAnchor:   [this.#depotSize/2, this.#depotSize/2], 
                    })
                );
            break;

            case "EMBASSY":
            this.icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/embassy.svg',
                    iconSize:     [this.#depotSize, this.#depotSize],
                    iconAnchor:   [this.#depotSize/2, this.#depotSize/2], 
                    })
                );
            break;

            case "Port":
                this.icon.setIcon(
                    L.icon({
                        iconUrl: 'assets/svgs/port.svg',
                        iconSize:     [this.#infrastructureSize, this.#infrastructureSize],
                        iconAnchor:   [this.#infrastructureSize/2, this.#infrastructureSize/2], 
                        })
                    );
               break;


            case "INSTALLATION":
                this.icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/base.svg',
                    iconSize:     [0, 0],
                    iconAnchor:   [0, 0], 
                    })
                );
            break;

            default:
                this.icon.setIcon(
                    L.icon({
                        iconUrl: 'assets/svgs/pin.svg',
                        iconSize:     [this.#depotSize, this.#depotSize],
                        iconAnchor:   [this.#depotSize/2, this.#depotSize/2], 
                    })
                );            
        }
    }
}
