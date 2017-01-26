Ember.Application.initializer({
    name: 'component-columntemplate',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This is the eventcategories component for the widget calendar
         *
         * @class eventcategories component
         * @memberOf canopsis.frontend.brick-calendar
         */
        var component = Ember.Component.extend({

            template: undefined,
            column: undefined,

            init: function() {
                this._super();
                if (!isNone(this.get('content'))) {
                    set(this, 'column', get(this, 'content.column'));
                    set(this, 'template', get(this, 'content.template'));
                }
              },

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