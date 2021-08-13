import Communication from './modules/Communication.js';
import EnemyPing from './modules/EnemyPing.js';
import Radar from './modules/Radar.js';
import IADS from './modules/IADS.js';
import Infrastructure from './modules/Infrastructure.js';
import Map from './modules/Map.js';
import UI from './modules/UI.js';
import Vehicle from './modules/Vehicle.js';

var REFRESHTIME = 2000

class DelphiMap{

    #teams;
    #mission;

    constructor(data){

        /**
         * Detect the working environment. Interns can see management functions in the code, but can't run any of them from Grayspace.
         */
        if(location.hostname.match('mgmt')){
            this.mgmtMode = true;
            console.log('Running in management context');
        } else{
            this.mgmtMode = false;
        }
        {
            let bounds = L.latLngBounds(L.latLng(data.exercise.minimum_latitude, data.exercise.minimum_longitude) , L.latLng(data.exercise.maximum_latitude, data.exercise.maximum_longitude));
            window.map = new Map(6.75, 15, bounds);
            window.ui = new UI();


            console.log(window.map, 'HEREHEREHERE')


            /* Save mission info */
            this.#mission = data.exercise;

            /* Define the teams */
            this.#teams = ["Valinar", "Halcyon", "Malazan", "Gallifrey", "Civilian"];

            /* Construct objects for infrastructure */
            this.infrastructure = {};
            for(let i = 0; i < data.infrastructure.length; i++){
                this.infrastructure[data.infrastructure[i].id] = new Infrastructure(data.infrastructure[i]);
            }

            /* Construct objects for vehicles */
            this.vehicles = {};
            for(let i = 0; i < data.vehicles.length; i++){
                if((data.vehicles[i].role != "SAM") && (data.vehicles[i].role != "RADAR")){
                    this.vehicles[data.vehicles[i].id] = new Vehicle(data.vehicles[i]);
                }
            }

            /*Construct objects for IADS */
            this.iads = {};
            for(let i = 0; i < data.vehicles.length; i++){
                if((data.vehicles[i].role == "SAM") || (data.vehicles[i].role == "RADAR")){
                    this.iads[data.vehicles[i].id] = new IADS(data.vehicles[i]);
                }
            }

            this.alertServer = new Communication();

        }
        window.ui.showSidenav();
        {
            // let builder = "<p><div class=\"text-center\"><img width=\"100px\"src=\"assets/img/logo_login.png\" class=\"img-fluid\"/></div></p>"
            // builder += "<b>What's New?!</b><p><ul><li><b>News Feed and Alerts</b> - Delphi provides alerts to keep you in the loop. The news feed keeps track of your mission.</li>";
            // builder += "<li><b>JavaScript ES6</b> - Delphi now uses multiple objects to provide the best SA envrionment which allows modification for your needs.</li>";
            // builder += "<li><b>Bug Fixes</b> - Delphi can handle errors better than ever, and sort of works.</li>";
            // builder += "<li><b>API Integrations</b> - Delphi is rolling out intgrations for your favourite Battlespace services!</li>";
            // builder += "<li><b>Keyboard shortcuts</b> - Press [c] to copy coordinates to the clipboard. Press [m] to add a marker to your map.</li>";
            // builder += "<li><b>Design your own map!</b> - You now have the option to upload your own JSON file to Delphi! Know your enemy!</li>";
            // builder += "</ul> <p><i> It's not the best choice, it's the Iron choice!</i></p>";
            // window.ui.createModal("Welcome to Delphi!", builder);
        }
       
    }

    updateData(data){
        /*Update the Vehciles */ 
        for(let i = 0; i < data.vehicles.length; i++){
            if((data.vehicles[i].role != "SAM") && (data.vehicles[i].role != "RADAR")){
                if((data.vehicles[i].id in this.vehicles)){
                    this.vehicles[data.vehicles[i].id].update(data.vehicles[i]);
                } else{
                    this.vehicles[data.vehicles[i].id] = new Vehicle(data.vehicles[i]);
                }
            } else{
                if((data.vehicles[i].id in this.iads)){
                    this.iads[data.vehicles[i].id].update(data.vehicles[i]);
                } else{
                    this.iads[data.vehicles[i].id] = new IADS(data.vehicles[i]);
                }
            }
        }
        for(let i = 0; i < data.infrastructure.length; i++){
            if((data.infrastructure[i].id in this.infrastructure)){
                this.infrastructure[data.infrastructure[i].id].update(data.infrastructure[i]);
            } else{
                this.infrastructure[data.infrastructure[i].id] = new Infrastructure(data.infrastructure[i]);
            }
        }

    }
    importJSON(){
        var myFile = $('#jsonupload').prop('files');
        var fileReader = new FileReader();
        fileReader.onload = function () {
            var data = fileReader.result;  // data <-- in this var you have the file data in Base64 format
            this.updateData(JSON.parse(data));
        }.bind(this);
        fileReader.readAsText($('#jsonupload').prop('files')[0]);
    }

    getTeams(){
        return this.#teams;
    }

    getMission(){
        return this.#mission;
    }

    isMGMT(){
        return this.mgmtMode;
    }

}


/**
* Initiates the script
* @returns {void}
*/
function init() {
    // let runtime =  location.hostname;
    $('#delphiModal').on('hidden.bs.modal', function () {
        $('#delphiModal').remove();
    });

    
    let runtime = "delphi.ace"
    if(runtime.match('delphi.mgmt')){
        var data_file = "http://delphi.mgmt" + "/data.json";
    } else{
        var data_file = "https://" + runtime + "/" + findGetParameter("auth") + ".json";
    }

    $.getJSON(data_file , function( data ) {
        // var delphiMap;
        window.delphiMap = new DelphiMap(data);
        window.enemyRadar = new Radar(window.map)
        setInterval(function(){
        $.getJSON(data_file, function( data ) {
            if (data!=""){ //If a Null result is not returned then show the result of the search
                window.delphiMap.updateData(data);
                }
            }).fail(function(){
                window.ui.alert("Delphi is offline. Cannot reach API.")
            });
        }, REFRESHTIME);

        setInterval(function(){
        $.getJSON("/assets/radar/radar.json", function(data){
            enemyRadar.update(data)
        });
        }, REFRESHTIME);

        $.getJSON("template.json", function(data){
            window.delphiMap.updateData(data)
        });
    }).fail(function(){
        location.reload();
    });

    

    // $(document).ready(function(){
    // $.getJSON("/assets/radar/radar.json", function(data){
    //     window.enemyRadar = new Radar(data, window.delphiMap)
    //     setInterval(function(){
    //     $.getJSON("/assets/radar/radar.json", function(data){
    //         enemyRadar.update(data)
    //     });
    //     }, 5000);
    // });
    // });
}


function findGetParameter(parameterName) {
    var result = null,
        tmp = [];
    location.search
        .substr(1)
        .split("&")
        .forEach(function (item) {
          tmp = item.split("=");
          if (tmp[0] === parameterName) result = decodeURIComponent(tmp[1]);
        });
    return result;
}
   
init();
  

