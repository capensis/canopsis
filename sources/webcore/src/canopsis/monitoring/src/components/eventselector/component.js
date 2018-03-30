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
    name: 'component-eventselector',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;

        /**
         * @component eventSelector
         * @description Displays an interface to search and select events
         *
         * ![Component preview](../screenshots/component-eventselector.png)
         */
        var component = Ember.Component.extend({
            /**
             * @property component
             * @type string
             * @description The value of the component input that is used to search for events
             */
            component: undefined,

            /**
             * @property resource
             * @type string
             * @description The value of the resource input that is used to search for events
             */
            resource: undefined,

            /**
             * @property componentDataStore
             * @type DS.Store
             * @description the store where searched events are disposed
             */
            componentDataStore: undefined,

            /**
             * @property content
             */
            content: undefined,

            /**
             * @property events
             */
            events: undefined,

            /**
             * @property labelled
             */
            labelled: undefined,

            /**
             * @property saveLabelsDone
             * @type boolean
             * @description Assigned to true when the labels have been saved
             */
            saveLabelsDone: undefined,

            /**
             * @property search_component
             */
            search_component: undefined,

            /**
             * @property search_resource
             */
            search_resource: undefined,

            /**
             * @property selectedEvents
             */
            selectedEvents: undefined,

            /**
             * @property selectors
             * @type boolean
             * @description When true, the search is only done on selector-based events
             */
            selectors: undefined,

            /**
             * @property topologies
             * @type boolean
             * @description When true, the search is only done on topology-based events
             */
            topologies: undefined,

            /**
             * @property type_label
             */
            type_label: undefined,

            /**
             * @method init
             */
            init: function() {
                this._super();
                this.set('componentDataStore', DS.Store.create({
                    container: this.get('container')
                }));
                console.log('Event selector init');

                this.set('selectedEvents', []);

                set(this, 'search_resource', __('Search for a resource'));
                set(this, 'search_component', __('Search for a component'));
                set(this, 'type_label', __('Associate a label to this event'));

                if (get(this, 'content') !== undefined) {
                    this.initializeEvents();
                }
            },

            /**
             * @method initializeEvents
             */
            initializeEvents: function () {

                var rks = [];

                var label_information = {};

                //takes care of how the events are loaded when a label is given to them.
                if (get(this, 'labelled')) {
                    var events_information = get(this, 'content');
                    var length = events_information.length;
                    for (var i=0; i<length; i++) {
                        rks.push(events_information[i].rk);
                        label_information[events_information[i].rk] = events_information[i].label;
                    }
                } else {
                    rks = get(this, 'content');
                }

                var eventselectorController = this;
                var query = get(this, 'componentDataStore').findQuery(
                    'event',
                    {
                        filter: JSON.stringify({_id: {'$in': rks}}),
                        limit: 0
                    }
                ).then(
                    function (data) {
                        console.log('Fetched initialization data from events', data.content);

                        //When labels are associated to events, they are set in them on load for this editor.
                        if (get(eventselectorController, 'labelled')) {

                            var length = data.content.length;
                            for (var i=0; i<length; i++) {
                                var rk = get(data.content[i], 'id');
                                set(data.content[i], 'label', label_information[rk]);
                            }

                        }
                        set(eventselectorController, 'selectedEvents', data.content);
                    }
                );
                void (query);
            },

            /**
             * @method findEvents
             */
            findEvents: function() {

                var filter = {};

                var excludeRks = this.getSelectedRks();

                //adding exclusion rks if any loaded
                if (excludeRks.length) {
                    filter._id = {'$nin': excludeRks};
                }

                //permissive search throught component and resource
                if (this.get('component')) {
                    filter.component = { '$regex' : '.*'+ this.get('component') +'.*', '$options': 'i' };
                }
                if (this.get('resource')) {
                    filter.resource = { '$regex' : '.*'+ this.get('resource') +'.*', '$options': 'i' };
                }

                filter.event_type = 'check';

                //does user selected selector or topology search
                if (this.get('selectors')) {
                    filter.event_type = {'$in' :['selector', 'sla']};
                }

                if (this.get('topologies')) {
                    filter.source_type = {'$eq' :'topo'};
                }

                if (!filter.resource && !filter.component) {
                    this.set('events', []);
                    //when user only wants topologies or selectors, query is done anyway with the right crecord type
                    if (!this.get('topologies') && !this.get('selectors')) {
                        return;
                    }
                }

                var query = get(this, 'componentDataStore').findQuery(
                    'event',
                    {
                        filter: JSON.stringify(filter),
                        limit: 10
                    }
                );

                var that = this;
                query.then(
                    function (data) {
                        console.log('Fetched data from events', data.content);
                        that.set('events', data.content);
                    });

                void (query);

            }.observes('component', 'resource'),

            /**
             * @method setSelector
             * @description Observer on "selectors". This observer is triggered when the user clicks on the "Selectors" checkbox. Refreshes the list of found events to filter only selectors
             */
            setSelector: function() {
                this.set('topologies', false);
                this.findEvents();
            }.observes('selectors'),

            /**
             * @method setTopologies
             * @description Observer on "topologies". This observer is triggered when the user clicks on the "Topologies" checkbox. Refreshes the list of found events to filter only topologies
             */
            setTopologies: function() {
                this.set('selectors', false);
                this.findEvents();
            }.observes('topologies'),

            /**
             * @method getSelectedRks
             * @returns {array} A list of selected event routing keys
             */
            getSelectedRks: function() {

                var selectedRks = [];
                var isLabelled = get(this, 'labelled');
                var selectedEvents = get(this, 'selectedEvents');

                if (!isNone(selectedEvents)) {

                    var length = selectedEvents.length;

                    for (var i=0; i<length; i++) {

                        if (isLabelled) {
                            var label = get(selectedEvents[i], 'label');
                            if (isNone(label)) {
                                //auto name event
                                label = 'default_' + i;
                            } else {
                                //avoid space in names
                                label = label.trim();
                            }
                            selectedRks.push({
                                rk: get(selectedEvents[i], 'id'),
                                label: label
                            });

                        } else {
                            selectedRks.push(get(selectedEvents[i], 'id'));
                        }
                    }
                }
                return selectedRks;
            },

            actions: {
                /**
                 * @method actions_saveLabels
                 * @description persists the label values into the computed result
                 */
                saveLabels: function () {
                    set(this, 'content', this.getSelectedRks());
                    set(this, 'saveLabelsDone', true);
                    var eventselectorController = this;
                    setTimeout(function () {
                        set(eventselectorController, 'saveLabelsDone', false);
                    }, 3000);
                },

                /**
                 * @method actions_add
                 * @param event
                 * @description Action triggered when the user adds an event from the selection
                 */
                add: function (event) {
                    console.log('Adding event', event);
                    get(this, 'selectedEvents').pushObject(event);
                    get(this, 'events').removeObject(event);
                    set(this, 'content', this.getSelectedRks());

                    if (get(this,'events').length === 0) {
                        this.findEvents();
                    }
                },

                /**
                 * @method actions_delete
                 * @param event
                 * @description Action triggered when the user removes an event from the selection
                 */
                delete: function (event) {
                    console.log('Rk to delete', event.id);
                    get(this, 'selectedEvents').removeObject(event);
                    this.findEvents();
                    set(this, 'content', this.getSelectedRks());
                }
            }
        });

        application.register('component:component-eventselector', component);
    }
});
