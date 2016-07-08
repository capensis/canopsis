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
    after: 'HashUtils',
    initialize: function(container, application) {

        var get = Ember.get,
            set = Ember.set,
            moment = window.moment;

        var component = Ember.Component.extend({
            timelineData: undefined,

            statusToName: {
                'ack':'Acknowledge by ',
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
                'cancel':{'icon':'fa-close','color':'bg-green'},
                'uncancel':{'icon':'fa-close','color':'bg-yellow'},
                'statusinc':{'icon':'fa-flag','color':undefined},
                'statusdec':{'icon':'fa-flag','color':undefined},
                'stateinc':{'icon':'fa-flag','color':undefined},
                'statedec':{'icon':'fa-flag','color':undefined},
                'changestate':{'icon':'fa-flag','color':undefined}
            },

            /**
             * @method didInsertElement
             * @description contains Rrule-editor initialisation and data binding
             */
            didInsertElement: function() {
                set(this,'steps',get(this,'timelineData').v.steps);

                for(var i = 0; i < get(this,'steps').length;i++){
                    var step = get(this,'steps')[i];

                    //build time related information
                    var date = new Date(get(this,'steps')[i].t*1000);
                    step.date = moment(date).format('LL');
                    step.time = moment(date).format('h:mm:ss a');
                    
                    //build color class
                    step.color = get(this,'iconsAndColors')[step._t].color;
                    //set icon
                    step.icon = get(this,'iconsAndColors')[step._t].icon;

                    //if no color, it's a state/value, so color
                    if(!step.color)
                        step.color = get(this,'colorArray')[step.val];

                    //if _t in ['stateinc','statedec','changestate']
                    if(step._t.indexOf('state') > -1)
                        step.state = get(this,'stateArray')[step.val];
                    
                    step.name = get(this,'statusToName')[step._t];
                }
            }
        });
        
        application.register('component:component-timeline', component);
    }
});