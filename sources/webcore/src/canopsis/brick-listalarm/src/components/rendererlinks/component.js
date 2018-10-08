Ember.Application.initializer({
    name: 'component-rendererlinks',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var component = Ember.Component.extend({

            /**
             * @method init
             */
            init: function() {
                this._super();
                var tableau = [];
                for(var category in this.value){
					this.value[category].forEach(function(item, index) {
						tableau.push([category, item]);
    				});
				}
                set(this, 'categories', tableau);
            }

        });

        application.register('component:component-rendererlinks', component);
    }
});
