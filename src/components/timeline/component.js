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

            statusToName: {
                "ack":"Acknowledge by ",
                "assocticket":"Ticket association by ",
                "declareticket":"Ticket declared by ",
                "cancel":"Canceled by ",
                "uncancel":"Uncanceled by ",
                "statusinc":"Status increased",
                "statusdec":"Status decreased",
                "stateinc":"State increased",
                "statedec":"State decreased",
                "changestate":"State changed"
            },

            stateArray: [
                "Ok",
                "Minor",
                "Major",
                "Critical"
            ],

            /**
             * @method didInsertElement
             * @description contains Rrule-editor initialisation and data binding
             */
            didInsertElement: function() {
                set(this,"steps",this.timelineData.v.steps)
                console.error(this.timelineData)
                for(var i = 0; i < this.steps.length;i++){
                    var date = new Date(this.steps[i].t*1000);
                    var step = this.steps[i]
                    step.date = moment(date).format('LL');
                    step.time = moment(date).format("h:mm:ss a");
                    step.color = "bg-"

                    switch(step._t){
                        case "ack":
                            step.icon = "fa-check"
                            step.color += "purple";
                            break;
                        case "assocticket":
                            step.icon = "fa-ticket"
                            step.color += "blue";
                            break;
                        case "declareticket":
                            step.icon = "fa-ticket"
                            step.color += "blue";
                            break;
                        case "cancel":
                            step.icon = "fa-close"
                            step.color += "green";
                            break;
                        case "uncancel":
                            step.icon = "fa-close"
                            step.color += "yellow";
                            break;
                        case "statusinc":
                        case "statusdec":
                            step.icon = "fa-flag"
                            break;
                        case "stateinc":
                        case "statedec":
                        case "changestate":
                            step.state = this.stateArray[step.val]
                            step.icon = "fa-flag"
                            break;
                        default:
                            break;
                    }


                    switch(step.val){
                        case 0:
                            step.color += "green";
                            break;
                        case 1:
                            step.color += "yellow";
                            break;
                        case 2:
                            step.color += "orange";
                            break;
                        case 3:
                            step.color += "red";
                            break;
                    }
                    
                    step.name = this.statusToName[step._t]
                }
            }
        })
        
        application.register('component:component-timeline', component);
    }
});