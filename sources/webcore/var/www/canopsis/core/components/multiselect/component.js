/*
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
*/

define([
	'ember',
	'app/application'
], function(Ember, Application) {

	Ember.Widgets.MultiSelectComponent3 = Ember.Widgets.MultiSelectComponent.extend({
		selectionItemView: Ember.Widgets.MultiSelectOptionView,
		select:0,

	  	keyDown: function(event) {
	  		// space = 32 ; enter = 13
	  		if (this.canCreate &&
		  	    event.keyCode === 13) {
				var query = this.get("query");
				if (!Ember.isEmpty(query)) {
			  	    var selections = this.get("selections");
			  	    query = this.get("query");
			  	    selections.pushObject({ name : query});//this.get("content")[0]);
				    Ember.set(this,'query',"");
				}
		  	}
		},

		init: function() {
		  	this._super();
		  	var SelectOption = ( this.get("select") === 0 )? "MultiSelectOptionView" : "MultiSelectOptionViewMY";
		  	this.set("selectionItemView",  Ember.Widgets[SelectOption]);
		  	//this.filteredContent;
		  	content = this.get( 'content' );
		  	var selections = this.get( "selections" );

		  //	var selectionsName = this.getNamesArray(selections);

		 	var newContent = Ember.copy(content , true);
	  		this.set("content" , newContent);

		  	var contentNames = this.getNamesArray( content );
		  	var selectionNOtInContent = this.filter( contentNames, "selections", false );

		  	var NameToRefOnContent = this.filter( selections , "content", true, "name" ) ;
		  	if ( selections ) {
				while ( selections.length > 0 ) {
					selections.pop();
				}
			}

			for ( var i=0 ; i < NameToRefOnContent.length ; i++ ) {
		    	selections.pushObject(NameToRefOnContent[i]);
			}

			for ( i=0 ; i < selectionNOtInContent.length ; i++ ) {
		    	selections.pushObject({ name: selectionNOtInContent[i]});
			}
		},

		getNamesArray: function(selections) {
			var selectionsName = [];
			for ( i=0 ; i < selections.length ; i++ ) {
				var selection = selections[i];
					selectionsName.push( selection.name );
			}
			return selectionsName;
		},

		filter: function( ToFilterWith , proprietyName , isON , fieldToFilter ) {
			var tabFiltered = this.get( proprietyName ).filter( function ( content ) {
				var toTestWith = ( fieldToFilter )? content[fieldToFilter] : content ;
				var keep = ToFilterWith.contains( toTestWith );
		    	return (isON)? keep : !keep;
		    });
		    return tabFiltered;
		},

		convertDictToArray:function(item) {
		    var fieldsArray = Ember.A();
			for (var attr in item) {
			//if option isn't in current model's options
				if (item.hasOwnProperty(attr)) {
				// I keep it for remenber how to keep track of true reference value
					var newObject = Ember.Object.create({value : item[attr] , field : attr });
					//newObject.addObserver('value',item, this.fooDidChange);

					fieldsArray.pushObject(newObject);
					console.log ( "Added "+ attr + " = " + item[attr] +" newObject = " + newObject[attr]);
				}
			}
			return fieldsArray;
		},

		modalShow: function(item) {
		    Ember.Widgets.ModalComponent.popup({
			    targetObject: this,
			    confirm: 'modalConfirm',
			    cancel: 'modalCancel',
			    content: this.convertDictToArray(item),
			    test:item,
			    //templateName: 'custom-modal-content',
			    contentViewClass: Ember.View.extend({
				    tagName: '',
				    template: Ember.Handlebars.compile('{{log test}} {{#each attr in content}} <div><span class="label" style="color:black"><label>{{attr.field}}:</label></span>{{input valueBinding="attr.value" }}</div> {{/each}}'),
				}),
				confirm: function(){
				    	console.log("Modal Confirm!",item);
				    	var fieldsArray = this.get("content");
				    	for (var i = 0; i < fieldsArray.length; i++) {
							Ember.set(item, fieldsArray[i].field, fieldsArray[i].value);
						}
					},


				cancel: function(){  console.log("Modal Cancel!") }
			});
		},
		actions: {
		 	removeSelectItem: function(item) {
				return this.removeSelectItem(item);
			},
			edit : function(item) {
				this.modalShow(item);
			}
		}
	});

	Ember.Widgets.MultiSelectOptionViewMY = Ember.Widgets.MultiSelectOptionView.extend({
	  	templateName: 'multi_select_itemMY',
	}),

	Ember.TEMPLATES["multi_select_itemMY"] = Ember.Handlebars.template(function anonymous(Handlebars,depth0,helpers,partials,data) {
		this.compilerInfo = [4,'>= 1.0.0'];
		helpers = this.merge(helpers, Ember.Handlebars.helpers); data = data || {};
		  var buffer = '', stack1, escapeExpression=this.escapeExpression;


		data.buffer.push("<div>");
		stack1 = helpers._triageMustache.call(depth0, "view.label", {hash:{},hashTypes:{},hashContexts:{},contexts:[depth0],types:["ID"],data:data});
		if(stack1 || stack1 === 0) { data.buffer.push(stack1); }
		data.buffer.push("</div><a class=\"ember-select-search-choice-close\" href=\"#\"");
		data.buffer.push(escapeExpression(helpers.action.call(depth0, "removeSelectItem", "view.content", {hash:{},hashTypes:{},hashContexts:{},contexts:[depth0,depth0],types:["STRING","ID"],data:data})));
		data.buffer.push(">×</a>");

		data.buffer.push("<div><button type=\"button\" class=\"btn btn-primary\"");
		data.buffer.push(escapeExpression(helpers.action.call(depth0, "edit", "view.content", {hash:{},hashTypes:{},hashContexts:{},contexts:[depth0,depth0],types:["STRING","ID"],data:data})));
		data.buffer.push(">EDIT</button></div>");

		return buffer;
	});

    Ember.Handlebars.helper( 'multi-select-component3', Ember.Widgets.MultiSelectComponent3 );

    return Application.ComponentExpandableaddbuttonComponent;
});