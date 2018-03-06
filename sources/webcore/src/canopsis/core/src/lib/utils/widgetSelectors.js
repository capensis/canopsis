/**
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
    name: 'WidgetSelectorsUtils',
    after: ['UtilityClass'],
    initialize: function(container, application) {
        var Utility = container.lookupFactory('class:utility');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;


        var widgetSelectors = Utility.create({

            name: 'widgetSelector',

            toTree: function (widget) {
                //doesn't work yet
                if (isNone(widget)) {
                    console.warn('Widget is undefined in widget selector toTree');
                    return {};
                }

                var tree = {};
                console.log('enter tree for widget' , widget.get('id'));
                var widgets = widgetSelectors.directChildren(widget);
                var sub_tree = [];

                var len = widgets.length;
                for (var i = 0, l = len; i < l; i++) {
                    console.log('iterating over widget ', widgets[i].get('id'));
                    sub_tree.push(widgetSelectors.toTree(widgets[i]));
                }

                tree[widget.get('id')] = {
                    children: sub_tree,
                    widget: widget
                };

                console.log('will return sub tree', tree);

                return tree;
            },

            directChildren: function(widget) {

                if (Ember.isNone(widget)) {
                    console.warn('Widget is undefined in widget selector direct children');
                    return [];
                }

                console.log('in direct children');

                var widgets = get(widget, 'items.content');
                var result = [];

                if(!Ember.isNone(widgets)) {
                    var len = widgets.length;
                    for (var i=0; i < len; i++) {
                        result.push(get(widgets[i], 'widget'));
                    }
                }
                return result;
            },

            /**
             * Recursively fetch widgets an returns a list of chilren
             *
             * @param widget
             */
            children: function(widget) {
                if (Ember.isNone(widget)) {
                    console.warn('Widget is undefined in widget selector in children');
                    return [];
                }

                var all = [];
                var widgets = widgetSelectors.directChildren(widget);

                var len = widgets.length;
                for (var i=0; i < len; i++) {
                    console.log('iterating over widget ',i);

                    var children = widgetSelectors.children(widgets[i]);

                    var len1 = children.length;
                    for(var j = 0, lj = len1; j < lj;j++){
                        console.log('iterating over sub children', j );
                        all.push(children[j]);
                    }

                    all.push(widgets[i]);
                }

                return all;
            },

            filter: function(widget, filter) {
                //testing if parameters exists
                if (Ember.isNone(widget)) {
                    console.warn('Widget is undefined in widget selector in filter');
                    return [];
                }
                if (Ember.isNone(filter)) {
                    console.warn('Filter is undefined in widget selector in filter');
                    return [];
                }

                //getting widget list for filtering purposes
                var widgets = [];

                if (filter.directChilrenOnly) {
                    widgets = widgetSelectors.directChildren(widget);
                    console.log('directChilrenOnly', widgets);
                } else {
                    widgets = widgetSelectors.children(widget);
                    console.log('not directChilrenOnly', widgets);
                }

                //keep track of selected element
                var selectedKeys = {};

                //widget filtered list
                var selection = [];

                if (filter.keyEquals) {
                    for (var key in filter.keyEquals) {
                        console.log('checking if key', key, 'is in widget');
                        for (var i = 0, l = widgets.length; i < l; i++) {
                            if (widgets[i].get(key) === filter.keyEquals[key]) {
                                if (Ember.isNone(selectedKeys[get(widgets[i],'id')])) {
                                    selection.push(widgets[i]);
                                    selectedKeys[get(widgets[i], 'id')] = true;
                                }
                            }
                        }
                    }
                    //TODO implement Key exist feature
                }

                return selection;

            }
        });

        application.register('utility:widget-selectors', widgetSelectors);
    }
});
