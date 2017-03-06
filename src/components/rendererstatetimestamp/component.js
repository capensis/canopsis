Ember.Application.initializer({
    name: 'component-rendererstatetimestamp',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var component = Ember.Component.extend({

            init: function() {
                this._super();
            },

            timestamp: function() {
                return new Date(2017, 0, 30).getTime() / 1000; 
            }.property('value'),

        });

        application.register('component:component-rendererstatetimestamp', component);
    }
});