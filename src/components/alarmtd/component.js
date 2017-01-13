Ember.Application.initializer({
    name: 'component-alarmtd',
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
            init: function() {
                this._super();

                set(this, 'alarm', get(this, 'alarm'));
                set(this, 'field', get(this, 'field'));                
              },

            value: function() {
                var alarm = get(this, 'alarm');
                var field = get(this, 'field');
                return get(alarm, field)
            }.property('alarm', 'field'),
        });

        application.register('component:component-alarmtd', component);
    }
});