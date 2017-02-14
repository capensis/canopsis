Ember.Application.initializer({
    name: 'component-rendererack',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var component = Ember.Component.extend({

            propertyMap: {
                't': 'date',
                'a': 'user'
            },

            properties: ['t', 'a'],

            propertiesMap: function () {
                var propertyMap = this.get('propertyMap');
                var val = this.get('value');
                return this.get('properties').map(function(prop) {
                    return {
                        'key': propertyMap[prop] || prop,
                        'value': get(val, prop)
                    }
                })  
            }.property('properties', 'propertyMap', 'value'),

            init: function() {
                this._super();
              },
            
        });

        application.register('component:component-rendererack', component);
    }
});