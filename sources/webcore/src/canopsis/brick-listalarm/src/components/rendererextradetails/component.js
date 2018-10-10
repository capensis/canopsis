Ember.Application.initializer({
    name: 'component-rendererextradetails',
    after: ['DatesUtils'],
    initialize: function(container, application) {
        var datesUtils = container.lookupFactory('utility:dates');
        var __ = Ember.String.loc;

        /**
         * This is the rendererextradetails component for the widget listalarm
         *
         * @class rendererextradetails component
         */

        var component = Ember.Component.extend({

            /**
             * @method init
             */
            init: function() {
                this._super();

                // set extraDetailsComponent to be able to rerender it later since
                // property does not triggers
                this.set('value.extraDetailsComponent', this);
            },

            /**
             * @property acktooltip
             */
            acktooltip: function() {
                if (this.get('hasAck')) {
                    return ['<center>',
                        '<b>' + __('Ack') + '</b><br/>',
                        '<i>' + __('Date') + '</i> : <br/>',
                        this.dateFormat(this.get('value.ack.t')) +' <br/> ',
                        __('By') +' : ' + this.get('value.ack.a') +' <br/><br/> ',
                        '<i>'+__('Comment') +' </i> : <br/>' + this.get('value.ack.m'),
                        '</center>'].join('');
                } else {
                    return '';
                }
            }.property('value.ack.t', 'hasAck'),

            /**
             * @property snoozetooltip
             */
            snoozetooltip: function() {
                if (this.get('hasSnooze')) {
                    return ['<center>',
                        '<b>' + __('Snooze') + '</b><br/>',
                        '<i>' + __('Since') + '</i> : <br/>',
                        this.dateFormat(this.get('value.snooze.t')) +' <br/> ',
                        '<i>' + __('To') + '</i> : <br/>',
                        this.dateFormat(this.get('value.snooze.val')) +' <br/> ',
                        __('By') +' : ' + this.get('value.snooze.a') +' <br/><br/> ',
                        '</center>'].join('');
                } else {
                    return '';
                }
            }.property('value.snooze.a', 'hasSnooze'),

            /**
             * @property tickettooltip
             */
            tickettooltip: function() {
                if (this.get('hasTicket')) {
                    return ['<center>',
                        '<b>' + __('Ticket declared') + '</b><br/>',
                        this.dateFormat(this.get('value.ticket.t')) +' <br/> ',
                        'Ticket number: ' + this.get('value.ticket.val') +' <br/> ',
                        __('By') +' : ' + this.get('value.ticket.a') +' <br/><br/> ',
                        '</center>'].join('');
                } else {
                    return '';
                }
            }.property('value.ticket.a', 'hasTicket'),

            /**
             * @property donetooltip
             */
            donetooltip: function() {
                if (this.get('hasDone')) {
                    return ['<center>',
                        '<b>' + __('Alarm done') + '</b><br/>',
                        this.dateFormat(this.get('value.done.t')) +' <br/> ',
                        __('Message') +' : ' + this.get('value.done.m') +' <br/> ',
                        __('By') +' : ' + this.get('value.done.a') +' <br/><br/> ',
                        '</center>'].join('');
                } else {
                    return '';
                }
            }.property('value.done.a', 'hasDone'),

            /**
             * @property pbehaviortooltip
             */
            pbehaviortooltip: function() {
                if (this.get('hasPBehavior')) {
                    var self = this;
                    return ['<center>',
                        '<b>' + __('Periodic behavior') + '</b><br/>',
                        this.get('value.pbehaviors').map(function(pbeh) {
                            pbeh.rrule = pbeh.rrule || __('No reccurence');
                            return pbeh.name + ' <br/>'
                                + self.dateFormat(pbeh.tstart) + ' - ' + self.dateFormat(pbeh.tstop) + ' <br/>'
                                + pbeh.rrule + ' <br/>';
                        }).join('') + ' <br/><br/> ',
                        '</center>'].join('');
                } else {
                    return '';
                }
            }.property('value.pbehaviors.@each', 'hasPBehavior'),

            /**
             * @property hasSnooze
             */
            hasSnooze: function() {
                return this.get('value.snooze') != null;
            }.property('value.snooze'),

            /**
             * @property hasTicket
             */
            hasTicket: function() {
                return this.get('value.ticket') != null;
            }.property('value.ticket'),

            /**
             * @property hasDone
             */
            hasDone: function() {
                return this.get('value.done') != null;
            }.property('value.done'),

            /**
             * @property hasAck
             */
            hasAck: function() {
                return this.get('value.ack') != null;
            }.property('value.ack'),

            /**
             * @property hasPBehavior
             */
            hasPBehavior: function() {
                if (this.get('value.pbehaviors') == null) {
                    return false;
                }
                return this.get('value.pbehaviors').length != 0;
            }.property('value.@each.pbehaviors.length'),

            /**
             * @property hasActivePBehavior
             */
            hasActivePBehavior: function() {
                if (this.get('value.pbehaviors') == null) {
                    return false;
                }
                for (i = 0; i < this.get('value.pbehaviors').length; i++){
                    pb_stop = this.get('value.pbehaviors')[i].tstop * 1000
                    pb_start = this.get('value.pbehaviors')[i].tstart * 1000
                    if (Date.now() < pb_stop && Date.now() > pb_start) {
                        return true;
                    }
                }
                return false;
            }.property('value.@each.pbehaviors.length'),

            /**
             * @property dateFormat
             */
            dateFormat: function (date) {
                var mEpoch = parseInt(date);
                return datesUtils.timestamp2String(mEpoch,'f', true);
            },
        });

        application.register('component:component-rendererextradetails', component);
    }
});
