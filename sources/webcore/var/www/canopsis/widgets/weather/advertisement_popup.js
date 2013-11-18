//need:app/lib/view/cpopup.js,app/lib/form/field/cdate.js
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

Ext.define('widgets.weather.advertisement_popup' , {
	extend: 'canopsis.lib.view.cpopup',
	alias: 'widget.weather.advertisement_popup',

	requires: [
		'canopsis.lib.form.field.cdate'
	],

	_component: undefined,
	referer: undefined,
	title: _('New advertisement'),
	width: 300,

	icon_wip: 'widgets/weather/icons/public_domain_icon/workman.png',
	icon_warning: 'widgets/weather/icons/public_domain_icon/slippery.png',
	icon_alert: 'widgets/weather/icons/public_domain_icon/alert.png',
	icon_class: 'widget-weather-form-icon',

	_buildForm: function() {
		this._form.add({
			xtype: 'fieldset',
			defaultType: 'radio',
			title: _('Manual state change'),
			layout: 'anchor',
			items: [
				{
					boxLabel: '<img src="' + this.icon_wip + '" class="' + this.icon_class + '"/>' + _('Planned shut down'),
					name: 'state',
					checked: true,
					anchor: '100%',
					inputValue: 'icon-wip'
				}, {
					boxLabel: '<img src="' + this.icon_warning + '" class="' + this.icon_class + '"/>' + _('Be cautious with application'),
					anchor: '100%',
					name: 'state',
					inputValue: 'icon-warning'
				}, {
					boxLabel: '<img src="' + this.icon_alert + '" class="' + this.icon_class + '"/>' + _('Standart alert'),
					anchor: '100%',
					name: 'state',
					inputValue: 'icon-alert'
				}
			]
		});

		this._form.add({
			xtype: 'displayfield',
			value: _('Alert comment') + ':'
		});

		this._form.add({
			xtype: 'textfield',
			name: 'alert_comment',
			anchor: '100%',
			emptyText: _('Type here the alert comment')
		});

		this._form.add({
			xtype: 'displayfield',
			value: _('Visible standart comment') + ':'
		});

		this._form.add({
			xtype: 'textfield',
			name: 'standart_comment',
			anchor: '100%',
			emptyText: _('Type here the visible comment')
		});

		this._form.add({
			xtype: 'fieldset',
			title: _('Period'),
			layout: {
				type: 'vbox',
				align: 'center'
			},
			items: [{
					xtype: 'cdate',
					name: 'startTs',
					label_text: _('From')
				},{
					xtype: 'cdate',
					name: 'stopTs',
					label_text: _('To')
				}]
		});

		return this._form;
	},

	_ok_button_function: function() {
		log.dump(this._form.getValues());
	}
});
