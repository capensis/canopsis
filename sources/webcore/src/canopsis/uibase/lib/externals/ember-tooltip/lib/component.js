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

define([], function() {

    var get = Ember.get,
        set = Ember.set,
        __ = Ember.String.loc;


    var component = Ember.Component.extend({
        tagName: 'span',
        template: Ember.HTMLBars.compile('{{yield}}'),

        htmlEnabled: true,
        placement: 'top',
        triggerEvent: 'hover',
        noTranslation: false,

        didInsertElement: function() {
            this.popoverTargets = this.$("[data-toggle=popover]");
            if(this.popoverTargets) {
                this.popoverTargets.popover();
            }

            var component = this;

            var options = {
                html : get(component, 'htmlEnabled'),
                placement: get(component, 'placement'),
                trigger: get(component, 'triggerEvent'),
                container:'body'
            };

            var content = get(component, 'content');

            if (!Ember.isNone(content)) {
                if(get(component, 'noTranslation')) {
                    options.content = content;
                } else {
                    options.content = __(content);
                }
            }

            var title = get(component, 'title');

            if (!Ember.isNone(title)) {
                if(get(component, 'noTranslation')) {
                    options.title = title;
                } else {
                    options.title = __(title);
                }
            }

            this.$().popover(options);
        },

        willDestroyElement: function() {
            if(this.popoverTargets) {
                this.popoverTargets.popover('destroy');
            }

            this.$().popover('destroy');
        }
    });

    Ember.Application.initializer({
        name:"component-tooltip",
        initialize: function(container, application) {
            application.register('component:component-tooltip', component);
        }
    });

    return component;
});
