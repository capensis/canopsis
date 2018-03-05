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
    name: 'component-arrayclassifiedcrecordselector',
    after: 'component-classifiedcrecordselector',
    initialize: function(container, application) {

        var Classifiedcrecordselector = container.lookupFactory('component:component-classifiedcrecordselector');
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @component arrayclassifiedcrecordselector
         */
        var component = Classifiedcrecordselector.extend({

            /**
             * @property multiselect
             */
            multiselect: true,

            /**
             * @method setInitialContent
             * @argument initialContent
             */
            setInitialContent: function(initialContent) {
                var valueKey = get(this, 'valueKey');

                console.log('setInitialContent', valueKey, typeof initialContent, initialContent);
                if(initialContent) {
                    if(valueKey) {
                        set(this, 'loadingInitialContent', 'true');
                    } else {
                        var selectionUnprepared = get(this, 'selectionUnprepared');
                        for (var i = 0, l = initialContent.length; i < l; i++) {
                            selectionUnprepared.pushObject({'name': initialContent[i]});
                        }
                    }
                }
            },

            /**
             * @method selectionChanged
             */
            selectionChanged: function(){
                this._super();
                //additional code ensuring single item selection and use of possible custom valueKey.
                var selection = get(this, 'selection');

                var valueKey = get(this, 'valueKey');
                if (isNone(valueKey)) {
                    valueKey = 'name';
                }

                //simple cache object to avoid duplicates values
                var cache = {};
                //no duplication selection list computation
                var new_selection = [];
                //simple content values computation
                var content = [];

                //iteraing over previous selection in order to recompute it.
                for (var i=0; i<selection.length; i++) {
                    var value = get(selection[i], valueKey);
                    if (!cache[value]) {
                        cache[value] = 1;
                        content.push(value);
                        new_selection.push(selection[i]);
                    }


                }
                set(this, 'selection', new_selection);
                set(this, 'content', content);

            }.observes('selectionUnprepared', 'selectionUnprepared.@each'),

            /**
             * @method extractItems
             * @arguments items
             */
            extractItems: function(items) {
                var valueKey = get(this, 'valueKey');
                var initialContent = get(this, 'content');

                if(valueKey) {
                    var resBuffer = [];
                    console.log('extractItems', arguments);
                    for (var i = 0, l = initialContent.length; i < l; i++) {
                        var correspondingExtractedItem = items.findBy('id', initialContent);

                        if(correspondingExtractedItem !== undefined) {
                            resBuffer.pushObject({ name: correspondingExtractedItem.get(valueKey)});
                        }
                    }
                    this.setProperties({
                        selectionUnprepared: resBuffer,
                        loadingInitialContent: false
                    });
                }
            }
        });

        application.register('component:component-arrayclassifiedcrecordselector', component);
    }
});
