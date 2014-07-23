/*
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
*/


define([
	'jquery',
	'app/lib/loaders/factories',
], function($, factories) {

	var timestampEditor = factories.Editor("timestamp", {

		init: function() {
			this._super();
			console.log("EditorTimestampEditorView init");
			this.set('pickTime', true);
			this.set('pickDate', true);
		},

		dateOnly: function () {
			this.set('pickTime', false);
		}.property('pickTime'),

		timeOnly: function () {
			this.set('pickDate', false);
		}.property('pickDate'),

		test: function () {
			console.log('test');
		},

		didInsertElement: function() {
			//@doc http://eonasdan.github.io/bootstrap-datetimepicker/
			$('#' + this.get('elementId')).datetimepicker({
				pickTime: this.get('pickTime'),
				pickDate: this.get('pickDate'),
				language: 'fr'
			});
		}
	});

	return timestampEditor;
});