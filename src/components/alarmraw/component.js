Ember.Application.initializer({
    name: 'component-alarmraw',
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
            tagName: 'tr',

            init: function() {
                this._super();

                set(this, 'alarm', get(this, 'alarm'));
                set(this, 'fields', get(this, 'fields'));                
            },

            expandClass: function () {
                var cl = '';
                if (this.get('alarm.isExpanded')) {
                    cl = 'glyphicon-minus';
                } else {
                    cl = 'glyphicon-plus';
                }
                return 'glyphicon ' + cl;
            }.property('alarm.isExpanded'),

            actions: {
                tdClick: function (alarm, field) {
                    this.sendAction('action', alarm, field);
                },

                sendAction: function (action, alarm) {
                    this.sendAction('saction', action, alarm);
                },

                expand: function () {
                    this.toggleProperty('alarm.isExpanded');
                }
            }

        });

        application.register('component:component-alarmraw', component);
    }
});