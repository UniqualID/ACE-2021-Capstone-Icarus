import EnemyPing from "./EnemyPing.js";

export default class Radar{
    
    constructor(map){
        this.dict = {}
        this.map = map
        console.log(this.map)
        console.log("Radar constructor")
    }

    update(data){
        // console.log("Radar update")

        var self = this;
        $.each(data, function(key, value){
            if (key in self.dict){
                self.dict[key].update(value)
            }else{
                self.dict[key] = new EnemyPing(key,value,self.map)
            }
        });
    }
}