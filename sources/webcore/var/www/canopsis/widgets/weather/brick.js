//need:widgets/weather/report_popup.js
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

widget_weather_template = Ext.create('Ext.XTemplate',
	'<table class="weather-table">',
		'<tr>',
			'<td style="vertical-align: top;" colspan=3>',
				'<span class="weather-title" id="{id}-title">{title}</span>',
				'<span class="weather-ts">{event_ts}</span></br>',
				'{output}',
			'</td>',
			'<td style="width: 14px;" id="{id}-edit_td">',
				'<tpl if="admin == true && derogation == true && exportMode == false">',
					'<div class="icon icon-edit weather-clickable" id="{id}-edit_button"></div>',
				'</tpl>',
			'</td>',
			'<td class="weather-td-image" style="width: 20%;">',
				'<div class="weather-relative">',
					'<tpl if="percent != undefined ">',
						'<div class="weather-percent">{percent}%</div>',
					'</tpl>',
					'<img class="weather-image" src="{icon_src}"/>',
				'</div>',
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
			//'<td><center>&nbsp;</center></td>',
		'</tr>',
		'<tr>',
			'<td></td>',
			'<td></td>',
		'</tr>',
	'</table>',
	{compiled: true}
);

widget_weather_template_left = Ext.create('Ext.XTemplate',
	'<table class="weather-table">',
		'<tr>',
			'<td class="weather-td-image" style="width: 20%;">',
				'<div class="weather-relative">',
					'<tpl if="percent != undefined ">',
						'<div class="weather-percent">{percent}%</div>',
					'</tpl>',
					'<img class="weather-image" src="{icon_src}"/>',
				'</div>',
			'</td>',
			'<td style="vertical-align: top;" colspan=3>',
				'<span class="weather-title">{title}</span>',
				'<span class="weather-ts">{event_ts}</span></br>',
				'{output}',
			'</td>',
			'<td style="width: 14px;" id="{id}-edit_td">',
				'<tpl if="admin == true && derogation == true && exportMode == false">',
					'<div class="icon icon-edit weather-clickable" id="{id}-edit_button"></div>',
				'</tpl>',
			'</td>',
		'</tr>',
		'<tr>',
			'<td><center>{legend}</center></td>',
			//'<td></td>',
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
		'</tr>',
		'<tr>',
			'<td></td>',
			'<td></td>',
		'</tr>',
	'</table>',
	{compiled: true}
);

