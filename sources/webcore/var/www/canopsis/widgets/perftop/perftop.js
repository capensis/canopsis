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
Ext.define('widgets.perftop.perftop' , {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.perftop',

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
		if(this.mfilter && this.mfilter !== '') {
			this.store = Ext.create('canopsis.store.Perfdatas', {
				model: 'canopsis.model.Perfdata',

				autoLoad: false,

				proxy: {
					type: 'rest',
					url: '/perfstore/perftop',
					extraParams: {
						'limit': this.limit,
						'sort': this.sort,
						'mfilter': this.mfilter,
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

			this.columns = [];

			if(this.show_position) {
				this.columns.push({
					header: '',
					width: 25,
					align: 'center',
					sortable: false,
					renderer: function(value, metaData, record, rowIndex) {
						void(value, metaData, record);

						return '<b>' + (rowIndex + 1) + '</b>';
					}
				});
			}

			if(this.show_source_type) {
				this.columns.push({
					header: '',
					width: 25,
					sortable: false,
					renderer: function() {
						return "<span class='icon icon-mainbar-perfdata' />";
					}
				});
			}

			if(this.show_component) {
				this.columns.push({
					header: _('Component'),
					flex: 1,
					sortable: false,
					dataIndex: 'co'
				});
			}

			if(this.show_resource) {
				this.columns.push({
					header: _('Resource'),
					flex: 1,
					sortable: false,
					dataIndex: 're'
				});
			}

			if(this.show_metric) {
				this.columns.push({
					header: _('Metric'),
					flex: 2,
					sortable: false,
					dataIndex: 'me'
				});
			}


			if(this.show_hr_value) {
				this.columns.push({
					header: _('Value'),
					width: 100,
					sortable: false,
					dataIndex: 'lv',
					align: 'right',
					renderer: function(value, metaData, record) {
						void(metaData);

						var me = this.cwidget;
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
				});
			}

			if(this.show_value) {
				this.columns.push({
					header: _('Value'),
					width: 100,
					sortable: false,
					dataIndex: 'lv',
					align: 'right'
				});
			}

			if(this.show_percent) {
				this.columns.push({
					header: _('Percent'),
					width: 100,
					sortable: false,
					dataIndex: 'pct',
					align: 'right',
					renderer: function(value) {
						if(Ext.isNumber(value)) {
							if(value === -1) {
								return _('N/A');
							}
							else {
								return value + '%';
							}
						}
						else {
							return _('N/A');
						}
					}
				});
			}

			if(this.show_unit) {
				this.columns.push({
					header: _('Unit'),
					width: 45,
					sortable: false,
					dataIndex: 'u',
					align: 'center'
				});
			}

			if(this.show_last_time) {
				this.columns.push({
					header: _('Last time'),
					width: 130,
					dataIndex: 'lts',
					align: 'center',
					renderer: rdr_tstodate
				});
			}

			this.grid = Ext.create('canopsis.lib.view.cgrid', {
				model: 'Perfdata',
				store: this.store,

				opt_bar: false,
				opt_paging: false,
				opt_menu_delete: false,
				opt_bar_add: false,
				opt_menu_rights: false,
				opt_bar_search: false,

				opt_tags_search: false,
				opt_simple_search: false,

				hideHeaders: this.hideHeaders,

				opt_cell_edit: false,

				columns: this.columns,
				cwidget: this
			});

			this.wcontainer.removeAll();
			this.wcontainer.add(this.grid);
		}

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
