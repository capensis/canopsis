Ember.Application.initializer({
    name: 'component-rendererack',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var component = Ember.Component.extend({

            init: function() {
                this._super();
              },

            
        });

        application.register('component:component-rendererack', component);
    }
});