widget_weather_simple_template = Ext.create('Ext.XTemplate',
	'<table class="weather-table" style="vertical-align:middle;">',
		'<tr>',
			'<td style="width:25%" class=""></td>',
			'<td class="weather-td-image" style="width:25%;">',
				'<div class="weather-relative">',
					'<img class="weather-image" src="{icon_src}">',
					'<tpl if="percent != undefined ">',
						'<div class="weather-percent">{percent}%</div>',
					'</tpl>',
				'</div>',
			'</td>',
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

	requires: [
		'widgets.weather.report_popup'
	],

	logAuthor: '[widget][weather][brick]',

	brick_number: undefined,
	iconSet: 1,
	icon_on_left: false,
	state_as_icon_value: false,
	bg_color: '#FFFFFF',

	display_name: undefined,
	display_report_button: false,
	display_derogation_icon: false,

	hide_title: false,
	simple_display: false,
	title_font_size: 14,

	alert_icon_basedir: 'widgets/weather/icons/alert/',
	alert_icon_name: ['workman.png', 'slippery.png', 'alert.png'],
	info_weather_icon: 'widgets/weather/icons/question.png',

	helpdesk: undefined,
	nodeId: undefined,
	component_name: undefined,

	fullscreenMode: false,

	initComponent: function() {
		log.debug(' + Initialize brick ' + this.data._id, this.logAuthor);

		if(this.bg_color) {
			if(this.bg_color.indexOf('#') === -1) {
				this.bg_color = '#' + this.bg_color;
			}

			this.style = {'background-color': this.bg_color};
		}

		this.event_type = this.data.event_type;
		this.component = this.data.component;
		this.resource = this.data.resource;

		this.callParent(arguments);

		this.on('resize', this.onResize, this);
	},

	afterRender: function() {
		// build widget base config

		if(this.simple_display) {
			this._html_template = widget_weather_simple_template;
		}
		else {
			if(this.icon_on_left) {
				this._html_template = widget_weather_template_left;
			}
			else {
				this._html_template = widget_weather_template;
			}
		}

		this.widget_base_config = {
			id: this.id,
			title_font_size: this.title_font_size,
			derogation: !this.fullscreenMode,
			exportMode: this.exportMode
		};

		//title

		if(this.display_name) {
			this.widget_base_config.title = this.display_name;
		}
		else if(this.data.display_name) {
			this.widget_base_config.title = this.data.display_name;
		}
		else if(this.component) {
			this.widget_base_config.title = this.component;
		}
		else {
			this.widget_base_config.title = 'Unknown';
		}

		if(this.hide_title) {
			this.widget_base_config.title = '';
		}

		var linkUrl = this.formatLink();

		if(this.fullscreenMode && linkUrl) {
			this.widget_base_config.title = '<a href="' + linkUrl + '" target="_newtab">' + this.widget_base_config.title + '</a>';
		}

		//check ressource admin
		if(global.accountCtrl.check_right(this.data, 'w')) {
			this.widget_base_config.admin = true;
		}

		// build html
		if (this.data) {
			this.build(this.data);
		}
		else {
			this.buildEmpty();
		}

		// get element
		this.edit_button = this.getEl().getById(this.id + '-edit_button');

		// bindings
		var report_button = this.getEl().getById(this.id + '-button');

		if(report_button) {
			report_button.on('click', this.report_issue, this);
		}

		var clickable_title = this.getEl().getById(this.id + '-title');

		if(clickable_title && (this.external_link || this.link) && !this.fullscreenMode) {
			clickable_title.addCls('weather-clickable');
			clickable_title.on('click', this.externalLink, this);
		}

		if(this.widget_base_config.admin && this.display_derogation_icon && this.edit_button) {
			var output = this.getEl().getById(this.id + '-edit_td');

			if(output) {
				output.hover(
					function() {
						this.edit_button.fadeIn();
					},
					function() {
						this.edit_button.fadeOut();
					},
					this
				);
			}

			if(this.edit_button) {
				this.edit_button.on('click', function() {
					if(!this.data.rk) {
						global.notify.notify(_('Information not found'), _("Please wait a moment, some informations aren't availables"), 'info');
					}
					else {
						var name = this.component;

						if(this.resource) {
							name += ' - ' + this.resource;
						}

						global.derogationCtrl.derogate(this.data.rk, name, true);
					}
				}, this);
			}
		}

		//Hack for removing scrolling bar on ie
		this.getEl().parent().setStyle('overflow-x', 'hidden');


	},

	onResize: function() {
		//very dirty hack, ie resize images after that
		if(Ext.isIE) {
			this.getEl().down('.weather-image').setStyle('width', this.getWidth() * 0.20);
			this.getEl().down('.weather-td-image').setStyle('width', this.getWidth() * 0.20);
		}
	},

	build: function(data) {
		log.debug('  +  Build html for ' + data._id, this.logAuthor);

		var widget_data = {};

		if(data.state !== undefined) {
			widget_data.icon_src = this.getIcon(data.state);
		}
		else {
			widget_data.icon_src = this.info_weather_icon;
		}

		if(data.last_state_change) {
			widget_data.legend = rdr_elapsed_time(data.last_state_change, true);
		}

		if(data.timestamp) {
			widget_data.event_ts = rdr_tstodate(data.timestamp, true);
		}

		if(data.output && data.output !== '') {
			widget_data.output = data.output;
		}

		if(data.event_type === 'sla' && data.perf_data_array && data.perf_data_array[0].unit === '%') {
			widget_data.percent = data.perf_data_array[0].value;
		}

		// alert && derog
		if(this.display_report_button) {
			widget_data.button_text = _('Report issue');
		}

		if(this.data.alert_msg) {
			widget_data.alert_msg = this.data.alert_msg;
		}
		else {
			widget_data.alert_msg = '&nbsp;';
		}

		if(this.data.alert_icon !== undefined) {
			widget_data.alert_icon = this.alert_icon_basedir + this.alert_icon_name[this.data.alert_icon];
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

		if(this.link) {
			if(this.link.indexOf('http://') !== -1 || this.link.indexOf('www.') !== -1 || this.link.indexOf('https://') !== -1) {
				if(this.link.indexOf('http://') === -1 && this.link.indexOf('https://') === -1) {
					this.link = 'http://' + this.link;
				}

				window.open(this.link, '_newtab');
			}
			else {
				Ext.getStore('Views').load({
					scope: this,
					callback: function() {
						store = Ext.getStore('Views');
						var record = store.findExact('crecord_name', this.link);

						if (record !== -1) {
							record = store.getAt(record);

							if(!this.fullscreenMode) {
								global.tabsCtrl.open_view({
									view_id: record.get('_id'),
									title: _(record.get('crecord_name'))
								});
							}
							else {
								var url = Ext.String.format('http://{0}/static/canopsis/display_view.html?view_id={1}&authkey={2}',
									$(location).attr('host'),
									record.get('_id'),
									global.account.authkey
								);

								window.open(url, '_newtab');
							}
						}
						else {
							global.notify.notify('Link is not valid', 'The specified link does not match any view or URL', 'info');
						}
					}
				});
			}
		}
	},

	//fast hack for freeze, open link in tab, will be changed in develop
	formatLink: function() {
		if(typeof(this.link) === 'string' && !this.exportMode) {
			if(this.link.indexOf('http://') !== -1 || this.link.indexOf('www.') !== -1 || this.link.indexOf('https://') !== -1) {
				if(this.link.indexOf('http://') === -1 && this.link.indexOf('https://') === -1) {
					this.link = 'http://' + this.link;
				}

				return this.link;
			}
			else {
				store = Ext.getStore('Views');
				var record = store.findExact('crecord_name', this.link);

				if(record !== -1) {
					record = store.getAt(record);

					var url = Ext.String.format('http://{0}/static/canopsis/display_view.html?view_id={1}&authkey={2}',
						$(location).attr('host'),
						record.get('_id'),
						global.account.authkey
					);
					return url;
				}
				else {
					log.debug('Link is not valid', 'The specified link does not match any view or URL', 'info');
					return undefined;
				}
			}
		}
		else {
			return undefined;
		}
	},

	getIcon: function(value) {
		return 'widgets/weather/icons/set' + this.iconSet + '/state_' + value + '.png';
	}
});
