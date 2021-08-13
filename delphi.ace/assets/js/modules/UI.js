export default class UI{
    
    #sidenav;
    #sidebar;
    #sidebarContent;
    #modal;
    #modalTitle;
    #modalBody;
    #alert;
    #feedArray;
    #intervalTimer;
    /**
     * Sets up the UI Object and enables interactions with the sidebar.
     */
    constructor(){
        this.#sidenav = $('#sidenav');
        this.#sidebar = $('#viewScreen');
        this.#sidebarContent = $('#viewScreenContent');
        this.#modal = $('#delphiModal');
        this.#modalTitle = $('#modalTitle');
        this.#modalBody = $('#modalBody');
        this.#alert = $('#alert');
        this.#feedArray = [];
        this.#intervalTimer = null;
    }

    /**
     * 
     * @param {String} item - Pushes an Alert to the news feed
     */
    pushToFeed(item){
        this.#feedArray.push(item);
    }


    showSidenav(){
        this.#sidenav.show(1000);
    }

    /**
     *  Opens the sidebar (Private to avoid outside interaction)
     */
    #showSidebar(){
        this.#sidebar.show(800);
    }


    /**
     * Bring up the Model Element
     */
    #showModal(){
        this.#modal.modal('show');
        console.log(this.#modal);
    }


    createModal(title, content){
        $('#modalTitle').html(title);
        $('#modalBody').html(content);
        this.#showModal();
    }

    /**
     * Hides the sidebar
     */
    hideSidebar(){
        this.#sidebar.hide();
        this.#resetInterval();
    }


    /**
     * 
     * @param {String} title - Title for the sidebar
     * @param {String} content - The contnet to be displayed
     */
    #fillSidebar(title, content){
        let builder =  "<h3>" + title + "</h3>";
        builder += content;
        this.#sidebarContent.html(builder);
        this.#showSidebar();
    }


    /**
     * Triggers an onscreen alert on the screen. Automatically pushes to news feed
     * @param {String} content - Alert content. String or HTML String
     */
    alert(content){
        this.#alert.html(content);
        let datetime = new Date()
        this.pushToFeed(datetime.getHours() + ":" + datetime.getMinutes() + ":" + datetime.getSeconds() + " - " + content);
        this.#alert.show()
        setTimeout(function(){
            this.#alert.hide();;
        }.bind(this), 2000);
    }

    /**
     * Displays a mission sidebar. Does not update.
     */
    showMission(){
        this.#resetInterval();
        let mission = window.delphiMap.getMission();
        let builder = "<table class=\"table\"><tr><th scope=\"row\">Mission Name:</th><td>" + mission.name + "</td></tr><tr><th scope=\"row\">Maximum Latitude:</th><td>" + mission.maximum_latitude + "</td></tr><tr><th scope=\"row\">Maximum Longitude:</th><td>" + mission.maximum_longitude + "</td></tr><tr><th scope=\"row\">Minimum Latitude:</th><td>" + mission.minimum_latitude + "</td></tr><tr><th scope=\"row\">Minimum Longitude:</th><td>" + mission.minimum_longitude + "</td></tr></table><p>";
        builder += "<center><button class=\"btn btn-primary\" onclick=\"window.ui.showJSONupload();\"> Upload Custom JSON </button></center>";
        this.#fillSidebar("Operation Details", builder);
    }


    /**
     * Launches Model to upload a new JSON file
     */
     showJSONupload(){
        let uploadHTML = "Please upload a JSON file: <br /><input id=\"jsonupload\" type=\"file\" class=\"form-control\"/><br /> ";
        uploadHTML += "<button data-dismiss=\"modal\" aria-label=\"Close\" onclick=\"window.delphiMap.importJSON()\" class=\"btn btn-primary\">Add Data</button><br />";
        window.ui.createModal("Upload a mission data file", uploadHTML);
    }

    /**
     * Show the news feed on the screen
     */
    showNewsFeed(){
        this.#resetInterval();
        this.#intervalTimer = null;
        let builder = "<table class=\"table\">";
        if(this.#feedArray.length == 0 ){
            builder += "<tr><td> The Feed is Empty - check back soon...</td></tr>"
        }
        for (var i = this.#feedArray.length - 1; i >= 0; i--) {
            builder += "<tr><td>" + this.#feedArray[i] + "</td></tr>"
        }
        builder += "</table>";
        this.#fillSidebar("Feed", builder);
        
        //Code to refresh a UI Window
        if(this.#intervalTimer == null){
            this.#intervalTimer = setInterval(function(){this.showNewsFeed()}.bind(this), 3000);
        }
    }


    /**
     * Shows all vehicle objects in a list in the sidebar
     * @param {String} - Team to show vehicles. Null means all.
     */
    showAllVehicles(team){
        this.#resetInterval();
        let sortedVehicles = this.#sortJSON(Object.values(window.delphiMap.vehicles), 'color');
        let builder = "";
        builder += "<a href=\"#\" onclick=\"window.ui.showAllVehicles()\">All</a> ";
        for(let z = 0; z < window.delphiMap.getTeams().length; z++){
            builder += "<a href=\"#\" onclick=\"window.ui.showAllVehicles('" + window.delphiMap.getTeams()[z] + "', 'infrastructure')\">" + window.delphiMap.getTeams()[z] + "</a> ";
        }
        builder += "<table class=\"table\">";
        for(let i = 0; i < sortedVehicles.length; i++){
            if((team == undefined) || (team == sortedVehicles[i].getTeam())){
                builder +=  "<tr class=\"asset team-" + sortedVehicles[i].getTeam() + "\" style=\"cursor: pointer;\" onclick=\"window.ui.showVehicle(" + sortedVehicles[i].getID() + ");\"><th scope=\"row\" style=\"color:" + sortedVehicles[i].getColor() +"\">";
                if(sortedVehicles[i].isAlive() == true){
                    builder += "<div class=\"dot bg-success\"></div> ";
                } 
                else{
                    builder += "<div class=\"dot bg-primary\"></div> ";
                }
                builder += sortedVehicles[i].getName() + " [" + sortedVehicles[i].getID() + "]</th><td>" + sortedVehicles[i].getRole() + "</td></tr>";
            }
        }
        builder += "</table>"
        this.#fillSidebar("Vehicles", builder);
        //Code to refresh a UI Window
        if(this.#intervalTimer == null){
            this.#intervalTimer = setInterval(function(){this.showAllVehicles(team)}.bind(this), 3000);
        }
    }

    /**
     * Shows all Infrastructure objects in a list in the sidebar
     * @param {String} - Team to show vehicles. Null means all.
     */
    showAllInfrastructure(team){
        //Code to stop another UI Refreshing loads
        this.#resetInterval();
        let sortedInfrastructure = this.#sortJSON(Object.values(window.delphiMap.infrastructure), 'color');
        let builder = "";
        builder += "<a href=\"#\" onclick=\"window.ui.showAllInfrastructure()\">All</a> ";
        for(let z = 0; z < window.delphiMap.getTeams().length; z++){
            builder += "<a href=\"#\" onclick=\"window.ui.showAllInfrastructure('" + window.delphiMap.getTeams()[z] + "', 'infrastructure')\">" + window.delphiMap.getTeams()[z] + "</a> ";
        }
        builder += "<table class=\"table\">";
        for(let i = 0; i < sortedInfrastructure.length; i++){
            if((team == undefined) || (team == sortedInfrastructure[i].getTeam())){
                builder +=  "<tr class=\"asset team-" + sortedInfrastructure[i].getTeam() + "\" style=\"cursor: pointer;\" onclick=\"window.ui.showInfrastructure(" + sortedInfrastructure[i].getID() + ")\"><th scope=\"row\" style=\"color:" + sortedInfrastructure[i].getColor() +"\">";
                if(sortedInfrastructure[i].isAlive() == true){
                    builder += "<div class=\"dot bg-success\"></div> ";
                } 
                else{
                    builder += "<div class=\"dot bg-primary\"></div> ";
                }
                builder += sortedInfrastructure[i].getName() + " [" + sortedInfrastructure[i].getID() + "]</th><td>" + sortedInfrastructure[i].getType() + "</td></tr>";
            }
        }
        builder += "</table>"
        this.#fillSidebar("Infrastructure", builder);
        //Code to refresh a UI Window
        if(this.#intervalTimer == null){
            this.#intervalTimer = setInterval(function(){this.showAllInfrastructure(team)}.bind(this), 3000);
        }
    }

    /**
     * Shows all Defenses objects in a list in the sidebar
     * @param {String} - Team to show IADS. Null means all.
     */
    showAllIads(team){
        //Code to stop another UI Refreshing loads
        this.#resetInterval();
        let sortedIads = this.#sortJSON(Object.values(window.delphiMap.iads), 'color');
        let builder = "";
        builder += "<a href=\"#\" onclick=\"window.ui.showAllIads()\">All</a> ";
        for(let z = 0; z < window.delphiMap.getTeams().length; z++){
            builder += "<a href=\"#\" onclick=\"window.ui.showAllIads('" + window.delphiMap.getTeams()[z] + "', 'infrastructure')\">" + window.delphiMap.getTeams()[z] + "</a> ";
        }
        builder += "<table class=\"table\">";
        for(let i = 0; i < sortedIads.length; i++){
            if((team == undefined) || (team == sortedIads[i].getTeam())){
                builder +=  "<tr class=\"asset team-" + sortedIads[i].getTeam() + "\" style=\"cursor: pointer;\" onclick=\"window.ui.showIad(" + sortedIads[i].getID() + ");\"><th scope=\"row\" style=\"color:" + sortedIads[i].getColor() +"\">";
                if(sortedIads[i].isAlive() == true){
                    builder += "<div class=\"dot bg-success\"></div> ";
                } 
                else{
                    builder += "<div class=\"dot bg-primary\"></div> ";
                }
                builder += sortedIads[i].getName() + " [" + sortedIads[i].getID() + "]</th><td>" + sortedIads[i].getRole() + "</td></tr>";
            }
        }
        builder += "</table>"
        this.#fillSidebar("IADS", builder);
        //Code to refresh a UI Window
        if(this.#intervalTimer == null){
            this.#intervalTimer = setInterval(function(){this.showAllIads(team)}.bind(this), 3000);
        }
    }


    /**
     * Fills the sidebar with info related to an individual vehicle
     * @param {int} id - ID value of the Vehicle
     */
    async showVehicle(id){
        //Code to stop another UI Refreshing loads
        this.#resetInterval();
        let vehicle = window.delphiMap.vehicles[id].toJSON();
        let builder = "";
        let header = "";
        if(vehicle.alive == true){
             header = "<div onclick=\"window.delphiMap.vehicles[" + vehicle.id + "].show();\"><div class=\"dot bg-success\" style=\"vertical-align: middle;\"></div> " + vehicle.name + "[" + vehicle.id + "]</div>";
        } 
        else{
            header = "<div onclick=\"window.delphiMap.vehicles[" + vehicle.id + "].show();\"><div class=\"dot bg-primary\" style=\"vertical-align: middle;\"></div> " + vehicle.name + "[" + vehicle.id + "]</div>";
        }
        builder += "<table class=\"table\"><tr><th scope=\"row\"> Team: </th><td style=\"color:" + vehicle.color +"\">" + vehicle.team  + "</td></tr>";
        builder += "<tr><th scope=\"row\"> Role: </th><td>" + vehicle.role + "</td></tr>";
        builder += "<tr><th scope=\"row\"> Latitude: </th><td>" + vehicle.latitude.toFixed(4) + "</td></tr>";
        builder += "<tr><th scope=\"row\"> Longitude: </th><td>" + vehicle.longitude.toFixed(4) + "</td></tr>";
        builder += "<tr><th scope=\"row\"> Altitude: </th><td>" + vehicle.altitude.toFixed(4) + "</td></tr>";
        let geoAPIInfo = await window.delphiMap.vehicles[id].geoAPI().catch((err) => {
            console.log("Failed API call");
            this.alert("Could not connect to GeoAPI");
        });
        if(!(geoAPIInfo === undefined)){ 
            builder += "<tr><th scope=\"row\"> Territory: </th><td>" + geoAPIInfo.nation + "</td></tr>";
            builder += "<tr><th scope=\"row\"> Ground Type: </th><td>" + geoAPIInfo.groundType + "</td></tr>";
        }
        if (window.delphiMap.isMGMT() == true){
            builder += "<tr><td><button class=\"btn btn-primary\" onclick=\"window.delphiMap.vehicles[" + vehicle.id + "].restart();\"> Restart Vehicle </button> </td><td><button class=\"btn btn-primary\" onclick=\"window.delphiMap.vehicles[" + vehicle.id + "].stop()\"> Stop Vehicle </button> </td></tr>";
            builder += "<tr><td><button class=\"btn btn-primary\" onclick=\"window.delphiMap.vehicles[" + vehicle.id + "].rehomeDroneUI();\"> Rehome Vehicle </button> </td><td><button class=\"btn btn-primary\" onclick=\"window.delphiMap.vehicles[" + vehicle.id + "].reload()\"> Reload Vehicle </button> </td></tr>";
        }
        builder += "</table>";

        window.ui.#fillSidebar(header, builder);

        //Starts the refresh process
        if(this.#intervalTimer == null){
            this.#intervalTimer = setInterval(function(){this.showVehicle(vehicle.id)}.bind(this), 3000);
        }
    }

    /**
     * Fills the sidebar with info related to an individual Infrastructure unit
     * @param {int} id - ID Value of the IAD Unit
     */
     async showInfrastructure(id){
        //Code to stop another UI Refreshing loads
        this.#resetInterval();
        let infrastructure = window.delphiMap.infrastructure[id].toJSON();
        let builder = "";
        let header = "";
        if(infrastructure.alive == true){
             header = "<div class=\"dot bg-success\" style=\"vertical-align: middle;\"></div> " + infrastructure.name + "[" + infrastructure.id + "]";
        } 
        else{
            header = "<div class=\"dot bg-primary\" style=\"vertical-align: middle;\"></div> " + infrastructure.name + "[" + infrastructure.id + "]";
        }
        builder += "<table class=\"table\"><tr><th scope=\"row\"> Team: </th><td style=\"color:" + infrastructure.color +"\">" + infrastructure.team  + "</td></tr>";
        builder += "<tr><th scope=\"row\"> Type: </th><td>" + infrastructure.type + "</td></tr>";
        builder += "<tr><th scope=\"row\"> Latitude: </th><td>" + infrastructure.latitude.toFixed(4) + "</td></tr>";
        builder += "<tr><th scope=\"row\"> Longitude: </th><td>" + infrastructure.longitude.toFixed(4) + "</td></tr>";
        builder += "</table>";

        window.ui.#fillSidebar(header, builder);

        //Starts the refresh process
        if(this.#intervalTimer == null){
            this.#intervalTimer = setInterval(function(){this.showInfrastructure(infrastructure.id)}.bind(this), 3000);
        }
    }


    /**
     * Fills the sidebar with info related to an individual IAD unit
     * @param {int} id - ID Value of the IAD Unit
     */
    async showIad(id){
        //Code to stop another UI Refreshing loads
        this.#resetInterval();
        let iad = window.delphiMap.iads[id].toJSON();
        let builder = "";
        let header = "";
        if(iad.alive == true){
             header = "<div class=\"dot bg-success\" style=\"vertical-align: middle;\"></div> " + iad.name + "[" + iad.id + "]";
        } 
        else{
            header = "<div class=\"dot bg-primary\" style=\"vertical-align: middle;\"></div> " + iad.name + "[" + iad.id + "]";
        }
        builder += "<table class=\"table\"><tr><th scope=\"row\"> Team: </th><td style=\"color:" + iad.color +"\">" + iad.team  + "</td></tr>";
        builder += "<tr><th scope=\"row\"> Role: </th><td>" + iad.role + "</td></tr>";
        builder += "<tr><th scope=\"row\"> Latitude: </th><td>" + iad.latitude.toFixed(4) + "</td></tr>";
        builder += "<tr><th scope=\"row\"> Longitude: </th><td>" + iad.longitude.toFixed(4) + "</td></tr>";
        builder += "<tr><th scope=\"row\"> Altitude: </th><td>" + iad.altitude.toFixed(4) + "</td></tr>";
        if (window.delphiMap.isMGMT() == true){
            builder += "<tr><td><button class=\"btn btn-primary\" onclick=\"window.delphiMap.iads[" + iad.id + "].restart();\"> Restart IAD </button> </td><td><button class=\"btn btn-primary\" onclick=\"window.delphiMap.iads[" + iad.id + "].stop()\"> Stop IAD </button> </td></tr>";
            builder += "<tr><td><button class=\"btn btn-primary\" onclick=\"window.delphiMap.iads[" + iad.id + "].rehomeDroneUI();\"> Rehome IAD </button> </td><td><button class=\"btn btn-primary\" onclick=\"window.delphiMap.iads[" + iad.id + "].reload()\"> Reload IAD </button> </td></tr>";        
        }
        builder += "</table>";

        window.ui.#fillSidebar(header, builder);

        //Starts the refresh process
        if(this.#intervalTimer == null){
            this.#intervalTimer = setInterval(function(){this.showIad(iad.id)}.bind(this), 3000);
        }
    }


  /**
   * Allow a user to search for a vehicle by name or value...
   */
   showSearch(query){
        this.#resetInterval();
        var builder = "";
        builder += "<div class=\"input-group mb-3\"><input id=\"searchBox\" type=\"text\" class=\"form-control\" aria-describedby=\"basic-addon2\"><div class=\"input-group-append\"><button onclick=\"window.ui.showSearch(document.getElementById('searchBox').value)\" class=\"btn btn-primary\" type=\"button\">Search</button></div></div>";
        if(!((query == null)||(query == undefined))){
            builder += "<p> Search Query: " + query ;
            builder += "<br/>Vehicles:</br><table class=\"table\">";

            let vehicle = this.#sortJSON(Object.values(window.delphiMap.vehicles), 'color');
            for(let i = 0; i < vehicle.length; i++){
                if((vehicle[i].getID() == query)||(vehicle[i].getName().includes(query))){
                    if(!((vehicle[i].getRole() == "SAM")||(vehicle[i].getRole() == "RADAR")||(vehicle[i].latitude == 0))){
                        if(true){
                            builder += "<tr class=\"asset team-" + vehicle[i].team + "\" style=\"cursor: pointer;\" onclick=\"window.ui.showVehicle(" + vehicle[i].getID() + ")\"><th scope=\"row\" style=\"color:" + vehicle[i].getColor() +"\">";
                            if(vehicle[i].isAlive() == true){
                                builder += "<div class=\"dot bg-success\"></div> ";
                            } else{
                                builder += "<div class=\"dot bg-primary\"></div> ";
                            }
                            builder += vehicle[i].getName() + " [" + vehicle[i].getID() + "]</th><td>" + vehicle[i].getRole() + "</td></tr>";
                        }
                    }
                }
            }
            //Infrastructure
            builder += "</table><br/>Infrastructure:</br><table class=\"table\">";
            let infrastructure = this.#sortJSON(Object.values(window.delphiMap.infrastructure), 'color');
            for(let i = 0; i < infrastructure.length; i++){
                if((infrastructure[i].getID() == query)||(infrastructure[i].getName().includes(query))){
                    if(infrastructure[i].isAlive() == true ){
                        builder += "<tr style=\"cursor: pointer;\" onclick=\"window.ui.showInfrastructure(" + infrastructure[i].getID() + ")\"><th  scope=\"row\" style=\"color:" + infrastructure[i].getColor() +"\"> <div class=\"dot bg-success\"></div> " + infrastructure[i].getName() + " [" + infrastructure[i].getID() + "]</th><td>" + infrastructure[i].getType() + "</td></tr>";
                    } else{
                        builder += "<tr  style=\"cursor: pointer;\" onclick=\"window.ui.showInfrastructure(" + infrastructure[i].getID() + ")\"><th scope=\"row\" style=\"color:" + infrastructure[i].getColor() +"\"> <div class=\"dot bg-primary\"></div> " + infrastructure[i].getName() +" [" + infrastructure[i].getID() + "]</th><td>" + infrastructure[i].getType() + "</td></tr>";
                    }
                }
            }
            builder += "</table>";
            //Defenses
            let defense = this.#sortJSON(Object.values(window.delphiMap.iads), 'color');
            builder += "</table><br/>IADS:</br><table class=\"table\">";
            for(let i = 0; i < defense.length; i++){
            if(((defense[i].getRole() == "SAM") || (defense[i].getRole() == "RADAR")) & (defense[i].latitude !== 0) ){
                if((defense[i].getID() == query)||(defense[i].getName().includes(query))){
                if(defense[i].isAlive() == true){
                    builder += "<tr style=\"cursor: pointer;\" onclick=\"window.ui.showIad(" + defense[i].getID() + ")\" ><th scope=\"row\" style=\"color:" + defense[i].getColor() +"\"> <div class=\"dot bg-success\"></div> " + defense[i].getName() + " [" + defense[i].getID() + "]</th><td>" + defense[i].getRole() + "</td></tr>";
                } else{
                    builder += "<tr style=\"cursor: pointer;\" onclick=\"window.ui.showIad(" + defense[i].getID() + ")\" ><th scope=\"row\" style=\"color:" + defense[i].getColor() +"\"> <div class=\"dot bg-primary\"></div> " + defense[i].getName() + " [" + defense[i].getID() + "]</th><td>" + defense[i].getRole() + "</td></tr>";
                }
                }
            }
            }
            builder += "</table>";
        }

        window.ui.#fillSidebar("Search", builder);
    }    

    /**
     * Show all assets in set bounds
     */
    showInBounds(boxStart, boxEnd){
        var builder = "";
        if ((boxStart == null)||(boxEnd == null)){
            builder += "<p> Use CTRL - Drag to select a region on the map </p>";
        } else{
            builder += "<p>Use ALT - CLICK to diable the region</p>";
            var minLat = null;
            var maxLat = null;
            var minLong = null;
            var maxLong = null;
            if (boxStart.lat > boxEnd.lat){
            minLat = boxEnd.lat;
            maxLat = boxStart.lat;
            } else{
            maxLat = boxEnd.lat;
            minLat = boxStart.lat;
            }
            if (boxStart.lng > boxEnd.lng){
            minLong = boxEnd.lng;
            maxLong = boxStart.lng;
            } else{
            maxLong = boxEnd.lng;
            minLong = boxStart.lng;
            }
            builder += "<h3> Vehicles:</h3><table class=\"table\">";
            let vehicles = this.#sortJSON(Object.values(window.delphiMap.vehicles), 'color');
            for(let i = 0; i < vehicles.length; i++){
                vehicle = window.delphiMap.vehicles[vehicle].getJSON();
                console.log(vehicle);
                if((vehicle[i].latitude >= minLat) && (vehicle[i].latitude <= maxLat) && (vehicle[i].longitude >= minLong) && (vehicle[i].longitude <= maxLong)){
                    if(!((vehicle[i].role == "SAM")||(vehicle[i].role == "RADAR")||(vehicle[i].latitude == 0))){
                        console.log("There is a vehicle");
                        if(vehicle[i].alive == true){
                            builder += "<tr style=\"cursor: pointer;\" onclick=\"window.ui.showVehicle(" + vehicle[i].id + ");\"><th scope=\"row\" style=\"color:" + vehicle[i].color +"\"> <div class=\"dot bg-success\"></div> ";
                        }
                        else{
                            builder += "<tr style=\"cursor: pointer;\" onclick=\"window.ui.showVehicle(" + vehicle[i].id + ");\"><th scope=\"row\" style=\"color:" + vehicle[i].color +"\"> <div class=\"dot bg-primary\"></div> ";
                        }
                        builder += vehicle[i].name + " [" + vehicle[i].id + "]</th><td>" + vehicle[i].role + "</td></tr>";
                        }
                    }
            }
            //Plot Infrastructure
            builder += "</table>";
            builder += "<h3> Infrastructure:</h3><table class=\"table\">";
            let infrastructures = this.#sortJSON(Object.values(window.delphiMap.infrastructure), 'color');
            for(let i = 0; i < infrastructures.length; i++){
                infrastructure = window.delphiMap.infrastructure[infrastructures].getJSON();
                if((infrastructure[i].latitude >= minLat) && (infrastructure[i].latitude <= maxLat) && (infrastructure[i].longitude >= minLong) && (infrastructure[i].longitude <= maxLong)){
                if(infrastructure[i].hitpoints > 0 ){
                    builder += "<tr><th scope=\"row\" style=\"color:" + infrastructure[i].color +"\"> <div class=\"dot bg-success\"></div> " + infrastructure[i].name + " [" + infrastructure[i].id + "]</th><td>" + infrastructure[i].type + "</td></tr>";
                } else{
                    builder += "<tr><th scope=\"row\" style=\"color:" + infrastructure[i].color +"\"> <div class=\"dot bg-primary\"></div> " + infrastructure[i].name + " [" + infrastructure[i].id + "]</th><td>" + infrastructure[i].type + "</td></tr>";
                }
                }
            }
            builder += "</table>";
        }
        window.ui.#fillSidebar("Selection", builder);

    }


    /**
     * A handy function for sorting JSON Objects by key. Cheers StackOverflow
     * @param {JSON} data JSON data to be sorted
     * @param {string} key The key to sort the data by
     * @returns {JSON} JSON data sorted by the defined key
     */

    #sortJSON(data, key) {
        return data.sort(function(a, b) {
            var x = a[key]; var y = b[key];
            return ((x < y) ? -1 : ((x > y) ? 1 : 0));
        });
    }

    #resetInterval(){
        clearInterval(this.#intervalTimer);
        this.#intervalTimer = null;
    }


    /**
     * Valinar custom sidebar content
     */
     async showValinar(){
        let builder = "";
        let header = "";
        header = "Valinar";

        /* Coordinates lookup */
        builder += "<div class=\"container\">"
        builder +=  "<div class=\"row\">"
        builder +=      "<span><strong>Lat/Long Lookup</strong></span>";
        builder +=      "<div class=\"md-form amber-textarea active-amber-textarea\"></div>"
        builder +=          "<textarea id=\"lat-lon-lookup\" class=\"md-textarea form-control\" rows=\"3\" placeholder=\"Input lat,long coordinates separated by a comma\"></textarea>"
        builder +=          "<button onclick=\"window.map.val_plotPoints(document.getElementById('lat-lon-lookup').value)\" id=\"submit\" type=\"button\" class=\"btn btn-primary\">PLOT</button>"
        builder +=      "</div>"
        //builder +=  "</div>"

        /* Plotting controls */
        builder +=  "<div class=\"row\">"
        builder +=      "<span><strong>Plot waypoints</strong></span><br>"
        builder +=      "<div class=\"btn-group btn-group-toggle\" data-toggle=\"buttons\">"
        builder +=          "<label class=\"btn btn-secondary active\" id=\"plot-inactive\">"
        builder +=              "<input type=\"radio\" name=\"options\" id=\"no-plot\" checked>Not plotting</input>"
        builder +=          "</label>"
        builder +=          "<label class=\"btn btn-secondary\" id=\"plot-active\">"
        builder +=              "<input type=\"radio\" name=\"options\" id=\"plot-active\">Plotting</input>"
        builder +=          "</label>"
        builder +=      "</div>"
        builder +=  "</div>"
        builder +=  "<div class=\"row\">"
        builder +=      "<button type=\"button\" class=\"btn btn-warning\" onclick=\"window.ui.clearMarkers();\">Clear Markers</button>"
        
        builder +=      "<select class=\"custom-select\" id=\"wp-vehicle-select\">"
        builder +=      "<option value=\"90\">ANCALAGON (Bomber)</option>"
        builder +=      "<option value=\"165\">GLAURUNG (Fighter)</option>"
        builder +=      "<option value=\"135\">KIRINKI (ISR)</option>"
        builder +=      "<option value=\"60\">SAGROTH (WMD)</option>"
        builder +=      "<option value=\"105\">THORONDOR (Multi)</option>"
        builder +=      "</select>"
        builder +=  "</div>"

        /* Waypoint table */
        builder +=  "<div class=\"row\" id=\"waypoints-list\">"
        builder +=      "<table class=\"table\">"
        builder +=          "<thead><tr>"
        builder +=              "<th scope=\"col\">Lat/Long</th>"
        builder +=              "<th scope=\"col\">Distance</th>"
        builder +=              "<th scope=\"col\">Time</th>"
        builder +=          "</tr></thead>"
        builder +=          "<tbody id=\"waypoints-body\">"
        builder +=          "</tbody>"
        builder +=      "</table>"
        builder +=  "</div>"


        builder += "</div>"
        window.ui.#fillSidebar(header, builder);
    }

    listWaypoints(waypointsList, distance) {
        var listContainer = document.getElementById("waypoints-body");
        listContainer.innerHTML = null;

        var vehicleSpeed = document.getElementById("wp-vehicle-select").value;

        var builder = "";
        for (var i = 0; i < waypointsList.length; i++) {
            var lat = waypointsList[i].lat.toFixed(4);
            var lng = waypointsList[i].lng.toFixed(4);
            var dist = distance[i];
            var time = Math.round(dist/vehicleSpeed);
            builder += "<tr><td>" + lat + "," + lng + "</td><td>" + dist + "</td><td>" + time + "</td></tr>"
        }
        listContainer.innerHTML = builder;
    }

    clearMarkers() {
        map.markerLayer.clearLayers();
        map.lastMarker = null;
        map.targetMarker = null;
        document.getElementById("waypoints-body").innerHTML = "";
        map.removeWaypoints();
    }


    
}