Ember.Application.initializer({
    name: 'component-search',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var component = Ember.Component.extend({
            classNames: ['pull-left', 'input-group', 'input-group-sm', 'col-xs-3', 'search', 'search-padding'],

            init: function() {
                this._super();
                console.error(this.get('isValid'))
              },

            removeInvalidSearchTextNotification: function () {
                this.set('isValid', true)
            }.observes('value'),

            actions: {
                search: function () {
                    if (this.get('value').length > 0) {
                        this.sendAction('action', this.get('value'));
                    }
                },

                resetValue: function () {
                    this.set('value', '');
                    this.sendAction('action', '');                   
                }
            }

        });

        application.register('component:component-search', component);
    }
});