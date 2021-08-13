import Location from './Location.js';

/**
 * Vehicle class that contains logic for Vehicle Operations
 * Extends the Location object
 */
export default class Vehicle extends Location{
    #id;
    #icon;
    #name;
    #role;
    #show;
    #alive;
    #showing;
    #team;
    #type;
    #scale;
    #vehicleSize;
    
    /**
     * Creates a Vehicle and sets its main variables
     * @param {Object} data - Individual Vehicle object taken from the JSON file
     */
    constructor(data){
        super(data.latitude, data.longitude, data.altitude);


        // Sets data variable of the class
        this.#name = data.name;
        this.#alive = data.alive;
        this.color = data.color;
        this.#team = data.team;
        this.#role = data.role;
        this.#type = data.type;
        this.#id = data.id;

        this.#scale = 1.0;
        this.#vehicleSize = 15*this.#scale;

        this.#showing = false;

        this.#icon = L.marker([this.latitude, this.longitude], 
            {icon: 
                L.icon({
                    iconUrl: 'assets/svgs/drone.svg',
                    iconSize:     [0, 0], 
                    iconAnchor:   [0, 0], 
                })
        }).addTo(window.map.vehicleLayer);


        this.#icon.bindPopup("<b onclick=\"window.ui.showVehicle(" + data.id + ");\" style=\"text-align:center\">" + data.name + "</b><br/>" + this.latitude.toFixed(4) + ", " + this.longitude.toFixed(4) + "<br/>" + this.altitude.toFixed(1) + "<br/>");
            
        //This logic defines whether flying vehicles show if they are on the ground
        if ( (this.#type == "copter" ) && (!(this.altitude > 0)) ){
            this.#show = false;
        } else{
            this.#show = true;
        }
        
        //Sends the icon to be classified
        this.classify();
        

    }

    /**
     * Move the vehicle on the map
     * @param {Object} data - Individual Vehicle object taken from the JSON file
     */
    update(data){

        /* Run a check to see if a vehicle has taken off and must be shown */
        if((this.altitude == 0) && (data.altitude !== 0 )){
            this.#show = true;
        }

        /* Run a check to see if a vehicle has landed and must be hidden */
        if(data.altitude == 0 ){
            this.#show = false;
        }

        /*Check if vehicle has died */
        if((this.#alive !== data.alive) && (this.#alive == true)){
            window.ui.alert(this.#name + " has been destroyed");
            this.#alive = false;
            this.#showing = true;
            this.#show = true;
        }

        /*Check if vehicle has been respawned */
        if((this.#alive !== data.alive) && (this.#alive == false)){
            window.ui.alert(this.#name + " has been replaced into the " + this.#team + " fleet");
            this.#alive = true;
            this.#showing = true;
            this.#show = true;
        }

        this.#alive = data.alive;

        this.classify();

        /* Update the class and it's L.Marker */ 
        super.update(data.latitude, data.longitude, data.altitude);
        this.#icon.setLatLng(this.latLng());
        
        var builder = "<b onclick=\"window.ui.showVehicle(" + data.id + ");\" style=\"text-align:center\">" + data.name + "</b><br/>" + this.latitude.toFixed(4) + ", " + this.longitude.toFixed(4) + "<br/>" + this.altitude.toFixed(1) + "<br/>"
        if (map.targetMarker != null){
            var dist = map.targetMarker.getLatLng().distanceTo(this.latLng()).toFixed(0)
            builder += "Dist to Marker: " + dist + "<br/>"
            var speed = 90;
            switch(this.#role){
                case "ISR":
                    speed = 120
                    break
                case "FIGHTER":
                    speed = 165
                    break;
                case "BOMBER":
                    speed = 90
                    break;
                case "ANCALAGON":
                    speed = 90
                    break;
                case "SAGROTH":
                    speed = 60;
                    break;
                case "MULTI":
                    speed = 105;
                    break;
                default:
                    speed = 90
            }

            var minutes = ("" + Math.floor(dist/speed/60)).padStart(2, "0")
            var seconds = ("" + Math.floor((dist/speed) - (minutes * 60))).padStart(2, "0")

            builder += "Time2Target: " + minutes + ":" + seconds
            this.#icon.bindPopup(builder);
        }else{
            this.#icon.bindPopup(builder);
        }
    }

    /**
     * Shows the popup on screen
     */
    show(){
        this.#icon.openPopup();
    }


    /**
     * Assigns the correct icon for the vehicle. Logic to show or hide a vehicle also exists here.
     */
    classify(){
            
        if((this.#showing == false) && (this.#show == true)){
            this.#showing = true;
            switch(this.#role){
                case "ISR":
                    this.#icon.setIcon(
                    L.icon({
                        iconUrl: 'assets/svgs/drone.svg',
                        iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                        iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                        })
                    );
                    
                break;

                case "BOMBER":
                this.#icon.setIcon(
                    L.icon({
                        iconUrl: 'assets/svgs/bomber.svg',
                        iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                        iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                    })
                );
                break;

                case "FIGHTER":
                    this.#icon.setIcon(
                        L.icon({
                            iconUrl: 'assets/svgs/fighter.svg',
                            iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                            iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                        })
                    );
                break;

                case "LOGISTICS":
                    this.#icon.setIcon(
                        L.icon({
                            iconUrl: 'assets/svgs/logistics.svg',
                            iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                            iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                        })
                    );
                break;

                case "MULTI":
                    this.#icon.setIcon(
                        L.icon({
                            iconUrl: 'assets/svgs/drone.svg',
                            iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                            iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                        })
                    );
                break;

                case "WMD":
                    this.#icon.setIcon(
                        L.icon({
                            iconUrl: 'assets/svgs/nuclear.svg',
                            iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                            iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                        })
                    );
                break;
            
                default:
                    this.#icon.setIcon(
                        L.icon({
                            iconUrl: 'assets/svgs/drone.svg',
                            iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                            iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                        })
                    );
                }

