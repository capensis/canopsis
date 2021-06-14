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
Ext.define('canopsis.lib.menu.cspinner' , {
	extend: 'Ext.button.Button',

	alias: 'widget.cspinner',

	logAuthor: '[cspinner]',

	width: 28,

	spinner: undefined,
	spinner_options: {
		color:'#fff',
		lines: 10,
		width: 2,
		length: 3,
		shadow: true,
		radius: 3,
		top: 1,
		left: 1
	},

	tEl: undefined,

	ajax_queue: 0,

	handleMouseEvents: false,

	initComponent: function() {
		log.debug('Initializing...', this.logAuthor);
	},

	start: function(){
		if(!this.tEl) {
			this.tEl = document.getElementById(this.id);
		}

		if(!this.spinner) {
			log.debug('Start spinner', this.logAuthor);
			this.spinner = new Spinner(this.spinner_options).spin(this.tEl);
		}
	},

	stop: function() {
		if(this.spinner && this.ajax_queue <= 0) {
			log.debug('Stop spinner', this.logAuthor);

			this.spinner.stop();
			delete this.spinner;
			this.spinner = undefined;
			this.ajax_queue = 0;
		}
	},

	bind_Ext_Ajax: function() {
		log.debug('Bind spinner on Ajax requests', this.logAuthor);

		Ext.Ajax.on('beforerequest', function() {
			this.ajax_queue += 1;
			this.start();
		}, this);

		Ext.Ajax.on('requestexception', function() {
			this.ajax_queue -= 1;
			this.stop();
		}, this);

		Ext.Ajax.on('requestcomplete', function() {
			this.ajax_queue -= 1;
			this.stop();
		}, this);
	},

	afterRender: function() {
		this.callParent(arguments);
		this.bind_Ext_Ajax();
	}
});
