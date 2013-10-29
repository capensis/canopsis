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
Ext.define('canopsis.controller.Notify', {
	extend: 'Ext.app.Controller',

	opacity: 0.9,
	history: false,
	logAuthor: '[controller][cnotify]',

	init: function() {
		global.notify = this;

		log.debug('[controller][cnotify] - Initialize ...');

		this.callParent(arguments);

		$.pnotify.defaults.history = false;
	},

	test: function() {
		this.notify(_('Title'), _('Description'));
		this.notify(_('Title'), _('Description'), 'info');
		this.notify(_('Title'), _('Description'), 'success');
		this.notify(_('Title'), _('Description'), 'warning');
		this.notify(_('Title'), _('Description'), 'error');
	},

	notify: function(title, text, type, icon, hide, closer, sticker) {
		if(type === undefined) {
			type = 'info';
		}

		if(hide === undefined) {
			hide = true;
		}

		if(closer === undefined) {
			closer = true;
		}

		if(sticker === undefined) {
			sticker = false;
		}

		$.pnotify({
			title: title,
			text: text,
			delay: 3500,
			type: type,
			history: this.history,
			notice_icon: icon,
			opacity: this.opacity,
			hide: hide,
			closer: closer,
			sticker: sticker
		});

		log.debug('Display notification', this.logAuthor);
	}
});
