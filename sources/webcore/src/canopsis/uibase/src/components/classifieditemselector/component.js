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

//TODO fuzzy search
//TODO hover effect


Ember.Application.initializer({
    name: 'component-classifieditemselector',
    initialize: function(container, application) {
        var set = Ember.set,
            get = Ember.get,
            isNone = Ember.isNone;

        /**
         * @component classifieditemselector
         */
        var component = Ember.Component.extend({

            multiselect: true,
            showselection: true,

            init: function(){
                this._super.apply(this, arguments);
                set(this, 'allCollapsed', true);

                if(!Ember.isArray(get(this, 'selection'))) {
                    console.warn('override selection property');
                    set(this, 'selection', Ember.A());
                }
            },

            actions: {
                setListMode: function() {
                    set(this, 'mode', 'list');
                },

                setIconMode: function() {
                    set(this, 'mode', 'icon');
                },

                unselectItem: function (item) {
                    get(this, 'selection').removeObject(item);

                    if(get(this, 'target')) {
                        get(this, 'target').send('unselectItem', item.name);
                    }
                },

                selectItem: function(item) {
                    console.log('selectItem', arguments);
                    if(get(this, 'multiselect') === false) {
                        set(this, 'selection', [item]);
                    } else {
                        //TODO use searchmethodsregistry instead of plain old static code
                        var search = get(this, 'selection').filter(function(loopItem, index, enumerable){
                            void(enumerable);

                            return loopItem === item;
                        });
                        if(search.length === 0){
                            get(this, 'selection').pushObject(item);
                        }
                    }

                    if(get(this, 'target')) {
                        get(this, 'target').send('selectItem', item.name);
                    } else {
                        console.warn('no target attribute for Classifieditemselector', this);
                    }
                },

                collapse: function(theClass){
                    if(theClass === 'all') {
                        if(get(this, 'allCollapsed') === true)
                            set(this, 'allCollapsed', false);
                        else
                            set(this, 'allCollapsed', true);
                    } else if(theClass === 'selection') {
                        if(get(this, 'selectionCollapsed') === true)
                            set(this, 'selectionCollapsed', false);
                        else
                            set(this, 'selectionCollapsed', true);
                    } else {
                        var originClass = get(this, 'classes').findBy('key', theClass.key);

                        console.log('collapse', theClass, theClass.key, originClass);

                        if(originClass.isCollapsed === true ||originClass.isCollapsed === undefined) {
                            set(originClass, 'isCollapsed', false);
                            set(theClass, 'isCollapsed', false);
                        } else {
                            set(originClass, 'isCollapsed', true);
                            set(theClass, 'isCollapsed', true);
                        }
                    }
                }
            },

            searchFilter: '',

            allCollapsed: false,
            selectionCollapsed: false,
            classesCollapsed: true,

            mode: 'list',

            defaultIcon: 'unchecked',

            iconModeButtonCssClass: function(){
                if(get(this, 'mode') === 'icon')
                    return 'btn btn-default active';
                else
                    return 'btn btn-default';
            }.property('mode'),

            listModeButtonCssClass: function(){
                if(get(this, 'mode') === 'list')
                    return 'btn btn-default active';
                else
                    return 'btn btn-default';
            }.property('mode'),

            listGroupClass: function() {
                return 'list-group ' + get(this, 'mode');
            }.property('mode'),

            classAllPanelId: function(){
                return get(this, 'elementId') + '_' + 'all';
            }.property(),

            classAllPanelHref: function(){
                return '#' + get(this, 'classAllPanelId');
            }.property(),

            allClasses: function() {
                var searchFilter = get(this, 'searchFilter');
                var res = get(this, 'content.all');

                if(!isNone(res)) {
                    var component = this;
                    res = get(this, 'content.all').filter(function(item, index, enumerable){
                        void(enumerable);

                        console.log('filter', item);
                        var systemClass = component.get('content.byClass.system');

                        if(!isNone(systemClass)) {
                            if(systemClass.filterBy('name', item.get('name')).length > 0) {
                                console.log('filtered!', item);
                                return false;
                            }
                        }

                        return true;
                    });

                    //TODO use searchmethodsregistry instead of plain old static code
                    if(searchFilter !== '') {
                        res = res.filter(function(item, index, enumerable){
                            void(enumerable);

                            var doesItStartsWithSearchFilter = item.name.slice(0, searchFilter.length) == searchFilter;
                            return doesItStartsWithSearchFilter;
                        });
                    }
                }


                console.log('recompute allClasses', res);
                return res;
            }.property('searchFilter'),

            collapsedPanelCssClass: 'list-group collapse',
            expandedPanelCssClass: 'list-group',

            classList: function(){
                var contentByClass = get(this, 'content.byClass');
                console.log('contentByClass', contentByClass);

                if(contentByClass !== null && contentByClass !== undefined) {
                    return $.map(contentByClass, function(value, key) {
                        if (contentByClass.hasOwnProperty(key)) {
                            return key;
                        }
                    });
                }
            }.property('content'),

            classesFiltered: function(){
                var classes = get(this, 'classes');
                var searchFilter = get(this, 'searchFilter');

                var res = Ember.A();

                //TODO use searchmethodsregistry instead of plain old static code
                var filterFunction = function(item, index, enumerable) {
                    void(enumerable);

                    var doesItStartsWithSearchFilter = item.name.indexOf(searchFilter) !== -1;
                    return doesItStartsWithSearchFilter;
                };

                for (var i = 0, l = classes.length; i < l; i++) {
                    var currentClass = Ember.Object.create({
                        key: classes[i].key,
                        items: classes[i].items,
                        id: classes[i].id,
                        titleHref: classes[i].titleHref,
                        isCollapsed: classes[i].isCollapsed || get(this, 'classesCollapsed')
                    });

                    var classItems = currentClass.items;
                    console.log('classItems', classItems);
                    if(searchFilter !== '') {
                        console.log('filter', classItems);
                        classItems = classItems.filter(filterFunction);

                        currentClass.items = classItems;

                        console.log('filterEnd', classItems);
                    }
                    res.push(currentClass);
                }
                console.log('classesFiltered CP end');
                return res;

            }.property('classes', 'searchFilter'),

            classes: function(){
                var component = this;
                var contentByClass = get(this, 'content.byClass');

                console.log('classes CP', arguments, this, contentByClass);

                if(contentByClass !== undefined) {
                    return $.map(contentByClass, function(value, key) {
                        if (contentByClass.hasOwnProperty(key)) {
                            var newObject = Ember.Object.create({
                                key: key,
                                items: value,
                                id: get(component, 'elementId') + '_' + key,
                                titleHref: '#' + get(component, 'elementId') + '_' + key
                            });
                            var res = [newObject];

                            set(res, 'isCollapsed', contentByClass[key].isCollapsed);

                            return res;
                        }
                    });
                } else {
                    set(component, 'allCollapsed', true);
                    return [];
                }
            }.property('content')
        });

        application.register('component:component-classifieditemselector', component);
    }
});
