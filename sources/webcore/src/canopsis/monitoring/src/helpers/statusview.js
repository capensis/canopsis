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
    name: 'StatusViewHelper',
    after: 'DatesUtils',
    initialize: function(container, application) {
        void(application);

        var datesUtils = container.lookupFactory('utility:dates');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;


        Ember.Handlebars.helper('statusview', function(status, crecord) {
            /**
                # status legend:
                # 0 == Ok
                # 1 == On going
                # 2 == Stealthy
                # 3 == Bagot
                # 4 == Canceled
            **/

            var statuses = {
                0: 'Off',
                1: 'On going',
                2: 'Stealthy',
                3: 'Bagot',
                4: 'Cancelled'
            };

            if (isNone(status)) {
                status = get(crecord, 'status');
            }

            var value = statuses[status] || '';
            set(crecord, 'statusvalue', __(value));

            if(status === 4) {

                //displays cancel information if any onto the status field
                var cancel = get(crecord, 'record.cancel');
                console.log('statusview', status, cancel);

                if(!isNone(cancel)) {

                    var timestamp = get(cancel, 'timestamp');
                    var comment = get(cancel, 'comment');
                    var author = get(cancel, 'author');

                    set(crecord, 'statushtml', [
                        '<center>',
                        '<i>' , __('Date') , '</i> : <br/>',
                        datesUtils.timestamp2String(timestamp) ,' <br/> ',
                        __('By'), ' : ' , author ,' <br/><br/> ',
                        '<i>', __('Comment') ,'</i> : <br/>' , comment,
                        '</center>'
                    ].join(''));
                }

            }
            return '';
        });
    }
});
