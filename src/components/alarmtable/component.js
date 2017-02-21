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
                this.set('isAllChecked', false);
            },

            // isAllChecked: function() {return truexxxx}.property(),

            allSelectionObserver: function() {
                var val = this.get('isAllChecked');
                this.get('alarms').setEach('isSelected', val);
            }.observes('isAllChecked'),

            currentSortColumn: function() {
                return get(this, 'defaultSortColumn');
            }.property('defaultSortColumn'),

            // sAlarms: function() {
            //     return this.get('alarms');
            // }.property('alarms'),

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
                },

                tdClick: function (alarm, field) {
                    // if cell is clickable
                    if (Ember.columnTemplates.findBy('columnName', field.humanName)) {       
                        this.set('clickedAlarm', alarm);
                        this.set('clickedField', field);
                        this.set('updater', (new Date()).getTime());
                    }
                },

                sendAction: function (action, alarm) {
                    this.sendAction('saction', action, alarm);
                }
            }
            
        });

        application.register('component:component-alarmtable', component);
    }
});