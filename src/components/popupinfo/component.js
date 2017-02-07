Ember.Application.initializer({
    name: 'component-popupinfo',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var component = Ember.Component.extend({

            init: function() {
                this._super();
                this.set('columnTemplate', Ember.columnTemplates.findBy('columnName', this.get('columnName')).columnTemplate)
            },

            upd: function () {
                if (this.get('columnName') == this.get('clickedField.humanName')) {
                    $('.popupinfo').hide();                                                
                    console.error(this.get('columnName'), this.get('columnTemplate'), this.get('clickedField'), this.get('clickedAlarm'))
                    this.$('.popupinfo').fadeIn(500);                    
                }
            }.observes('updater'),

            actions: {
                hide: function () {
                    this.$('.popupinfo').fadeOut(500);
                }
            }

        });

        application.register('component:component-popupinfo', component);
    }
});