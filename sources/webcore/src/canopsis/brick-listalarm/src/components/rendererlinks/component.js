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

                // The linkbuilders used to return the links directly as
                // strings. They can now also return objects with the
                // properties 'label' and 'link', allowing to change the link's
                // label.
                // The following code converts both representations into a
                // couple (label, link), so they can be displayed in the same
                // manner by the template.
                var tableau = [];
                for(var category in this.value){
                    this.value[category].forEach(function(item, index) {
                        if (typeof(item) == 'object' && item.hasOwnProperty('link') && item.hasOwnProperty('label')) {
                            tableau.push([category + " - " + item['label'], item['link']]);
                        } else {
                            tableau.push([category, item]);
                        }
                    });
                }
                set(this, 'categories', tableau);
            }

        });

        application.register('component:component-rendererlinks', component);
    }
});
