/*
 * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
 *
 * This file is part of Canopsis.
 *
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Canopsis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
 */

Ember.Application.initializer({
    name: 'component-ack',
    after: 'DatesUtils',
    initialize: function(container, application) {

        var datesUtils = container.lookupFactory('utility:dates');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;

        /**
         * @component ack
         * @description Displays informations on a given ack
         *
         * ![Component preview](../screenshots/component-ack.png)
         *
         * @example {{component-ack crecord=this value=value}}
         */
        var component = Ember.Component.extend({
            /**
             * @property crecord
             * @type object
             */
            crecord: undefined,

            /**
             * @property value
             * @type object
             */
            value: undefined,

            /**
             * @property tickettooltip
             * @type string
             */
            tickettooltip: undefined,

            /**
             * @property acktooltip
             * @type string
             */
            acktooltip: undefined,

            /**
             * @property ackcolor
             */
            ackcolor: undefined,

            /**
             * @property acktitle
             */
            acktitle: undefined,

            /**
             * @method init
             */
            init: function() {
                this._super();

                var crecord = get(this, 'crecord');
                var value = get(this, 'value');

                console.log('Ack crecord', crecord, 'Ack value', value);

                //displays ticket information if any onto the status field
                var ticket_declared_author = get(crecord, 'record.ticket_declared_author');
                var ticket_declared_date = get(crecord, 'record.ticket_declared_date');

                var ticketNumber = get(crecord, 'record.ticket');
                var ticketDate = get(crecord, 'record.ticket_date');

                var ticketNumberHtml = '';
                //Generate ticket declared html information
                if (!isNone(ticketNumber) && !isNone(ticketDate)) {
                    ticketNumberHtml = [
                        '<b>' + __('Ticket number') + '</b><br/>',
                        datesUtils.timestamp2String(ticketDate) +' <br/> ',
                        '<i>' + ticketNumber +'</i><br/> '
                    ].join('');
                }

                var tickettooltip = '';
                //Generate html display for ticket declared and ticket number when possible
                if(!isNone(ticket_declared_date) && !isNone(ticket_declared_author)) {
                    tickettooltip = ['<center>',
                        '<b>' + __('Ticket declared') + '</b><br/>',
                        datesUtils.timestamp2String(ticket_declared_date) +' <br/> ',
                        __('By') +' : ' + ticket_declared_author +' <br/><br/> ',
                        ticketNumberHtml,
                        '</center>'
                    ].join('');
                    set(this, 'tickettooltip', tickettooltip);

                } else if (!isNone(ticketNumber) && !isNone(ticketDate)) {

                    //When no ticket declared, then ticket date was saved.
                    console.debug('ticket date is ', get(crecord, 'record.ticket_date'));

                    tickettooltip = ['<center>',
                        ticketNumberHtml,
                        '</center>'
                    ].join('');
                    set(this, 'tickettooltip', tickettooltip);
                }

                var ts = value && value.timestamp ? datesUtils.timestamp2String(value.timestamp): '';
                var author = value && value.author ? value.author: '';
                var comment =  value && value.comment ? value.comment: '';

                var acktooltip = ['<center>',
                    '<b>' + __('Ack') + '</b><br/>',
                    '<i>' + __('Date') + '</i> : <br/>',
                    ts +' <br/> ',
                    __('By') +' : ' + author +' <br/><br/> ',
                    '<i>'+__('Comment') +' </i> : <br/>' + comment,
                    '</center>'].join('');

                console.log('ack value', value);

                set(this, 'acktooltip', acktooltip);

                if(value && value.isCancel) {
                    set(this, 'ackcolor', '');
                    set(this, 'acktitle', __('Cancelled'));
                } else {
                    set(this, 'ackcolor', 'bg-purple');
                    set(this, 'acktitle', __('Acknowleged'));
                }
            }
        });

        application.register('component:component-ack', component);
    }
});
