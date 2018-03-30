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


/*
'sendPbehavior': function(){

                    var actionContext = get(this,'actionContext');
                    var me = this;

                    // For now, type_ field is set by default
                    var event = {
                        "author":window.username,
                        "comments":[{"author":window.username,"message":me.get("pbehaviorComment")}],
                        "filter":{"_id": actionContext['entity_id']},
                        "tstart": parseInt(Date.now()/1000,10),
                        "tstop": 2147483647,
                        "name": "downtime",
                        "reason": me.get("pbehaviorReason"),
                        "type_": "pause"
                    }

                    // send request to create the pbehavior
                    action = function() {
                        return $.ajax({
                            type: 'POST',
                            url: '/api/v2/pbehavior',
                            data: JSON.stringify(event),
                            contentType: 'application/json',
                            dataType: 'json',
                            success: function(){
                                console.log("Ticket sent");
                                me.togglePbehavior()
                            },
                            error: function(){console.error("Failure to send downtime")}
                        });
                    }

                    // Notify the action as done
                    actionContext.get('posponedActions').pushObject(Ember.Object.create(
                        {
                            'actionName': 'pbehavior',
                            'action': action
                        })
                    )
                    this.get('availableActions').findBy('name', 'pause').set('isSaved', true)
                    this.togglePbehavior()

                }*/

Ember.Application.initializer({
    name: 'component-crudentitydetails',
    after: ['DataUtils', 'HashUtils'],
    initialize: function(container, application) {
        var dataUtils = container.lookupFactory('utility:data');
        var formsUtils = container.lookupFactory('utility:forms');

        var get = Ember.get,
            set = Ember.set,
            moment = window.moment;
            __ = Ember.String.loc;
        
        var component = Ember.Component.extend({

            init: function () {
                this._super();
                set(this, 'componentDataStore', DS.Store.create({
                    container: get(this, 'container')
                }));
            },

            didInsertElement: function() {
                var component = this;
                set(component, 'entity', this.this);
                this.findPbehaviors();
            },

            findPbehaviors: function () {

                var component = this,
                    adapter = dataUtils.getEmberApplicationSingleton()
                    .__container__.lookup('adapter:ccpbehavior');
                var entity = get(component, 'entity')
                
                adapter.findQuery('ccpbehavior', encodeURI(entity._id))
                .then(function (queryResults) {
                    set(component, 'pbehaviors', queryResults)
                });
            },

            removepbehavior: function (id) {
                var component = this
                var recordWizard = formsUtils.showNew('confirmform', {}, {
                    title: __('Do you want to remove this pbehavior?')
                });

                recordWizard.submit.then(function (form) {
                    // Send request to delete it
                    return $.ajax({
                        type: 'DELETE',
                        url: '/api/v2/pbehavior/' + id,
                        contentType: 'application/json',
                        dataType: 'json',
                        success: function () {
                            console.log("pbehavior removed");
                        },
                        error: function () { console.error("Failure to send remove pbehavior") }
                    });
                });
            },

            actions: {
                'removepbehavior': function (id) {
                    this.removepbehavior(id)
                }
            }

        });

        application.register('component:component-crudentitydetails', component);
    }
});
