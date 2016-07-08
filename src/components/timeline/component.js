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

        var __ = Ember.String.loc,
            get = Ember.get,
            set = Ember.set,
            isArray = Ember.isArray

        var component = Ember.Component.extend({
            timelineData: undefined,

            /**
             * @method didInsertElement
             * @description contains Rrule-editor initialisation and data binding
             */
            didInsertElement: function() {
                set(this,"steps",this.timelineData.v.steps)
                console.error(this.timelineData)
                for(var i = 0; i < this.steps.length;i++){
                    var date = new Date(this.steps[i].t*1000);
                    this.steps[i].date = moment(date).format('LL');
                    this.steps[i].time = moment(date).format("h:mm:ss a");
                    this.steps[i].color = "bg-"

                    switch(this.steps[i]._t){
                        case "ack":
                            this.steps[i].icon = "fa-check"
                            this.steps[i].color += "purple";
                            break;
                        case "assocticket":
                            this.steps[i].icon = "fa-ticket"
                            this.steps[i].color += "blue";
                            break;
                        case "declareticket":
                            this.steps[i].icon = "fa-ticket"
                            this.steps[i].color += "blue";
                            break;
                        case "cancel":
                            this.steps[i].icon = "fa-close"
                            this.steps[i].color += "green";
                            break;
                        case "uncancel":
                            this.steps[i].icon = "fa-close"
                            this.steps[i].color += "yellow";
                            break;
                        case "statusinc":
                        case "statusdec":
                            this.steps[i].icon = "fa-envelope"
                            break;
                        case "stateinc":
                        case "statedec":
                        case "changestate":
                            this.steps[i].icon = "fa-envelope"
                            break;
                        default:
                            break;
                    }


                    switch(this.steps[i].val){
                        case 0:
                            this.steps[i].color += "green";
                            break;
                        case 1:
                            this.steps[i].color += "yellow";
                            break;
                        case 2:
                            this.steps[i].color += "orange";
                            break;
                        case 3:
                            this.steps[i].color += "red";
                            break;
                    }
                    console.error(this.steps[i]._t)
                }
            }
        })
        
        application.register('component:component-timeline', component);
    }
});