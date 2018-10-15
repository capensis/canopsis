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
            __ = Ember.String.loc;

        var DEFAULT_AUTHOR = __('unknown');

        var component = Ember.Component.extend({
                        timelineData: undefined,

            statusToName: {
                'ack': __('Acknowledged by '),
                'ackremove': __('Ack removed by '),
                'assocticket': __('Ticket association by '),
                'declareticket': __('Ticket declared by '),
                'cancel': __('Canceled by '),
                'comment': __('Comment by '),
                'uncancel': __('Restored by '),
                'statusinc': __('Status increased'),
                'statusdec': __('Status decreased'),
                'stateinc': __('State increased'),
                'statedec': __('State decreased'),
                'changestate': __('State changed'),
                'snooze': __('Snoozed by '),
                'statecounter': __('Cropped states (since last change of status)'),
                'done': __('Mark as done by '),
                'hardlimit': __('Hard limit reached !'),
                'long_output': __("update_output")
            },

            authoredName: [
                'ack',
                'ackremove',
                'assocticket',
                'declareticket',
                'cancel',
                'comment',
                'uncancel',
                'snooze',
                'done'
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
                'done': {'icon': 'fa-check-square', 'color': 'bg-olive'},
                'hardlimit': {'icon': 'fa-warning', 'color': 'bg-red'},
                'long_output': {'icon': 'fa-pencil', 'color': 'bg-grey'}
            },

            addAuthor: ['stateinc', 'statedec', 'changestate', 'statusinc', 'statusdec', 'long_output'],

            /**
             * @method didInsertElement
             * @description contains Rrule-editor initialisation and data binding
             */
            didInsertElement: function() {
                var component = this;
                var entityData = get(component, 'timelineData');
                var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:timeline');
                var encoded_query = JSON.stringify({'d': entityData.entity_id,
                                                    'v.creation_date': entityData["v.creation_date"]})
                var query =  {'filter': encoded_query, 'opened':true, 'resolved':true,'sort_key':'t','sort_dir':'DESC','limit':1,'with_steps':true};


                adapter.findQuery('alarm', 'get-current-alarm', undefined, query).then(function (result) {
                    // onfullfillment
                    var previousDate = undefined;
                    var steps = [];
                    if (result.data.length > 0 && result.data[0].alarms.length > 0){

                        for (var i = result.data[0].alarms[0].v.steps.length - 1 ; i >= 0 ; i--) {
                            var step = result.data[0].alarms[0].v.steps[i];
                            var index_state_check = 0
                            var index_status_check = 0
                            var statusToName = ''

                            //build time related information
                            var date = new Date(step.t*1000);
                            step.date = moment(date).format('LL');
                            if(step.date != previousDate)
                                step.showDate = true;
                            else
                                step.showDate = false;

                            step.time = moment(date).format('HH:mm:ss');

							step.a = __(step.a)

                            if (!(step._t in get(component, 'iconsAndColors'))) {
                                console.warn('Unknown step "' + step._t + '" : skipping');
                                continue;
                            }

                            //build color class
                            step.color = get(component, 'iconsAndColors')[step._t].color;
                            //set icon
                            step.icon = get(component, 'iconsAndColors')[step._t].icon;

                            //if no color, it's a state/value, so color
                            if (!step.color)
                                step.color = get(component,'colorArray')[step.val];

                            if (step._t.indexOf('state') >= index_state_check)
                                step.state = get(component,'stateArray')[step.val];

                            if (step._t.indexOf('status') >= index_status_check)
                                step.status = get(component,'statusArray')[step.val];

                            if (step._t === 'snooze') {
                                var message = new Date(step.val * 1000);
                                step.m = __("snooze_until") + moment(message).format('HH:mm:ss');
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

                            statusToName = get(component, 'statusToName')[step._t];
                            if (component.get('addAuthor').includes(step._t)) {
                                if (step.a.indexOf('.') > -1){
                                    step.name = statusToName + ' ' + __('by ') + DEFAULT_AUTHOR;
                                } else {
                                    step.name = statusToName + ' ' + __('by ') + step.a;
                                }
                            } else {
                                step.name = statusToName;
                                if (get(component, 'authoredName').indexOf(step._t) != -1) {
                                    step.name += step.a;
                                }
                            }

                            // add value to a ticket's message
                            if (step._t === 'assocticket') {
								step.m = __("assoc_ticket") + step.val
                            }

                            steps.push(step);

                            //stock previous date
                            previousDate = step.date;

                        }

                        /* steps can be null if entity has no current alarm. */
                        set(component, 'steps', steps);
                    }
                }, function (reason) {
                    // onrejection
                    console.error('ERROR in the adapter: ', reason);
                });
            }
        });

        application.register('component:component-timeline', component);
    }
});
