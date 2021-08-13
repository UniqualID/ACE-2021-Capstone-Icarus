export default class Communication{

    
    /**
     * Stream alerts
     **/
    constructor(){
        var webSocket = new WebSocket("ws://delphi.ace:3000/");
        webSocket.onopen = function (event) {
            webSocket.send("Starting Client Connection");
        };


        // Listen for messages
        webSocket.addEventListener('message', function (event) {
            window.ui.alert(event.data);
        });

        // Listen for possible errors
        webSocket.addEventListener('error', function (event) {
            window.ui.alert("Could not connect to the Delphi event stream");
        });
    }


}