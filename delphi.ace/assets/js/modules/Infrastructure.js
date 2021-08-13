import Location from "./Location.js";

export default class Infrastructure extends Location{

    color;
    #hitpoints;
    #name;
    #team;
    #type;
    #scale;
    #infrastructureSize
    #depotSize;
    #icon;
    #id;
    
    /**
     * 
     * @param {Float} latitude - The Longitude of the Location
     * @param {Float} longitude - The Latitude of the Location
     * @param {Float} altitude - The Altitude of the Location
     */
    constructor(data){
        super(data.latitude, data.longitude, 0);


        // Sets data variable of the class
        this.#name = data.name;
        this.color = data.color;
        this.#hitpoints = data.hitpoints;
        this.#team = data.team;
        this.#type = data.type;
        this.#id = data.id;

        //Sort out icon scaling
        this.#scale = 1.0;
        this.#infrastructureSize = 20*this.#scale;
        this.#depotSize = 15*this.#scale;


        //Create an empty icon in the position it is supposed to be
        this.#icon = L.marker([this.latitude, this.longitude], 
            {icon: 
                L.icon({
                    iconUrl: 'assets/svgs/fort.svg',
                    iconSize:     [0, 0], 
                    iconAnchor:   [0, 0], 
                })
        }).addTo(window.map.infrastructureLayer);

        this.#icon.bindPopup("<b onclick=\"window.ui.showInfrastructure(" + data.id + ");\"  style=\"text-align:center\">" + data.name + "</b><br/>" + this.latitude.toFixed(4) + ", " + this.longitude.toFixed(4) + "<br/>");


        //Sends the icon to be classified
        this.classify();
    }




    /**
     * Check if the infrastructure has been attacked. If it has, alert the user.

     */
    update(data){
        /*Check if the infrastructure has been attacked */
        if(this.#hitpoints > data.hitpoints){
            console.log( (this.#hitpoints - data.hitpoints) + " points of damage taken on " + this.#name);
            window.ui.alert((this.#hitpoints - data.hitpoints) + " points of damage taken on " + this.#name);
        }
    }



    classify(){
        switch(this.#type){
            case "Airbase":
                this.#icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/airbase.svg',
                    iconSize:     [this.#infrastructureSize, this.#infrastructureSize],
                    iconAnchor:   [this.#infrastructureSize/2, this.#infrastructureSize/2], 
                    })
                );
                // L.DomUtil.addClass(this.#icon._icon, "blinking")

                
                L.marker(this.latLng(), {icon: L.divIcon({className: 'leaflet-text-labels', html: this.#name}), iconsize: 100, zIndexOffset: 1100}).addTo(window.map.infrastructureLayer);
            break;

            case "Depot":
                this.#icon.setIcon(
                    L.icon({
                        iconUrl: 'assets/svgs/supply.svg',
                        iconSize:     [this.#infrastructureSize, this.#infrastructureSize],
                        iconAnchor:   [this.#infrastructureSize/2, this.#infrastructureSize/2], 
                        })
                    );
            break;
    
            case "Fort":
            this.#icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/fort.svg',
                    iconSize:     [this.#infrastructureSize, this.#infrastructureSize],
                    iconAnchor:   [this.#infrastructureSize/2, this.#infrastructureSize/2], 
                    })
                );
            break;

            case "ENCAMPMENT":
            this.#icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/encampment.svg',
                    iconSize:     [this.#infrastructureSize, this.#infrastructureSize],
                    iconAnchor:   [this.#infrastructureSize/2, this.#infrastructureSize/2], 
                    })
                );
            break;

            case "WMD-PRODUCTION":
            this.#icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/nuclear.svg',
                    iconSize:     [this.#depotSize, this.#depotSize],
                    iconAnchor:   [this.#depotSize/2, this.#depotSize/2], 
                    })
                );
            break;

            case "AnimalHospital":
            this.#icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/animalhospital.svg',
                    iconSize:     [this.#depotSize, this.#depotSize],
                    iconAnchor:   [this.#depotSize/2, this.#depotSize/2], 
                    })
                );
            break;

            case "Datacenter":
            this.#icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/datacenter.svg',
                    iconSize:     [this.#depotSize, this.#depotSize],
                    iconAnchor:   [this.#depotSize/2, this.#depotSize/2], 
                    })
                );
            break;

            case "NETWORK-NODE":
            this.#icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/networknode.svg',
                    iconSize:     [this.#depotSize, this.#depotSize],
                    iconAnchor:   [this.#depotSize/2, this.#depotSize/2], 
                    })
                );
            break;

            case "EMBASSY":
            this.#icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/embassy.svg',
                    iconSize:     [this.#depotSize, this.#depotSize],
                    iconAnchor:   [this.#depotSize/2, this.#depotSize/2], 
                    })
                );
            break;

            case "Port":
                this.#icon.setIcon(
                    L.icon({
                        iconUrl: 'assets/svgs/port.svg',
                        iconSize:     [this.#infrastructureSize, this.#infrastructureSize],
                        iconAnchor:   [this.#infrastructureSize/2, this.#infrastructureSize/2], 
                        })
                    );
               break;


            case "INSTALLATION":
                this.#icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/base.svg',
                    iconSize:     [0, 0],
                    iconAnchor:   [0, 0], 
                    })
                );
            break;

            default:
                this.#icon.setIcon(
                    L.icon({
                        iconUrl: 'assets/svgs/pin.svg',
                        iconSize:     [this.#depotSize, this.#depotSize],
                        iconAnchor:   [this.#depotSize/2, this.#depotSize/2], 
                    })
                );
        }
        this.teamifyIcon();
    }

    isAlive(){
        return (this.#hitpoints != 0);
    }


    getName(){
        return this.#name;
    }

    getID(){
        return this.#id;
    }

    getTeam(){
        return this.#team;
    }

    getColor(){
        return this.color;
    }

    getType(){
        return this.#type;
    }

    toJSON(){
        return {
            "id" : this.#id,
            "name" : this.#name,
            "color" : this.color,
            "team" : this.#team,
            "latitude" : this.latitude,
            "longitude" : this.longitude,
            "altitude" : this.altitude,
            "type" : this.#type,
            "alive" : (this.#hitpoints != 0)
        }
    }




    /**
     * Private function that takes the private #icon varibale/ object and assigns a team color
     */
    teamifyIcon(){
    if(this.color == "red"){
        L.DomUtil.addClass(this.#icon._icon, 'eastIcon');
    } else
    if(this.color == "blue"){
        L.DomUtil.addClass(this.#icon._icon, 'westIcon');
    } else
    if(this.color == "green"){
        L.DomUtil.addClass(this.#icon._icon, 'centralIcon');
    }else 
    if(this.color == "purple"){
      L.DomUtil.addClass(this.#icon._icon, 'gallifrayIcon');
    }
    }
    

    
}