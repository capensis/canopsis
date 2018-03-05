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
    name: 'component-mixinselector',
    after: 'MixinsRegistry',
    initialize: function(container, application) {
        var mixinsRegistry = container.lookupFactory('registry:mixins');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @component mixinselector
         * @description A mixin selector. Fills a classifieditemselector with data from the mixins registry
         *
         * ![Component preview](../screenshots/component-mixinselector.png)
         */
        var component = Ember.Component.extend({

            /**
             * @method init
             */
            init: function() {
                this._super.apply(this, arguments);

                if(isNone(get(this, 'content'))) {
                    set(this, 'content', Ember.A());
                }

                set(this, 'selectionPrepared', Ember.A());

                var content = get(this, 'content');

                if(content) {
                    for (var i = 0, l = content.length; i < l; i++) {
                        var mixinSchema = window.schemasRegistry.getByName(content[i].name);

                        if(mixinSchema && mixinSchema.modelDict && mixinSchema.modelDict.metadata && mixinSchema.modelDict.metadata.description) {
                            content[i].description = mixinSchema.modelDict.metadata.description;
                        }

                        if(typeof content[i] === 'string') {
                            content[i] = { name: content[i] };
                        }
                    }
                }
                set(this, 'selectionPrepared', content);
            },

            /**
             * @property classifiedItems
             * @description The mixins registry
             * @type MixinsRegistry
             * @default mixinsRegistry
             */
            classifiedItems: mixinsRegistry,

            /**
             * @property selectionPrepared
             * @description Contains the selection managed by the classifieditemselector
             */
            selectionPrepared: undefined,

            /**
             * @property content
             * @description Contains the user selection, extracted from the classifieditemselector that can be used outside of the component
             */
            content: undefined,

            /**
             * @property selectionUnprepared
             * @default Ember.computed.alias('content')
             */
            selectionUnprepared: Ember.computed.alias('content'),

            /**
             * @method recomputeSelection
             * @description recalculates the selection and update the "content" property
             */
            recomputeSelection: function() {
                var selection = get(this, 'selectionPrepared');
                console.log('recomputeSelection', selection, get(this, 'selectionPrepared'));

                var content = get(this, 'content');

                var resBuffer = Ember.A();
                if(selection) {
                    for (var i = 0, l = selection.length; i < l; i++) {
                        var currentItem = selection[i];
                        var currentItemName = get(currentItem, 'name');
                        var newResBufferItem;

                        var existingContentItem = content.findBy('name', currentItemName);
                        if(existingContentItem) {
                            newResBufferItem = existingContentItem;
                        } else {
                            newResBufferItem = {
                                name: currentItemName
                            };
                        }
                        resBuffer.pushObject(newResBufferItem);
                    }
                }

                set(this, 'content', resBuffer);
            },

            actions: {
                /**
                 * @method actions_selectItem
                 * @description Calls the recomputeSelection method when the user selects an item
                 */
                selectItem: function() {
                    this.recomputeSelection();
                },

                /**
                 * @method actions_selectItem
                 * @description Calls the recomputeSelection method when the user unselects an item
                 */
                unselectItem: function(){
                    this.recomputeSelection();
                }
            }
        });
        application.register('component:component-mixinselector', component);
    }
});
