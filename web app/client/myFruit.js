

function displayAll(){

    $.getJSON('Fruit.json', function(data){
        let arr = data

    console.log("Button activated \n"+ "#########################")
    var out = "Print all fruits: <br>";
    var i;
    if(arr.length==0){
        out += 'no fruit in the list';
    }
    for(i = 0; i< arr.length; i++){
        out += 'id: '+arr[i].id + ', name: '+arr[i].Name + ', price: '+arr[i].Price + ', Date: '+ arr[i].Date + '<br>'
    }
    document.getElementById("food").innerHTML = out;
    })

}