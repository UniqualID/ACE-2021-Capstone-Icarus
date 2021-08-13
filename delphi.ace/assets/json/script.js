var smth = {}
smth["infrastructure"] = []
smth["vehicles"] = []
function handleFormSubmit() {
  //event.preventDefault();
  
  const data = new FormData(event.target);
  
  const formJSON = Object.fromEntries(data.entries());
  formJSON["id"] = parseInt(formJSON["id"])
  formJSON.hitpoints = parseInt(formJSON.hitpoints)
  formJSON.longitude = parseFloat(formJSON.longitude)
  formJSON.latitude = parseFloat(formJSON.latitude)
  
  // for multi-selects, we need special handling
  // formJSON.snacks = data.getAll('snacks');
  switch (formJSON.type){
    case "Airbase":
    case "Depot":
    case "Datacenter":
    case "Fort":
    case "WMD-PRODUCTION":
    case "ENCAMPMENT":
    case "AnimalHospital":
    case "EMBASSY":
    case "NETWORK-NODE":
    case "INSTALLATION":
    case "MYSTERY":
      smth["infrastructure"].push(formJSON);
      break;
    default:
      smth["vehicles"].push(formJSON)
  }
  const results = document.querySelector('.results pre');
  // results.innerText = JSON.stringify(formJSON, null, 2);
  results.innerText = JSON.stringify(smth, null, 2)
}


function download(content, fileName, contentType) {
    var a = document.createElement("a");
    var file = new Blob([content], {type: contentType});
    a.href = URL.createObjectURL(file);
    a.download = fileName;
    a.click();
}

function save(){
  download(JSON.stringify(smth, null, 2), "newassets.json", "text/plain")
}
const form = document.querySelector('.contact-form');
//form.addEventListener('submit', handleFormSubmit);
