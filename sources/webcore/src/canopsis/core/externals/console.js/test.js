console = Console.create({
	_baseConsole: console
});

console.log("m1", console);
console.tags.add("t1");
console.tags.add("t2");
console.log("m2", console);
console.tags.remove("t1");
console.log("m3", console);
console.tags.flush();
console.log("m4", console);


console.tags.add("t3");
console.log("Object instantiated");

window.setInterval(function(){
	console.log("m5");
	if(Math.random() > 0.25)
		console.log("m6");
},1000);
