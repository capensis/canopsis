Ember.Application.initializer({
    name: 'component-columntemplate',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This is the columntemplate component for the widget listalarm
         *
         * @class columntemplate component
         */
        var component = Ember.Component.extend({

            /**
             * @property template
             */
            template: undefined,

            /**
             * @property column
             */
            column: undefined,

            /**
             * @method init
             */
            init: function() {
                this._super();
                if (!isNone(this.get('content'))) {
                    set(this, 'column', get(this, 'content.column'));
                    set(this, 'template', get(this, 'content.template'));
                }
              },

            /**
             * @method onUpdate
             */
            onUpdate: function() {
                this.set('content', {
                    column: get(this, 'column'),
                    template: get(this, 'template')
                });
            }.observes('template', 'column')

        });

        application.register('component:component-columntemplate', component);
    }
});