/*
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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
    'datetimepicker'
], function() {

    var get = Ember.get,
        set = Ember.set;

    var component = Ember.Component.extend({

        template: Ember.HTMLBars.compile('{{input class="form-control"}}'),

        init: function () {
            this._super.apply(this, arguments);
        },

        didInsertElement: function (){
            //@doc http://eonasdan.github.io/bootstrap-datetimepicker/
            var timepicker = this.$();

            var timepickerComponent = this;

            console.log('timestamp init', get(this, 'content'));
            if (get(this, 'content') === 2147483647 || get(this, 'content') === undefined) {
                timepicker.datetimepicker({
                    useSeconds: true, //en/disables the seconds picker
                    useCurrent: true, //when true, picker will set the value to the current date/time
                    language: 'fr'
                });
            } else {
                timepicker.datetimepicker({
                    useSeconds: true, //en/disables the seconds picker
                    useCurrent: false, //when true, picker will set the value to the current date/time
                    defaultDate: new Date(get(this, 'content')*1000),
                    language: 'fr'
                });
            }

            console.log('timestamp dom init complete');

            timepicker.on("dp.change",function (e) {
                var timestamp = new Date(e.date).getTime() / 1000;
                set(timepickerComponent, 'content', timestamp);
                console.log('timestamp date set', timestamp);
            });
        }
    });

    Ember.Application.initializer({
        name:"component-datetimepicker",
        initialize: function(container, application) {
            application.register('component:component-datetimepicker', component);
        }
    });

    return component;
});