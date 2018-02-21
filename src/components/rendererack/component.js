Ember.Application.initializer({
    name: 'component-rendererack',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;
        /**
         * This is the rendererack component for the widget listalarm
         *
         * @class rendererack component
         */

        var component = Ember.Component.extend({

            /**
             * @property propertyMap
             */
            propertyMap: {
                't': 'date',
                'a': 'user'
            },

            /**
             * @property properties
             */
            properties: ['t', 'a'],

            /**
             * @property propertiesMap
             */
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

            /**
             * @method init
             */
            init: function() {
                this._super();
              },
            
        });

        application.register('component:component-rendererack', component);
    }
});