Ember.Application.initializer({
    name: 'component-alarmstate',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This is the alarmstate component for the widget listalarm
         *
         * @class alarmstate component
         */
        var component = Ember.Component.extend({
            /**
             * @property state
             */
            state: undefined,

            /**
             * @property isSelected
             */
            isSelected: 0,

            /**
             * @method init
             */
            init: function() {
                this._super();
                if (!isNone(this.get('content'))) {
                    set(this, 'state', get(this, 'content.state'));
                }
              },

            /**
             * @method onUpdate
             */
            onUpdate: function() {
                this.set('content', {
                    state: get(this, 'state'),
                });
            }.observes('state')

        });

        application.register('component:component-alarmstate', component);
    }
});