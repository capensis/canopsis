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

if(typeof(i18n) === 'undefined') {
	i18n = {};
}

function _(text, context) {
	var ttext = undefined;

	// Check if locales is loaded
	if(typeof(i18n) === 'undefined') {
		return text;
	}

	if(context) {
		ttext = i18n[context + '.' + text];

		if(ttext) {
			return ttext;
		}

		ttext = i18n[context + '.' + Ext.String.capitalize(text)];

		if(ttext) {
			return Ext.String.uncapitalize(ttext);
		}

		ttext = i18n[context + '.' + Ext.String.uncapitalize(text)];

		if(ttext) {
			return Ext.String.capitalize(ttext);
		}
	}

	ttext = i18n[text];

	if(ttext) {
		return ttext;
	}

	ttext = i18n[Ext.String.capitalize(text)];

	if(ttext) {
		return Ext.String.uncapitalize(ttext);
	}

	ttext = i18n[Ext.String.uncapitalize(text)];

	if (ttext) {
		return Ext.String.capitalize(ttext);
	}

	// Translate failed
	if(global && global.log.level > 4 && !Ext.Array.contains(global.untranslated, text)) {
		global.untranslated.push(text);
	}

	if(global && global.log.level > 4) {
		return "->> " + text + " <<-";
	}
	else {
		return text;
	}
}
