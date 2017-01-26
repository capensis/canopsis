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

              actions: {
                  tdClick: function (alarm, field) {
                      this.sendAction('action', alarm, field);
                  }
              }

        });

        application.register('component:component-alarmraw', component);
    }
});