                if (this.#alive == false){
                    this.#icon.setIcon(
                        L.icon({
                            iconUrl: 'assets/svgs/dead.svg',
                            iconSize:     [this.#vehicleSize, this.#vehicleSize],
                            iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2],
                        })
                    );
                }
        
        
            this.teamifyIcon();
        } 
        else if ((this.#showing == true) && (this.#show == false)) {
            this.#showing = false;
            this.#icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/drone.svg',
                    iconSize:     [0, 0], 
                    iconAnchor:   [0, 0], 
                })
            );
        }
    }


    isAlive(){
        return this.#alive;
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

    getRole(){
        return this.#role;
    }

    toJSON(){
        return {
            "id" : this.#id,
            "name" : this.#name,
            "role" : this.#role,
            "color" : this.color,
            "team" : this.#team,
            "latitude" : this.latitude,
            "longitude" : this.longitude,
            "altitude" : this.altitude,
            "type" : this.#type,
            "alive" : this.#alive
        }
    }


    /**
     * Uses the GeoAPI to determine which terratory this vehicle is in and what ground type it is on
     */
    async geoAPI(){
        return new Promise((resolve, reject) => {
            $.ajax({
                type: 'GET',
                url: "http://maps.ace:3000/" + this.latitude + "/" + this.longitude + "/" + this.altitude,
                success: function(data) {
                    resolve({
                        "nation": data.nation,
                        "groundType": data.groundType
                    });
                },
                error: function(data){
                    reject(data);
                }
            });
        });
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


    
    /** 
     * This is a mangement function, interns can ignore this as you cannot reach API endpoints.
     * Attempts to reload a vehicle using the vehicles.mgmt API
     */
     reload(){
        var confirmation = confirm("Do you really want to restart " + this.#name + "?");
        if (confirmation == true) {
            const url='http://vehicles.mgmt/api/reload/' + this.#name;
            var xmlHttp = new XMLHttpRequest();
            xmlHttp.onreadystatechange = function (){
            if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
                console.log(xmlHttp.responseText);
            }
            xmlHttp.open("GET", url, true);
            xmlHttp.send();
            return;
        } 
        else {
            return;
        }
    }

    
    /** 
     * This is a mangement function, interns can ignore this as you cannot reach API endpoints.
     * Attempts to restart a vehicle using the vehicles.mgmt API
     */
    restart(){
        var confirmation = confirm("Do you really want to restart " + this.#name + "?");
        if (confirmation == true) {
            const url='http://vehicles.mgmt/api/restart/' + this.#name;
            var xmlHttp = new XMLHttpRequest();
            xmlHttp.onreadystatechange = function (){
            if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
                console.log(xmlHttp.responseText);
            }
            xmlHttp.open("GET", url, true);
            xmlHttp.send();
            return;
        } 
        else {
            return;
        }
    }

    /** 
     * This is a mangement function, interns can ignore this as you cannot reach API endpoints.
     * Attempts to stop a vehicle using the vehicles.mgmt API
     */
    stop(){
        var confirmation = confirm("Do you really want to stop " + this.#name + "?");
        if (confirmation == true) {
            const url='http://vehicles.mgmt/api/stop/' + this.#name;
            var xmlHttp = new XMLHttpRequest();
            xmlHttp.onreadystatechange = function (){
            if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
                console.log(xmlHttp.responseText);
            }
            xmlHttp.open("GET", url, true);
            xmlHttp.send();
            return;
        } 
        else {
            return;
        }
    }


    /**
     * Launches Model that provides options for rehoming a drone
     */
    rehomeDroneUI(){
        let optionDropdown = "<select id=\"baseSelect\" class=\"form-control\">";
        for(let infrastructure in window.delphiMap.infrastructure){
            infrastructure = window.delphiMap.infrastructure[infrastructure];
            if(infrastructure.getTeam() == this.#team){
                console.log(infrastructure);
                optionDropdown += "<option value=\" " + infrastructure.getID() + "\">" + infrastructure.getName() + "</option>";
            }
        }
        optionDropdown += "</select><br />";
        optionDropdown += "<button data-dismiss=\"modal\" aria-label=\"Close\" onclick=\"window.delphiMap.vehicles[" + this.#id + "].rehome()\" class=\"btn btn-primary\">Rehome</button><br />";
        $('#modalBody').html(optionDropdown);
        window.ui.createModal("Rehome " + this.#name, optionDropdown);
    }

    /**
     * Attempts to rehome a vehicle using the vehicles.mgmt API - requires UI
     */
    rehome(){
        console.log($('#baseSelect :selected').val());
        window.delphiMap.infrastructure.forEach(function(element){
        if(element.getID() == $('#baseSelect :selected').val()){
            console.log("REHOME " + this.#name + " TO " + element.getName() + " AT LAT: " + element.latitude + " AT LONG: " + element.longitude);
            const url='http://vehicles.mgmt/api/rehome/' + this.#name + "/" + element.latitude + "/" + element.longitude;
            var xmlHttp = new XMLHttpRequest();
            xmlHttp.onreadystatechange = function (){
            if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
                console.log(xmlHttp.responseText);
            }
            xmlHttp.open("GET", url, true);
            xmlHttp.send();
            return;
        }
        });
    } 



    

    
}