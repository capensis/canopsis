Ember.Application.initializer({
    name: 'component-alarmtable',
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

            init: function() {
                this._super();
            },

            currentSortColumn: function() {
                return get(this, 'defaultSortColumn');
            }.property('defaultSortColumn'),

            sAlarms: function() {
                return this.get('alarms');
            }.property('alarms'),

            actions: {
                click: function (field) {
                    if (field == this.get('currentSortColumn')) {
                        this.set('currentSortColumn.isASC', !this.get('currentSortColumn.isASC'));
                    } else {
                        this.set('currentSortColumn.isSortable', false);
                        this.set('currentSortColumn', field);
                        this.set('currentSortColumn.isSortable', true);
                        this.set('currentSortColumn.isASC', true);
                    }
                    this.sendAction('action', this.get('currentSortColumn'));
                }
            }
            
        });

        application.register('component:component-alarmtable', component);
    }
});