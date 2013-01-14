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

widget_weather_template = Ext.create('Ext.XTemplate',
		'<table class="weather-table">',
			'<tr>',
				'<td style="vertical-align: top;" colspan=3>',
					'<span class="weather-title">{title}</span>',
					'<span class="weather-ts">{event_ts}</span></br>',
					'{output}',
				'</td>',
				'<td style="width: 14px;" id="{id}-edit_td">',
					'<tpl if="admin == true && derogation == true">',
						'<div class="icon icon-edit weather-clickable" id="{id}-edit_button"></div>',
					'</tpl>',
				'</td>',
				'<td style="width: 20%; position:relative">',
					'<tpl if="percent != undefined ">',
						'<div class="weather-percent">{percent}%</div>',
					'</tpl>',
					'<img class="weather-image" src="{icon_src}"/>',
				'</td>',
			'</tr>',
			'<tr>',
				'<td rowspan=2 style="width: 90px;">',
					'<tpl if="button_text != undefined">',
						'<button class="weather-button"  type="button" id="{id}-button">{button_text}</button>',
					'</tpl>',
				'</td>',
				'<td rowspan=2 style="width: 30px;">',
					'<tpl if="alert_icon" != undefined">',
						'<img src="{alert_icon}">',
					'</tpl>',
				'</td>',
				'<td rowspan=2>',
					'<tpl if="alert_msg" != undefined">',
						'<p class="weather-alert-message" id="{id}-alert_message">{alert_msg}</p>',
					'</tpl>',
				'</td>',
				'<td></td>',
				'<td><center>{legend}</center></td>',
			'</tr>',
			'<tr>',
				'<td></td>',
				'<td></td>',
			'</tr>',
		'</table>',
		{compiled: true}
	);

widget_weather_simple_template = Ext.create('Ext.XTemplate',
		'<table class="weather-table" style="height:100%;vertical-align:middle;">',
			'<tr>',
				'<td style="width:25%" class=""></td>',
				'<td style="width:25%" class=""><img class="weather-image" src="{icon_src}"></td>',
				'<td style="width:30%;font-size:{title_font_size}px" class="" id="{id}-title">',
					'<div><span>{title}</span></div>',
				'</td>',
				'<td style="width:15%" class=""></td>',
		'</tr>',
		{compiled: true}
	);

