!function(){window.App=Ember.Application.create()}(),function(){App.IndexController=Ember.Controller.extend({modes:["tree","view","form","code","text"],mode:"tree",name:"JSONEditor"})}(),function(){App.ApplicationAdapter=DS.FixtureAdapter}(),function(){Ember.Handlebars.helper("pretty-print",function(a,b,c){return console.log(arguments),JSON.stringify(a,b,c)})}(),function(){App.IndexRoute=Ember.Route.extend({model:function(){return{array:[1,2,3],"boolean":!0,"null":null,number:123,object:{a:"b",c:"d",e:"f"},string:"Hello World"}}})}(),function(){App.Router.map(function(){})}();