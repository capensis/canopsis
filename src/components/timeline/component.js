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
                'ack':'Acknowledged by ',
                'assocticket':'Ticket association by ',
                'declareticket':'Ticket declared by ',
                'cancel':'Canceled by ',
                'uncancel':'Uncanceled by ',
                'statusinc':'Status increased',
                'statusdec':'Status decreased',
                'stateinc':'State increased',
                'statedec':'State decreased',
                'changestate':'State changed'
            },

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

            colorArray:[
                'bg-green',
                'bg-yellow',
                'bg-orange',
                'bg-red'
            ],

            iconsAndColors: {
                'ack':{'icon':'fa-check','color':'bg-purple'},
                'assocticket':{'icon':'fa-ticket','color':'bg-blue'},
                'declareticket':{'icon':'fa-ticket','color':'bg-blue'},
                'cancel':{'icon':'glyphicon glyphicon-ban-circle','color':'bg-gray'},
                'uncancel':{'icon':'glyphicon glyphicon-ban-circle','color':'bg-gray'},
                'statusinc':{'icon':'fa-sort-amount-asc','color':'bg-gray'},
                'statusdec':{'icon':'fa-sort-amount-desc','color':'bg-gray'},
                'stateinc':{'icon':'fa-flag','color':undefined},
                'statedec':{'icon':'fa-flag','color':undefined},
                'changestate':{'icon':'fa-flag','color':undefined}
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

                    /* result can be null if entity has no current alarm. */
                    if (result === null) {
                        set(component,'steps',[]);
                    } else {
                        set(component,'steps',result.data[0].value.steps);

                        for(var i = 0; i < get(component,'steps').length;i++){
                            var step = get(component,'steps')[i];

                            //build time related information
                            var date = new Date(step.t*1000);
                            step.date = moment(date).format('LL');
                            step.time = moment(date).format('h:mm:ss a');

                            //build color class
                            step.color = get(component,'iconsAndColors')[step._t].color;
                            //set icon
                            step.icon = get(component,'iconsAndColors')[step._t].icon;

                            //if no color, it's a state/value, so color
                            if(!step.color)
                                step.color = get(component,'colorArray')[step.val];

                            //if _t in ['stateinc','statedec','changestate']
                            if(step._t.indexOf('state') > -1)
                                step.state = get(component,'stateArray')[step.val];

                            if(step._t.indexOf('status') > -1)
                                step.status = get(component,'statusArray')[step.val];

                            step.name = get(component,'statusToName')[step._t];
                        }
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
