Ember.Application.initializer({
    name: 'component-alarmtable',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This is the alarmtable component for the widget listalarm
         *
         * @class alarmtable component
         */
        var component = Ember.Component.extend({

            /**
             * @method isChangedByUser
             */
            init: function() {
                this._super();
                this.set('isAllChecked', false);
            },

            /**
             * @method allSelectionObserver
             */
            allSelectionObserver: function() {
                var val = this.get('isAllChecked');
                this.get('alarms').setEach('isSelected', val);
            }.observes('isAllChecked'),

            /**
             * @property currentSortColumn
             */
            currentSortColumn: function() {
                return get(this, 'defaultSortColumn');
            }.property('defaultSortColumn'),

            /**
             * @property columnsAmount
             */
            columnsAmount: function () {
                return this.get('fields.length') + 3;
            }.property('fields.length'),

            actions: {

                /**
                 * @method click
                 */
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

                /**
                 * @method tdClick
                 */
                tdClick: function (alarm, field) {
                    // if cell is clickable
                    if (Ember.columnTemplates.findBy('columnName', field.humanName)) {
                        this.set('clickedAlarm', alarm);
                        this.set('clickedField', field);
                        this.set('updater', (new Date()).getTime());
                    }
                },

                /**
                 * @method sendAction
                 */
                sendAction: function (action, alarm) {
                    this.sendAction('saction', action, alarm);
                }
            }

        });

        application.register('component:component-alarmtable', component);
    }
});