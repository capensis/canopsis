Ember.Application.initializer({
    name: 'component-rendererstatusval',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var component = Ember.Component.extend({

            list: {
                0: 'Off',
                1: 'On going',
                2: 'Stealthy',
                3: 'Bagot',
                4: 'Cancel'
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