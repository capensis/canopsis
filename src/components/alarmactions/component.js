Ember.Application.initializer({
    name: 'component-alarmactions',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This is the eventcategories component for the widget calendar
         *
         * @class eventcategories component
         * @memberOf canopsis.frontend.brick-calendar
         */
        var component = Ember.Component.extend({
            tagName: 'td',
            classNames: ['action-cell'],

            actionsMap: Ember.A([
                {
                    class: 'glyphicon glyphicon-saved',
                    internal_states: ['unacked'],
                    name: 'ack',
                    mixin_name: 'ack'
                },
                {
                    class: 'glyphicon glyphicon-ok',
                    internal_states: ['unacked'],
                    name: 'fastack',
                    mixin_name: 'fastack'
                },
                {
                    class: 'glyphicon glyphicon-ban-circle',
                    internal_states: ['acked'],
                    name: 'unack',
                    mixin_name: 'ackremove'
                },
                {
                    class: 'fa fa-ticket',
                    internal_states: ['acked'],
                    name: 'declareticket',
                    mixin_name: 'declareticket'
                },
                {
                    class: 'fa fa-thumb-tack',
                    internal_states: ['acked', 'removed'],
                    name: 'assocticket',
                    mixin_name: 'assocticket'
                },
                {
                    class: 'glyphicon glyphicon-trash',
                    internal_states: ['acked'],
                    name: 'cancelalarm',
                    mixin_name: 'cancel'
                },
                {
                    class: 'fa fa-exclamation-triangle',
                    internal_states: ['acked'],
                    name: 'changecriticity',
                    mixin_name: 'changestate'
                },
                {
                    class: 'glyphicon glyphicon-share-alt',
                    internal_states: ['removed'],
                    name: 'restorealarm',
                    mixin_name: 'recovery'
                },
                {
                    class: 'fa fa-clock-o',
                    internal_states: ['unacked', 'acked','removed'],
                    name: 'snoozealarm',
                    mixin_name: 'snooze'
                }
                // TODO declare incident
            ]),

            init: function() {
                this._super();
              },

            availableActions: function() {
                var intState = this.get('internalState');
                return this.get('actionsMap').filter(function(item, index, enumerable) {
                    return item.internal_states.includes(intState)
                });
            }.property('internalState'),

            internalState: function() {
                if (this.get('state') == 0 && !this.get('isAcked')) {
                    return 'immutable';
                }
                if (this.get('state') > 0 && !this.get('isAcked')) {
                    return 'unacked';
                }
                if (this.get('isAcked')) {
                    return 'acked';
                }
                // if (!this.get('isCancelled')) {
                //     return 'cancelled';
                // }
                return 'removed';
                

                // if (this.get('state') == 0 && this.get('status') == 3) {
                //     return 'acked';
                // } else {
                //     if (this.get('state') == 1) {
                //         return 'cancelled'
                //     } else {
                //         return 'unacked'
                //     }
                // }
            }.property('alarm.state', 'alarm.status'),

            isAcked: function() {
                return this.get('alarm.ack._t') != undefined;
            }.property('alarm.ack._t'),

            // isCancelled: function () {
            //     return this.get('alarm.cancelled') != undefined;
            // }.property('alarm.cancelled'),

            hasLinks: function() {
                return this.get('alarm.linklist.event_links.length') > 0;
            }.property('alarm.linklist.event_links'),

            actions: {
                sendAction: function (action) {
                    this.sendAction('action', action, this.get('alarm'));
                }
            }
        });

        application.register('component:component-alarmactions', component);
    }
});