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
    name:"component-periodicbehaviormanager",
    after: 'CrudMixin',
    initialize: function(container, application) {
        var CrudMixin = container.lookupFactory('mixin:crud');

        var get = Ember.get,
            set = Ember.set;
            __ = Ember.String.loc;


        var CrudEventedComponent = Ember.Component.extend(Ember.Evented, CrudMixin);

        var component = CrudEventedComponent.extend({
            init: function() {
                /* mixin options for mixins */
                set(this, 'mixinOptions', {
                    /* specific to crud mixin */
                    crud: {
                        form: 'modelform'
                    }
                });

                this._super.apply(this, arguments);

                /* store for pbehaviors fetching */
                var store = DS.Store.create({
                    container: get(this, 'container')
                });

                set(this, 'widgetDataStore', store);

                this.refreshContent();
                this.on('refresh', this.refreshContent);
            },

            refreshContent: function() {
                console.group('Fetching periodic behaviors');
                console.log('context:', get(this, 'contextId'));

                var store = get(this, 'widgetDataStore'),
                    ctrl = this;

                store.findQuery(
                    'pbehavior',
                    {
                        entity_ids: get(this, 'contextId')
                    }
                ).then(
                    function(result) {
                        set(ctrl, 'behaviors', get(result, 'content'));
                        console.log('behaviors:', get(ctrl, 'behaviors'));
                    }
                );

                console.groupEnd();
            },

            onRecordReady: function(record) {
                this._super.apply(this, arguments);

                var contextId = get(this, 'contextId');
                set(record, 'source', contextId);
            },

            tableColumns: [
                {name: 'dtstart', title: __('From')},
                {name: 'dtend', title: __('To')},
                {name: 'rrule', title: __('Recursion')},
                {name: 'duration', title: __('Duration')},
                {name: 'behaviors', title: __('Behaviors')},
                {
                    title: new Ember.Handlebars.SafeString('<span class="glyphicon glyphicon-plus-sign"></span>'),
                    action: 'edit',
                    actionAll: 'addBehavior',
                    style: 'text-align: center;'
                },
                {
                    title: new Ember.Handlebars.SafeString('<span class="glyphicon glyphicon-trash"></span>'),
                    action: 'remove',
                    style: 'text-align: center;'
                }
            ],

            actions: {
                addBehavior: function() {
                    this.send('add', 'pbehavior');
                }
            }
        });

        application.register('component:component-periodicbehaviormanager', component);
    }
});
