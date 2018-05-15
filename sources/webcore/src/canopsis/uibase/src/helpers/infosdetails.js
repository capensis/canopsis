/*
 * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
 *
 * This file is part of Canopsis.
 *
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Canopsis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
 */

Ember.Application.initializer({
    name: 'InfosDetailsHelper',
    after: 'DatesUtils',
    initialize: function(container, application) {
        void(application);

        var datesUtils = container.lookupFactory('utility:dates');

        var get = Ember.get,
            isNone = Ember.isNone;


        Ember.Handlebars.helper('infosdetails', function(info, value) {

            var details = '';
            for (var prop in info) {
				if (info.hasOwnProperty(prop)) {

					if (info[prop]["description"] === undefined || info[prop]["description"] === ""){
						if (info[prop]["value"] !== ""){
							details = details + '<ul><li>' + prop +  " : " + info[prop]["value"] + '</li></ul>';
						} else {
							details = details + '<ul><li>' + prop + '</li></ul>';
						}

					} else{
						details = details + '<ul><li>' + prop + '</li>';
						if (info[prop]["value"] !== ""){
							details = details + '<ul><li>' + info[prop]["description"] + ' : ' + info[prop]["value"] + '</li></ul>';
						} else {
							details = details + '<ul><li>' + info[prop]["description"] + '</li></ul>';
						}
						details = details + '</ul>';
					}
				}
            }
            return new Ember.Handlebars.SafeString(details);
        });
    }
});
