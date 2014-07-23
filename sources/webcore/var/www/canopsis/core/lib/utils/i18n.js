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
	'canopsis/canopsisConfiguration',
	'app/application'
], function(conf) {

	var i18n = {
		lang: 'fr',
		translations: {
			todo: {},
			'fr': {}
		},
		newTranslations: true,
		_: function(word) {
			if (i18n.translations[i18n.lang] && i18n.translations[i18n.lang][word]) {
				return i18n.translations[i18n.lang][word];
			} else {
				//adding translation to todo list
				if (typeof(word) === 'string' && !i18n.translations.todo[word]) {
					i18n.translations.todo[word.toLowerCase().trim()] = 1;
					i18n.newTranslations = true;
				}
				//returns original not translated string
				return word;
			}
		},
		todo: {},
		uploadDefinitions: function () {

			$.ajax({
				url: '/rest/misc/i18n',
				type: 'POST',
				data: JSON.stringify({
					id: 'translations',
					translations: i18n.translations,
					crecord_type: 'i18n'
				}),
				success: function(data) {
					if (data.success) {
						console.log('Upload lang upload complete');
					}
				},
				async: false
			});
		},
		downloadDefinitions: function () {
			$.ajax({
				url: '/rest/misc/i18n/translations',
				success: function(data) {
					if (data.success) {
						i18n.translations = data.data[0].translations;
					}
				},
				async: false
			}).fail(function () {
				console.log('initialization case. translation is now ready');
				i18n.uploadDefinitions();
			});
		},

	};

	i18n.fields = {
		ack_this_alert: i18n._('Acknowlege this alert'),
		cancel_this_alert: i18n._('Cancel this alert'),
		uncancel_this_alert: i18n._('Undo alert cancellation')
	};

	i18n.downloadDefinitions();

	window._ = i18n._;
	window.tr = i18n.fields;

	if (conf.DEBUG && conf.TRANSLATE) {
		setInterval(function () {
			if (i18n.newTranslations) {
				console.log('Uploading new translations');
				i18n.newTranslations = false;
				i18n.uploadDefinitions();
				i18n.downloadDefinitions();
			}

		}, 10000);
	}

	return i18n;
});