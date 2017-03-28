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
    name: 'component-timeline',
    after: ['DataUtils', 'HashUtils'],
    initialize: function(container, application) {

        var dataUtils = container.lookupFactory('utility:data');

        var get = Ember.get,
            set = Ember.set,
            moment = window.moment;

        var component = Ember.Component.extend({
                        timelineData: undefined,

            statusToName: {
                'ack': 'Acknowledged by ',
                'ackremove': 'Ack removed by ',
                'assocticket': 'Ticket association by ',
                'declareticket': 'Ticket declared by ',
                'cancel': 'Canceled by ',
                'comment': 'Comment by ',
                'uncancel': 'Restored by ',
                'statusinc': 'Status increased',
                'statusdec': 'Status decreased',
                'stateinc': 'State increased',
                'statedec': 'State decreased',
                'changestate': 'State changed',
                'snooze': 'Snoozed by ',
                'statecounter': 'Cropped states (since last change of status)',
                'hardlimit': 'Hard limit reached !'
            },

            authoredName: [
                'ack',
                'ackremove',
                'assocticket',
                'declareticket',
                'cancel',
                'comment',
                'uncancel',
                'snooze'
            ],

            stateArray: [
                'Ok',
                'Minor',
                'Major',
                'Critical'
            ],

            statusArray: [
                'off',
                'ongoing',
                'stealthy',
                'bagot',
                'canceled'
            ],

            colorArray: [
                'bg-green',
                'bg-yellow',
                'bg-orange',
                'bg-red'
            ],

            iconsAndColors: {
                'ack': {'icon': 'fa-check', 'color': 'bg-purple'},
                'ackremove': {'icon': 'glyphicon glyphicon-ban-circle', 'color': 'bg-purple'},
                'assocticket': {'icon': 'fa-ticket', 'color': 'bg-blue'},
                'declareticket': {'icon': 'fa-ticket', 'color': 'bg-blue'},
                'cancel': {'icon': 'glyphicon glyphicon-trash', 'color': 'bg-gray'},
                'comment': {'icon': 'fa-comment-o', 'color': 'bg-teal'},
                'uncancel': {'icon': 'glyphicon glyphicon-share', 'color': 'bg-gray'},
                'statusinc': {'icon': 'fa-chevron-up', 'color': 'bg-gray'},
                'statusdec': {'icon': 'fa-chevron-down', 'color': 'bg-gray'},
                'stateinc': {'icon': 'fa-flag', 'color': undefined},
                'statedec': {'icon': 'fa-flag', 'color': undefined},
                'changestate': {'icon': 'fa-flag', 'color': undefined},
                'snooze': {'icon': 'fa-clock-o', 'color': 'bg-fuchsia'},
                'statecounter': {'icon': 'fa-scissors', 'color': 'bg-black'},
                'hardlimit': {'icon': 'fa-warning', 'color': 'bg-red'}
            },

            /**
             * @method didInsertElement
             * @description contains Rrule-editor initialisation and data binding
             */
            didInsertElement: function() {
                var component = this;

                var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:alarm');
                var query = {'entity_id': get(component, 'timelineData').entity_id};

                adapter.findQuery('alarm', 'get-current-alarm', query).then(function (result) {
                    // onfullfillment

                    var steps = [];
                    for (var i = result.data[0].value.steps.length - 1 ; i >= 0 ; i--) {
                        var step = result.data[0].value.steps[i];

                        //build time related information
                        var date = new Date(step.t*1000);
                        step.date = moment(date).format('LL');
                        step.time = moment(date).format('h:mm:ss a');

                        if (!(step._t in get(component, 'iconsAndColors'))) {
                            console.warning('Unknown step "' + step._t + '" : skipping');
                            continue;
                        }

                        //build color class
                        step.color = get(component, 'iconsAndColors')[step._t].color;
                        //set icon
                        step.icon = get(component, 'iconsAndColors')[step._t].icon;

                        //if no color, it's a state/value, so color
                        if (!step.color)
                            step.color = get(component,'colorArray')[step.val];

                        if (step._t.indexOf('state') > -1)
                            step.state = get(component,'stateArray')[step.val];

                        if (step._t.indexOf('status') > -1)
                            step.status = get(component,'statusArray')[step.val];

                        if (step._t === 'snooze') {
                            var until = new Date(step.val * 1000);
                            step.until = moment(until).format('h:mm:ss a');
                        }

                        if (step._t === 'statecounter') {
                            step.m = '<table class="table table-hover"><tbody>';

                            step.m += '<tr><th>State increases</th><th>' + step.val.stateinc + '</th></tr>';
                            step.m += '<tr><th>State decreases</th><th>' + step.val.statedec + '</th></tr>';

                            for (v in step.val) {
                                if (v.startsWith('state:')) {
                                    var state = parseInt(v.replace('state:', ''), 10);
                                    var state_label = get(component, 'stateArray')[state];
                                    /* Custom states (other than 0, 1, 2, 3) are not supported */
                                    step.m += '<tr><th>State ' + state_label + '</th><th>' + step.val[v] + '</th></tr>';
                                } 
                            }

                            step.m += '</tbody></table>';
                        }

                        step.name = get(component, 'statusToName')[step._t];
                        if (get(component, 'authoredName').indexOf(step._t) != -1) {
                            step.name += step.a;
                        }

                        steps.push(step);
                    }

                    /* steps can be null if entity has no current alarm. */
                    set(component, 'steps', steps);
                }, function (reason) {
                    // onrejection
                    console.error('ERROR in the adapter: ', reason);
                });
            }
        });
        
        application.register('component:component-timeline', component);
    }
});
