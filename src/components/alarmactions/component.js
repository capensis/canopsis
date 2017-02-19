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

            init: function() {
                this._super();
              },

            hasLinks: function() {
                return this.get('alarm.linklist.event_links.length') > 0;
            }.property('alarm.linklist.event_links'),
        });

        application.register('component:component-alarmactions', component);
    }
});