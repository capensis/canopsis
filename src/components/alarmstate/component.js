Ember.Application.initializer({
    name: 'component-alarmstate',
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

            state: undefined,
            isSelected: 0,

            init: function() {
                this._super();
                if (!isNone(this.get('content'))) {
                    set(this, 'state', get(this, 'content.state'));
                }
              },

            onUpdate: function() {
                this.set('content', {
                    state: get(this, 'state'),
                });
            }.observes('state')

        });

        application.register('component:component-alarmstate', component);
    }
});