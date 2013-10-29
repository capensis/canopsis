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
var log = {
	/* 0: none, 1: info, 2: error, 3: error + warning, 4: error + warning + debug, 5: error + warning + debug + dump */
	level: global.log.level,
	buffer: global.log.buffer,

	console: true,

	info: function(msg, author) {
		if(this.level >= 1) {
			this.writeMsg(msg, 1, author);
		}
	},

	debug: function(msg, author) {
		if(this.level >= 4) {
			this.writeMsg(msg, 4, author);
		}
	},

	warning: function(msg, author) {
		if(this.level >= 3) {
			this.writeMsg(msg, 2, author);
		}
	},

	error: function(msg, author) {
		if(this.level >= 2) {
			this.writeMsg(msg, 3, author);
		}
	},

	dump: function(msg, author) {
		if(this.level >= 5) {
			this.writeMsg(msg, 5, author);
		}
	},

	writeMsg: function(msg, level, author) {
		var level_msg = '';

		if(level === 1) {
			level_msg = 'INFO';
		}
		else if(level === 2) {
			level_msg = 'WARNING';
		}
		else if(level === 3) {
			level_msg = 'ERROR';
		}
		else if(level === 4) {
			level_msg = 'DEBUG';
		}
		else if(level === 5) {
			level_msg = 'DUMP';
		}

		if(author) {
			msg = author + ' - ' + msg;
		}

		// Chech if firebug console ...
		if(this.console) {
			try {
				console.log(msg);
			}
			catch (err) {
				this.console = false;
			}
		}
	}
};