Ext.define('widgets.weather.brick' , {
	extend: 'Ext.Component',
	alias: 'widget.weather.brick',

	logAuthor: '[widget][weather][brick]',

	brick_number: undefined,
	iconSet: 1,
	icon_on_left: false,
	state_as_icon_value: false,
	bg_color: '#FFFFFF',

	display_report_button: false,
	display_derogation_icon: false,

	simple_display: false,
	title_font_size: 14,

	alert_icon_basedir: 'widgets/weather/icons/alert/',
	alert_icon_name: ['workman.png','slippery.png','alert.png'],
	info_weather_icon: 'widgets/weather/icons/info-icon.png',

	helpdesk: undefined,
	nodeId: undefined,
	component_name: undefined,

	fullscreenMode: false,

	initComponent: function() {
		log.debug(' + Initialize brick ' + this.data._id, this.logAuthor);
		if (this.bg_color) {
			if (this.bg_color.indexOf('#') == -1)
				this.bg_color = '#' + this.bg_color;

			this.style = {'background-color': this.bg_color};
		}

		//log.dump(this.data)
		this.event_type = this.data.event_type;
		this.component = this.data.component;
		this.resource = this.data.resource;

		this.callParent(arguments);
	},

	afterRender: function() {
		//------------------build widget base config--------------
		if (this.simple_display)
			this._html_template = widget_weather_simple_template;
		else
			this._html_template = widget_weather_template;

		this.widget_base_config = {
			id: this.id,
			title_font_size: this.title_font_size,
			derogation: !this.fullscreenMode
		};

		//title

		if (this.data.display_name) {
			this.widget_base_config.title = this.data.display_name;
		}else {
			if (this.component)
				this.widget_base_config.title = this.component;
			else
				this.widget_base_config.title = 'Unknown';
		}

		var linkUrl = this.formatLink();

		if (this.fullscreenMode && linkUrl)
			this.widget_base_config.title = '<a href="' + linkUrl + '" target="_newtab">' + this.widget_base_config.title + '</a>';

		//icons

		if (this.icon_on_left) {
			this.widget_base_config.first_panel_float = 'right';
			this.widget_base_config.second_panel_float = 'left';
		}else {
			this.widget_base_config.first_panel_float = 'left';
			this.widget_base_config.second_panel_float = 'right';
		}

		//check ressource admin
		if (global.accountCtrl.check_right(this.data, 'w'))
			this.widget_base_config.admin = true;

		//----------------------build html------------------------
		if (this.data) {
			if (!this.exportMode)
				this.build(this.data);
		}else {
			this.buildEmpty();
		}

		//-----------------------get element----------------------
		this.edit_button = this.getEl().getById(this.id + '-edit_button');
		//-----------------------bindings-------------------------
		var report_button = this.getEl().getById(this.id + '-button');
		if (report_button)
			report_button.on('click', this.report_issue, this);

		var clickable_title = this.getEl().getById(this.id + '-title');
		if (clickable_title && (this.external_link || this.link) && !this.fullscreenMode) {
			clickable_title.addCls('weather_clickable');
			clickable_title.on('click', this.externalLink, this);
		}
		if (this.widget_base_config.admin && this.display_derogation_icon && this.edit_button) {
			var output = this.getEl().getById(this.id + '-edit_td');
			if (output) {
				output.hover(
					function() {this.edit_button.fadeIn()},
					function() {this.edit_button.fadeOut()},
					this
				);
			}

			if (this.edit_button) {
				this.edit_button.on('click', function() {
					if (!this.data.rk) {
						global.notify.notify(_('Information not found'), _("Please wait a moment, some informations aren't availables"), 'info');
					}else {
						var name = this.component;
						if (this.resource)
							name += ' - ' + this.resource;
						global.derogationCtrl.derogate(this.data.rk, name, true);
					}
				},this);
			}
		}

		//Hack for removing scrolling bar on ie
		this.getEl().parent().setStyle('overflow-x','hidden')
	},

	build: function(data) {
		log.debug('  +  Build html for ' + data._id, this.logAuthor);

		var widget_data = {
			legend: rdr_elapsed_time(data.last_state_change, true),
			event_ts: rdr_tstodate(data.timestamp, true)
		};

		if (data.output && data.output != '')
			widget_data.output = data.output;

		if (data.event_type == 'selector') {
			var icon_value = 100 - (data.state / 4 * 100);
			widget_data.icon_src = this.getIcon(icon_value);
		}else {
			if (this.state_as_icon_value || this.selector) {
				if (!this.selector) {
					var icon_value = 100 - (data.state / 4 * 100);
					widget_data.icon_src = this.getIcon(icon_value);
				}else {
					log.debug('  +  This brick is using its selector state as icon', this.logAuthor);
					var icon_value = 100 - (this.selector.state / 4 * 100);
					widget_data.icon_src = this.getIcon(icon_value);
				}
			}else {
				if (data.perf_data_array[0])
					widget_data.icon_src = this.getIcon(data.perf_data_array[0].value);
				else
					widget_data.icon_src = this.info_weather_icon;
			}
			if (data.perf_data_array)
				widget_data.percent = data.perf_data_array[0].value;
		}

		//----------------alert && derog-------------
		if (this.display_report_button)
			widget_data.button_text = _('Report issue');

		if (this.data.alert_msg)
			widget_data.alert_msg = this.data.alert_msg;
		else
			widget_data.alert_msg = "&nbsp;"

		if (this.data.alert_icon != undefined)
			widget_data.alert_icon = this.alert_icon_basedir + this.alert_icon_name[this.data.alert_icon];

		var config = Ext.Object.merge(widget_data, this.widget_base_config);
		var _html = this._html_template.applyTemplate(config);
		this.getEl().update(_html);
	},


	buildReport: function(data) {
		log.debug(' + Build html report for ' + this.event_type + ' ' + this.component + ':', this.logAuthor);
		var widget_data = {	};

		if (data && data.values.length > 0) {
			var timestamp = data.values[0][0];
			var nb_points = data.values.length;
			var last_timestamp = data.values[nb_points - 1][0];
			var last_value = data.values[nb_points - 1][1];

			if (this.event_type == 'selector' || this.selector_state_as_icon_value) {
				var state = demultiplex_cps_state(last_value).state;
				log.debug(' + State of ' + this.component + ' is: ' + state, this.logAuthor);
				log.debug(' + ' + nb_points + ' points returned by server', this.logAuthor);
				log.debug('  +  First value ts: ' + timestamp, this.logAuthor);
				log.debug('  +  Last value ts: ' + last_timestamp, this.logAuthor);

				var icon_value = 100 - (state / 4 * 100);
				widget_data.icon_src = this.getIcon(icon_value);
				widget_data.output = _('State on') + ' ' + rdr_tstodate(last_timestamp / 1000);
			}else {
				var cps_pct_by_state_0 = last_value;
				widget_data.percent = cps_pct_by_state_0;
				widget_data.icon_src = this.getIcon(cps_pct_by_state_0);
				widget_data.output = _('SLA on') + ' ' + rdr_tstodate(last_timestamp / 1000);
			}
		} else {
			widget_data.icon_src = this.info_weather_icon;
			widget_data.output = _('No data available');
		}

		var config = Ext.Object.merge(widget_data, this.widget_base_config);
		var _html = this._html_template.applyTemplate(config);
		this.getEl().update(_html);
	},

	buildEmpty: function() {
		log.debug('  +  Build empty brick for ' + this.event_type + ' ' + this.component, this.logAuthor);
		var widget_data = {
			output: _('No data for the selected information'),
			icon_src: this.info_weather_icon
		};

		var config = Ext.Object.merge(widget_data, this.widget_base_config);
		var _html = this._html_template.applyTemplate(config);
		this.getEl().update(_html);
	},

	report_issue: function() {
		var config = {
			_component: this.component,
			display_name: this.data.display_name,
			referer: this.data.rk,
			title: _('Report issue for ') + this.event_type + ' ' + this.component,
			renderTo: Ext.getCmp('main-tabs').getActiveTab().id,
			helpdesk: this.helpdesk
		};

		var popup = Ext.create('widgets.weather.report_popup', config);
		popup.show();
	},

	externalLink: function() {
		log.debug(' + Clicked on title, follow specified link', this.logAuthor);
		if (this.link) {
			if (this.link.indexOf('http://') != -1 || this.link.indexOf('www.') != -1 || this.link.indexOf('https://') != -1) {
				if (this.link.indexOf('http://') == -1 && this.link.indexOf('https://') == -1)
					this.link = 'http://' + this.link;
				window.open(this.link, '_newtab');
			}else {
				Ext.getStore('Views').load({
					scope: this,
					callback: function(records, operation, success) {
						store = Ext.getStore('Views');
						var record = store.findExact('crecord_name', this.link);
						if (record != -1) {
							record = store.getAt(record);
							if (!this.fullscreenMode) {
								global.tabsCtrl.open_view({ view_id: record.get('_id'), title: _(record.get('crecord_name')) });
							}else {
								var url = Ext.String.format('http://{0}/static/canopsis/display_view.html?view_id={1}&auth_key={2}',
									$(location).attr('host'),
									record.get('_id'),
									global.account.authkey
								);
								window.open(url, '_newtab');
							}
						}else {
							global.notify.notify('Link is not valid', 'The specified link does not match any view or URL', 'info');
						}
					}
				});
			}
		}
	},

	//fast hack for freeze, open link in tab, will be changed in develop
	formatLink: function() {
		if (this.link) {
			if (this.link.indexOf('http://') != -1 || this.link.indexOf('www.') != -1 || this.link.indexOf('https://') != -1) {
				if (this.link.indexOf('http://') == -1 && this.link.indexOf('https://') == -1)
					this.link = 'http://' + this.link;
				return this.link;
			}else {
				store = Ext.getStore('Views');
				var record = store.findExact('crecord_name', this.link);
				if (record != -1) {
					record = store.getAt(record);
					var url = Ext.String.format('http://{0}/static/canopsis/display_view.html?view_id={1}&auth_key={2}',
						$(location).attr('host'),
						record.get('_id'),
						global.account.authkey
					);
					return url;
				}else {
					log.debug('Link is not valid', 'The specified link does not match any view or URL', 'info');
					return undefined;
				}
			}
		}else {
			return undefined;
		}
	},

	getIcon: function(value) {
		value = Math.floor(value / 10) * 10;
		return 'widgets/weather/icons/set' + this.iconSet + '/' + value + '.png';
	}
});
