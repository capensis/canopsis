//need:app/lib/view/cwidget.js,widgets/stepeue/feature.js
/*
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------
*/
Ext.define('widgets.stepeue.stepeue', {
	extend: 'canopsis.lib.view.cwidget',
	alias: 'widget.stepeue',

	requires: [
		'widgets.stepeue.feature'
	],

	logAuthor: '[widget][stepeue]',
	scroll: true,
	useScreenShot: true,

	initComponent: function() {
		log.debug('initialization of the eue\'s widget', this.logAuthor);

		pnl = Ext.create('Ext.Panel', {
			xtype: 'panel',
			width: '100%',
			height: '100%'
		});

		pnl.on('afterrender', function() {
			pnl.setLoading(true, true);
		});

		this.callParent(arguments);
		this.wcontainer.add(pnl);
		this.wcontainer.setLoading(true, true);
	},

	destroyObject: function() {
		log.debug('object destroyed', this.logAuthor);

		for(var i = 0; i < this.features.length; i++) {
			this.features[i].destroyFeature();
			Ext.destroyMembers(this.features[i]);
		}
	},

	doRefresh: function(from, to) {
		log.debug('do Refresh', this.logAuthor);

		if(this.features != undefined) {
			this.destroyObject();
		}

		if(this.nodes.length == 0) {
			this.buildHtml();
		}
		else {
			this.urlPerfStore = this.makeUrl(from, to);
			this.last_from = to;
			this.features = new Array();

			for(var i = 0; i < this.nodes.length; i++) {
				var feature = Ext.create('widgets.stepeue.feature');  //creation of empty feature
				this.features.push(feature);
			}
		}

		this.buildHtml();
	},

	makeUrl: function(from, to) {
		var url = '/perfstore/values';

		if(!to) {
			url += '/' + parseInt(from/1000);
		}

		if(from && to) {
			url += '/' + parseInt(from/1000) + '/' + parseInt(to/1000);
		}

		return url;
	},

	buildHtml: function() {
		if (this.nodes.length != 0) {
			//If there are nodes it means that there are features to print, otherwise we print "no features selected"
			log.debug('Build the view of the widget', this.logAuthor);

			var listItems = new Array();
			var me = this;

			for(var i = 0; i < this.nodes.length; i++) { 
				//each nodes mean each feature
				var title = this.nodes[i].split('.')[5]; //title of feature
				var object = {
					title: title,
					id: me.wcontainer.id + ':feature:' + i,
					layout: 'fit',
					border: false,
					listeners: {
						activate: function(tab) {
							log.debug('activated');

							var idString = tab.id.split(':');
							var id = idString[2]; //id of the features

							me.features[id].init(me.nodes[id], me, tab); //we build the feature view after initialization

						}
					}
				};

				listItems.push(object); 
			}

			//if the array listItems has one item it means that there is no tabs otherwise it means that there are.
			if(listItems.length == 1) {
				listItems[0].xtype = 'panel';

				listItems[0].listeners.afterrender = function() {
					log.debug('after render', this.logAuthor);
					me.features[0].init(me.nodes[0], me, this);
				};

				this.content = Ext.create('Ext.Panel', {
					xtype: 'panel',
					items: listItems,
					layout: 'fit',
					border: false
				});
			}
			else {
				var tabsPanel = Ext.create('Ext.tab.Panel', {
					xtype: 'panel',
					width: '100%',
					height: '100%',
					items: listItems,
					border: false
				});

				this.content = Ext.create('Ext.Panel', {
					xtype: 'panel',
					items: tabsPanel,
					layout: 'fit',
					border: false
				});
			}

			this.wcontainer.removeAll();
			this.wcontainer.add(this.content);
		}
		else {
			this.wcontainer.removeAll();

			this.wcontainer.add({
				html: 'no feature selected'
			});
		}
	}
});
