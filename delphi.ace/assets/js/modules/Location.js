export default class Location{
    
    /**
     * Creates an Infrastructure object and sets its main variables
     * @param {Object} data - Individual Infrastructure object taken from the JSON file
     */
    constructor(latitude, longitude, altitude){
        this.latitude = latitude;
        this.longitude = longitude;
        this.altitude = altitude;
    }

    /**
     * Checks if the infrastructure is alive. This function does not move the object
     * @param {Object} data - New Individual Infrastructure object taken from the JSON file
     */
    update(latitude, longitude, altitude){
        this.latitude = latitude;
        this.longitude = longitude;
        this.altitude = altitude;
    }

    /**
     * Returns the location in Leaflet Format
     * @returns A Leaflet LatLng for the given Location Object
     */
    latLng(){
        return L.latLng(this.latitude, this.longitude);
    }
    

    
}