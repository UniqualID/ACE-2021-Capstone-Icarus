import Location from './Location.js';

/**
 * Vehicle class that contains logic for Vehicle Operations
 * Extends the Location object
 */
export default class IADS extends Location{

    #icon;
    #radiusCircle;
    #name;
    #role;
    #alive;
    color;
    #team;
    #type;
    #scale;
    #vehicleSize;
    #id;
    
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

        //Ranges of assets

        {
            let rangeSAM = {};
            rangeSAM['red'] = 30000;
            rangeSAM['blue'] = 40000;
            rangeSAM['green'] = 25000;
            rangeSAM['purple'] = 40000;
            let rangeRADAR = {};
            rangeRADAR['red'] = 65000;
            rangeRADAR['blue'] = 80000;
            rangeRADAR['green'] = 80000;
            rangeRADAR['purple'] = 80000; 
            let radarColor = '#7F00FF';
            let samColor = '#e54f2a';
            this.#radiusCircle;
            switch(this.#role){
                case "RADAR":
                this.#radiusCircle = L.circle(this.latLng(), rangeRADAR[this.color], {color: radarColor, weight: 0.75, fill: false}).addTo(window.map.iadsLayer);
                break;
        
                case "SAM":
                this.#radiusCircle = L.circle(this.latLng(), rangeSAM[this.color], {color: samColor, weight: 0.75, fill: false}).addTo(window.map.iadsLayer);
                break;
            }

        }


        this.#scale = 1.0;
        this.#vehicleSize = 15*this.#scale;

        this.#icon = L.marker([this.latitude, this.longitude], 
            {icon: 
                L.icon({
                    iconUrl: 'assets/svgs/radar.svg',
                    iconSize:     [0, 0], 
                    iconAnchor:   [0, 0], 
                })
        }).addTo(window.map.iadsLayer);


        this.#icon.bindPopup("<b onclick=\"window.ui.showIad(" + data.id + ");\" style=\"text-align:center\">" + data.name + "</b><br/>" + this.latitude.toFixed(4) + ", " + this.longitude.toFixed(4) + "<br/>" + this.altitude.toFixed(1) + "<br/>");
            

        //Sends the icon to be classified
        this.classify();


    }

    /**
     * Move the vehicle on the map
     * @param {Object} data - Individual Vehicle object taken from the JSON file
     */
    update(data){
        /*Check if vehicle has died */
        if((this.#alive !== data.alive) && (this.#alive == true)){
            this.#alive = false;
            this.#radiusCircle.remove();
        }
        if(!this.#alive){
            this.#radiusCircle.remove();
        }
        this.#alive = data.alive;
    }


    /**
     * Assigns the correct icon for the vehicle. Logic to show or hide a vehicle also exists here.
     */
    classify(){
        switch(this.#role){
            case "SAM":
                this.#icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/sam.svg',
                    iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                    iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                    })
                );

            break;

            case "RADAR":
            this.#icon.setIcon(
                L.icon({
                    iconUrl: 'assets/svgs/radar.svg',
                    iconSize:     [this.#vehicleSize, this.#vehicleSize], // size of the icon
                    iconAnchor:   [this.#vehicleSize/2, this.#vehicleSize/2], // point of the icon which will correspond to marker's location
                })
            );
            break;
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
        optionDropdown += "<button data-dismiss=\"modal\" aria-label=\"Close\" onclick=\"window.delphiMap.iads[" + this.#id + "].rehome()\" class=\"btn btn-primary\">Rehome</button><br />";
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