Ember.Application.initializer({
    name: 'component-rendererextradetails',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;
            __ = Ember.String.loc;

        var component = Ember.Component.extend({

            init: function() {
                this._super();
            },

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

            snoozetooltip: function() {
                if (this.get('hasSnooze')) {
                    return ['<center>',
                        '<b>' + __('Snooze') + '</b><br/>',
                        '<i>' + __('Since') + '</i> : <br/>',
                        this.dateFormat(this.get('value.snooze.t')) +' <br/> ',
                        __('By') +' : ' + this.get('value.snooze.a') +' <br/><br/> ',
                        '</center>'].join('');
                } else {
                    return '';
                }
            }.property('value.snooze.a', 'hasSnooze'),

            tickettooltip: function() {
                if (this.get('hasTicket')) {
                    return ['<center>',
                        '<b>' + __('Ticket declared') + '</b><br/>',
                        this.dateFormat(this.get('value.ticket.t')) +' <br/> ',
                        __('By') +' : ' + this.get('value.ticket.a') +' <br/><br/> ',
                        '</center>'].join('');
                } else {
                    return '';
                }
            }.property('value.ticket.a', 'hasTicket'),

            pbehaviortooltip: function() {
                if (this.get('hasPBehavior')) {
                    var self = this;
                    return ['<center>',
                        '<b>' + __('Periodic behavior') + '</b><br/>',
                        this.get('value.pbehaviors').map(function(pbeh) {
                            return pbeh.behavior + ' <br/>'
                                + self.dateFormat(pbeh.dtstart) + ' - ' + self.dateFormat(pbeh.dtstop) + ' <br/>'
                                + pbeh.rrule + ' <br/>'
                        }).join('') + ' <br/><br/> ',
                        '</center>'].join('');
                } else {
                    return '';
                }
            }.property('value.pbehaviors.@each.behavior', 'hasPBehavior'),

            hasSnooze: function() {
                return this.get('value.snooze') != null;
            }.property('value.snooze'),

            hasTicket: function() {
                return this.get('value.ticket') != null;
            }.property('value.ticket'),

            hasAck: function() {
                return this.get('value.ack') != null;
            }.property('value.ack'),

            hasPBehavior: function() {
                return this.get('value.pbehaviors') != null;
            }.property('value.pbihaviors'),

            dateFormat: function (date) {
                var mEpoch = parseInt(date);
                if (mEpoch < 10000000000) mEpoch *= 1000; // convert to milliseconds (Epoch is usually expressed in seconds, but JS uses milliseconds)
                var dDate = new Date(mEpoch);
                return moment(dDate).format('MM/DD/YY hh:mm:ss');
            },
        });

        application.register('component:component-rendererextradetails', component);
    }
});