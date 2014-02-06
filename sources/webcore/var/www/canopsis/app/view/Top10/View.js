//need:app/lib/view/cwidget.js,app/store/Perfdatas.js,app/lib/view/cgrid.js
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
Ext.define('canopsis.view.Top10.View' , {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.top10',

	requires: [
		'canopsis.store.Perfdatas',
		'canopsis.lib.view.cgrid'
	],

	//Default options
	limit: 10,
	mfilter: undefined,
	sort: -1,
	time_window: 86400,
	threshold: undefined,
	threshold_direction: -1,

	show_source_type: true,
	show_last_time: false,
	show_metric: true,
	show_component: true,
	show_resource: true,
	show_value: false,
	show_unit: false,
	show_hr_value: true,
	show_percent: false,
	hideHeaders: false,
	show_position: true,
	expand: false,

	afterContainerRender: function() {
		this.store = Ext.create('canopsis.store.Perfdatas', {
			model: 'canopsis.model.Perfdata',

			autoLoad: false,

			proxy: {
				type: 'rest',
				url: '/perfstore/perftop',
				extraParams: {
					'limit': this.limit,
					'sort': this.sort,
					'mfilter': Ext.JSON.encode({
						'me': {'$in': [
							'cps_statechange_1',
							'cps_statechange_2'
						]}
					}),
					'output': true,
					'threshold': this.threshold,
					'threshold_direction': this.threshold_direction,
					'expand': this.expand,
					'percent': this.show_percent,
					'threshold_on_pct': this.threshold_on_pct,
					'report': this.reportMode || this.exportMode
				},
				reader: {
					type: 'json',
					root: 'data',
					totalProperty: 'total',
					successProperty: 'success'
				}
			}
		});

		this.columns = [{
			header: '',
			width: 25,
			align: 'center',
			sortable: false,
			renderer: function(value, metaData, record, rowIndex) {
				void(value, metaData, record);

				return '<b>' + (rowIndex + 1) + '</b>';
			}
		},{
			header: '',
			width: 25,
			align: 'center',
			sortable: false,
			renderer: function(value, metaData, record) {
				void(value, metaData);

				if(record.get('me') === 'cps_statechange_1') {
					return '<span class="icon icon-state-1" />';
				}
				else {
					return '<span class="icon icon-state-2" />';
				}
			}
		},{
			header: '',
			width: 25,
			sortable: false,
			renderer: function() {
				return "<span class='icon icon-mainbar-perfdata' />";
			}
		},{
			header: _('Component'),
			flex: 1,
			sortable: false,
			dataIndex: 'co'
		},{
			header: _('Resource'),
			flex: 1,
			sortable: false,
			dataIndex: 're'
		},{
			header: _('Most recurrent output'),
			flex: 2,
			sortable: false,
			dataIndex: 'output'
		},{
			header: _('Number of alerts'),
			width: 100,
			sortable: false,
			dataIndex: 'lv',
			align: 'right',
			renderer: function(value, metaData, record) {
				void(metaData);

				var me = this.cwidget();
				var unit = record.get('u');

				if(me.humanReadable) {
					return rdr_humanreadable_value(value, unit);
				}
				else {
					if(unit) {
						return value + ' ' + unit;
					}
					else {
						return value;
					}
				}
			}
		}];

		this.grid = Ext.create('canopsis.lib.view.cgrid', {
			model: 'Perfdata',
			store: this.store,

			opt_bar: true,
			opt_paging: false,
			opt_menu_delete: false,
			opt_bar_add: false,
			opt_menu_rights: false,
			opt_bar_search: false,

			opt_tags_search: false,
			opt_simple_search: false,

			opt_bar_customs: [{
				xtype: 'button',
				iconCls: 'icon-export',
				text: _('Export'),
				handler: function() {
					var params = [
						{name: 'csv', value: true},
						{name: 'fields', value: Ext.JSON.encode(['co', 're', 'me', 'output', 'lv'])}
					];

					for(var key in this.store.proxy.extraParams) {
						params.push({
							name: key,
							value: this.store.proxy.extraParams[key]
						});
					}

					console.log(params);

					getDataFromURL(this.store.proxy['url'], params);
				}.bind(this)
			}],

			hideHeaders: this.hideHeaders,

			opt_cell_edit: false,

			columns: this.columns,

			cwidget: function() {
				return this;
			}.bind(this)
		});

		this.wcontainer.removeAll();
		this.wcontainer.add(this.grid);

		this.ready();
	},

	doRefresh: function(from, to) {
		this.store.proxy.extraParams['time_window'] = to - from;
		this.store.proxy.extraParams['report'] = this.reportMode || this.exportMode;
		var url  = this.store.proxy['url'];

		if(this.store.proxy.extraParams['report']) {
			this.store.proxy['url'] = url + '/' + parseInt(from/1000) +'/' + parseInt(to/1000);
		}

		if(this.grid) {
			this.grid.store.load();
		}

		this.store.proxy['url'] = url;
	}
});
