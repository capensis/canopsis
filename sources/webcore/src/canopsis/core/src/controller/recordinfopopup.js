/**
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
    name: 'RecordinfopopupController',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            __ = Ember.String.loc;


        var eventedController = Ember.Controller.extend(Ember.Evented);

        var controller = eventedController.extend({
            init: function () {
                set(this, 'title',__('Information'));
                console.log('initializing recordinfopopup controller');
            },

            actions: {
                show: function(crecord, template) {
                    console.log('Show recordinfopopup', crecord, template);

                    var html;

                    try {
                        html = Handlebars.compile(template)(crecord[0].toJson());
                    } catch (err) {
                        html = '<i>An error occured while compiling the template with the record. please if check the template is correct</i>';
                    }

                    set(this, 'content', new Ember.Handlebars.SafeString(html));

                    //FIXME do not use jquery for that kind of things on a controller
                    var left = ($(window).width() - $('#recordinfopopup').outerWidth()) / 2;
                    $('#recordinfopopup').css('left', left);
                    $('#recordinfopopup').fadeIn(500);
                },

                hide: function() {
                    console.log('hiding recordinfopopup');
                    $('#recordinfopopup').fadeOut(500);
                },
            }

        });
        application.register('controller:recordinfopopup', controller);
    }
});
