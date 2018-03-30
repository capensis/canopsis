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
    name: 'component-stringclassifiedcrecordselector',
    after: 'component-classifiedcrecordselector',
    initialize: function(container, application) {

        var Classifiedcrecordselector = container.lookupFactory('component:component-classifiedcrecordselector');

        var get = Ember.get,
            set = Ember.set;

        /**
         * @component stringclassifiedcrecordselector
         */
        var component = Classifiedcrecordselector.extend({
            multiselect:false,

            selectionChanged: function(){
                var selectionUnprepared = get(this, 'selectionUnprepared')[0];
                var res;

                var valueKey = get(this, 'valueKey');

                if(!Ember.isNone(selectionUnprepared)) {
                    if(valueKey) {
                        res = get(selectionUnprepared, 'value');
                    } else {
                        res = get(selectionUnprepared, 'name');
                    }
                    console.log('selection changed', res);
                }
                set(this, 'selection', res);
            }.observes('selectionUnprepared', 'selectionUnprepared.@each')
        });

        application.register('component:component-stringclassifiedcrecordselector', component);
    }
});
