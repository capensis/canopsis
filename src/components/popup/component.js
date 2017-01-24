Ember.Application.initializer({
    name: 'component-popup',
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
                // set(this, 'tt', true);

                var cView = get(this, 'cView');
                var column = get(this, 'column');

                set(this, 'column', column);
                set(this, 'columnView', cView);
            },

            // columnView: function () {
            //     return this.get('cV');
            // }.property('cV'),

            popup: function () {
                // console.error('popup Modal', this.get('popupToggleModal'));
                // console.error('pop', this.get('clickedAlarm'), this.get('clickedField'));

                if (this.get('clickedField.name') == this.get('column')) {
                    this.set('columnViewContext', this.get('clickedAlarm'));
                    this.$('.modal').modal("show"); 
                }
                // this.set('columnViewContext', this.get('popupToggleModal'));
                // if (this.get('popupToggleModal.state') > 1) {
                //     this.set('columnView', this.get('arr')[0]);               
                // } else {
                //     this.set('columnView', this.get('arr')[0]);               
                // }
                
                // this.rerender();
                // this.$('.modal').modal("show");               
            }.observes('listener'),

            // actions: {
            //     open: function () {
            //         console.error('open');
            //         // this.set('tt', !this.get('tt'));
            //         this.$('.modal').modal("show");
            //     }
            // }

        });

        application.register('component:component-popup', component);
    }
});