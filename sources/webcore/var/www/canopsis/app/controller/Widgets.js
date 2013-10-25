/*
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
*/
Ext.define('canopsis.controller.Widgets', {
	extend: 'Ext.app.Controller',

	stores: ['Widgets'],
	models: ['Event'],

	item_to_translate: ['title', 'fieldLabel', 'boxLabel', 'text', 'emptyText', 'header'],

	logAuthor: '[controller][Widgets]',

	init: function() {
		Ext.Loader.setPath('widgets', '/static/canopsis/widgets');
		Ext.Loader.setPath('widgets.thirdparty', '/static/widgets');

		this.store = this.getStore('Widgets');

		log.debug('parsing Widget store', this.logAuthor);

		if(!global.minimified) {
			this.store.on('load', function() {
				this.store.each(function(record) {
					var name = undefined;

					if(record.get('thirdparty')) {
						name = 'widgets.thirdparty.' + record.get('xtype') + '.' + record.get('xtype');
					}
					else {
						name = 'widgets.' + record.get('xtype') + '.' + record.get('xtype');
					}

					log.debug('loading ' + record.get('xtype') + ' (' + name + ')', this.logAuthor);
					Ext.require(name);
				}, this);

				this.clean_disabled_widget();

				//translate the store
				this.check_translate();

				// small hack
				Ext.Function.defer(function() {
					this.fireEvent('loaded');
				}, 1000, this);

			}, this, {single: true});
		}
		else {
			this.store.on('load', function() {
				this.clean_disabled_widget();
				this.fireEvent('loaded');
			}, this);
		}
	},

	clean_disabled_widget: function() {
		var records = [];

		this.store.each(function(record) {
			if(record.get('disabled') === true) {
				log.debug('Remove ' + record.get('xtype') + ' from widget store', this.logAuthor);
				records.push(record);
			}
		}, this);

		this.store.remove(records);
	},

	check_translate: function() {
		log.debug('Attempting to translate widget in store', this.logAuthor);

		this.store.each(function(record) {
			var options = record.get('options');

			if(options !== undefined) {
				for(var i = 0; i < options.length; i++) {
					this.translate(record.get('xtype'), options[i]);
				}
			}
		}, this);
	},

	//recursive translate function for widget records
	translate: function(xtype, data) {
		// for every item
		var me = this;

		Ext.Object.each(data, function(key, value) {
			if((key === 'items' || key === 'store' || key === 'data' || key === 'additional_field' || key === 'customForm' || key >= 0) && typeof(value) !== 'string') {
				me.translate(xtype, value);
			}

			if(Ext.Array.contains(me.item_to_translate, key)) {
				data[key] = _(value, xtype);
			}
		});
	}
});
