Ember.Application.initializer({
    name: 'component-rendererstatusval',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var component = Ember.Component.extend({

            list: {
                0: 'undef',
                1: 'on going',
                2: 'undef',
                3: 'undef',
                4: 'cancelled'
            },

            init: function() {
                this._super();
            },

            status: function() {
                return this.get('list')[this.get('value')];
            }.property('value'),
        });

        application.register('component:component-rendererstatusval', component);
    }
});