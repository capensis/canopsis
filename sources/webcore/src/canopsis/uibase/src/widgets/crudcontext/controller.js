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
    name: 'CrudContextWidget',
    after: ['WidgetFactory', 'PaginationMixin', 'InspectableArrayMixin', 'UserconfigurationMixin', 'DraggablecolumnsMixin', 'RoutesUtils', 'FormsUtils', 'ListlineView', 'DataUtils'],
    initialize: function(container, application) {
        var WidgetFactory = container.lookupFactory('factory:widget');
        var PaginationMixin = container.lookupFactory('mixin:pagination');
        var InspectableArrayMixin = container.lookupFactory('mixin:inspectable-array');
        var UserConfigurationMixin = container.lookupFactory('mixin:userconfiguration');
        var ListlineController = container.lookupFactory('view:listline');
        var DataUtils = container.lookupFactory('utility:data');

        var formsUtils = container.lookupFactory('utility:forms');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;


        var listOptions = {
            mixins: [
                InspectableArrayMixin,
                PaginationMixin,
                UserConfigurationMixin
            ]
        };

        /**
         * @widget List
         *
         * @description
         *
         * # Overview
         *
         * Displays a list of records. This widget can be enhanced with a wide range of mixins, to provide additionnal content management features such as :
         * - filtering
         * - content display
         * - data edition
         * - monitoring related features
         * - and so on
         *
         * # Basic usage
         *
         * By default, the widget is able to display a paginated listing of all records found for a specified type (see the "listed_crecord_type" property on the schema).
         * For all of this records, it will show a configurable list of columns (see the "displayed_columns" schema property).
         *
         * For each displayed column of each listed record, the widget will try to find a correct and fancy way to display the value. It will thus try to find a correct renderer to display the data key.
         *
         * # Screenshots
         *
         * ![Simple list](../screenshots/widget-list-simple.png)
         * ![Event view](../screenshots/widget-list-events.png)
         */
        var widget = WidgetFactory('crudcontext', {
            css :'table table-striped table-bordered table-hover dataTable sortable',
            needs: ['login', 'application', 'recordinfopopup'],
            /**
             * @property standardListDisplay
             * @description Whether to display the list as the regular table or not.
             * Used with mixin that fill the partial slot "alternativeListDisplay", this can help to provide alternative displays
             */
            standardListDisplay: true,

            /**
             * @property dynamicTemplateName
             */
            dynamicTemplateName: 'loading',

            /**
             * @property listlineControllerClass
             * @todo test if this is needed (used in cloaked mode)
             */
            listlineControllerClass: ListlineController,

            /**
             * @property user
             */
            user: Ember.computed.alias('controllers.login.record._id'),

            /**
             * @property rights
             */
            rights: Ember.computed.alias('controllers.login.record.rights'),

            /**
             * @property safeMode
             */
            safeMode: Ember.computed.alias('controllers.application.frontendConfig.safe_mode'),

            /**
             * @method init
             */
            init: function() {
                var me = this;
                //prepare user configuration to fetch customer preference by reseting data.
                //properties are set here to avoid same array references between instances.
                Ember.setProperties(this, {
                    custom_filters: [],
                    widgetData: [],
                    findOptions: {},
                    loaded: false
                });
                crudContextAdapter = DataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:crudcontext');
                set(me, 'crudContextAdapter', crudContextAdapter);
                set(me, 'listed_crecord_type', 'ccontext');

                get(this,'partials.selectionToolbarButtons').pushObject('actionbutton-createpbehavior');

                this._super.apply(this, arguments);

            },

            /**
             * @method rollbackableChanged
             * @description observes if the model is rollbackable
             */
            rollbackableChanged: function() {
                var list = this;
                if(get(this, 'model.rollbackable') === false) {
                    Ember.run.scheduleOnce('afterRender', this, function() { list.refreshContent(); });
                }
            }.observes('model.rollbackable'),

            /**
             * @method generateListlineTemplate
             * @argument shown_columns
             */
            generateListlineTemplate: function (shown_columns) {
                //TODO temporarily removed create button
                get(this,'partials.actionToolbarButtons').removeObject('actionbutton-create');

                var html = '<td>{{#if pendingOperation}}<i class="fa fa-cog fa-spin"></i>{{/if}}{{component-checkbox checked=isSelected class="toggle"}}</td>';

                if(get(this, '_partials.columnsLine')) {
                    html += '{{#each columns in controller._partials.columnsLine}}<td>{{partial columns}}</td>{{/each}}';
                }

                if(shown_columns === undefined || shown_columns.length === 0) {
                    return undefined;
                }

                for (var i = 0, l = shown_columns.length; i < l; i++) {
                    var currentColumn = shown_columns[i];
                    console.log('currentColumn', currentColumn);

                    if(get(currentColumn, 'options.show')) {
                        if(currentColumn.renderer && get(this, 'model.useRenderers')) {
                            html += ['<td class="', currentColumn.field, '">{{component-renderer rendererType="', currentColumn.renderer, '" value=this.', currentColumn.field, ' record=this field="', currentColumn.field, '" shown_columns=controller.shown_columns}}</td>'].join('');
                        } else {
                            html += ['<td class="', currentColumn.field, '">{{this.', currentColumn.field, '}}</td>'].join('');
                        }
                    }
                }

                var itemActionbuttons = get(this, '_partials.itemactionbuttons');
                if(itemActionbuttons) {
                    console.log('itemactionbuttons', itemActionbuttons);
                    html += '<td style="padding-left:0; padding-right:0"><div class="btn-group" style="display:flex">';

                    for (var j = 0, lj = itemActionbuttons.length; j < lj; j++) {
                        html += ['{{partial "', itemActionbuttons[j], '"}}'].join('');
                    }

                    html += '</div></td>';
                }

                console.log('generatedListlineTemplate', html);
                return html;
            },

            /**
             * @method updateInterval
             * @description Manages how time filter is set to the widget
             * @argument interval
             */
            updateInterval: function (interval){
                console.warn('Set widget list time interval', interval);
                set(this, 'timeIntervalFilter', interval);
                this.refreshContent();
            },

            /**
             * @method getTimeInterval
             * @description Manages how time filter is get from the widget for refresh purposes
             */
            getTimeInterval: function () {
                var interval = get(this, 'timeIntervalFilter');
                if (isNone(interval)) {
                    return {};
                } else {
                    return interval;
                }
            },

            /**
             * @property itemType
             */
            itemType: function() {
                var listed_crecord_type = get(this, 'listed_crecord_type');
                console.info('listed_crecord_type', listed_crecord_type);
                if(listed_crecord_type !== undefined && listed_crecord_type !== null ) {
                    return get(this, 'listed_crecord_type');
                } else {
                    return 'event';
                }
            }.property('listed_crecord_type'),

            /**
             * @method isAllSelectedChanged
             * @description Observer on "isAllSelected". Automatically select every record when the header checkbox is clicked
             */
            isAllSelectedChanged: function(){
                get(this, 'widgetData').content.setEach('isSelected', get(this, 'isAllSelected'));
            }.observes('isAllSelected'),

            //Mixin aliases
            //history

            /**
             * @property historyMixinFindOptions
             */
            historyMixinFindOptions: Ember.computed.alias('findOptions.useLogCollection'),
            //inspectedDataItemMixin

            /**
             * @property inspectedDataArray
             */
            inspectedDataArray: Ember.computed.alias('widgetData'),
            //pagination

            /**
             * @property paginationMixinFindOptions
             */
            paginationMixinFindOptions: Ember.computed.alias('findOptions'),

            /**
             * @property widgetDataMetas
             */
            widgetDataMetas: {},

            /**
             * @method findItems
             */
            findItems: function() {
                var me = this;
                var controller = this;
                var crudContextAdapter = get(this, 'crudContextAdapter');

                var appController = get(this, 'controllers.application');
                appController.addConcurrentLoading('list-data');
                var store = get(controller, 'widgetDataStore'),
                    promise;

                if (isNone(get(controller, 'widgetDataStore'))) {
                    set(controller, 'widgetDataStore', DS.Store.create({
                        container: get(controller, 'container')
                    }));
                }

                var itemType = get(this, 'itemType');

                var findParams = this.computeFindParams();

                if(isNone(container.lookupFactory('model:' + itemType))) {
                    if(get(this, 'content.displayedErrors') === undefined) {
                        set(this, 'content.displayedErrors', Ember.A());
                    }

                    get(this, 'content.displayedErrors').pushObject('No model was found for '+ itemType);
                    set(this, 'loaded', true);
                    appController.removeConcurrentLoading('list-data');
                } else {
                    get(this, 'widgetDataStore').findQuery(itemType, findParams).then(function(queryResults) {
                        //retreive the metas of the records
                        var listed_crecord_type = get(me, 'listed_crecord_type');
                        var crecordTypeMetadata = get(me, 'widgetDataStore').metadataFor(listed_crecord_type);
                        console.log('crecordTypeMetadata', crecordTypeMetadata);

                        Ember.setProperties(me, {
                            widgetDataMetas: crecordTypeMetadata,
                            loaded: true
                        });

                        me.extractItems.apply(me, [queryResults]);

                        for(var i = 0, l = queryResults.content.length; i < l; i++) {
                            //This value reset spiner display for record in flight status
                            queryResults.content[i].set('pendingOperation', false);
                        }

                        appController.removeConcurrentLoading('list-data');
                    });
                }
            },

            /**
             * @event widgetDataChanged
             */
            widgetDataChanged: function() {
                this.trigger('refresh');
            }.observes('widgetData'),

            /**
             * @property attributesKeysDict
             * @description Computed property dependant on "attributesKeys"
             */
            attributesKeysDict: function() {
                var res = {};
                var attributesKeys = get(this, 'attributesKeys');
                var sortedAttribute = get(this, 'sortedAttribute');

                for (var i = 0, l = attributesKeys.length; i < l; i++) {
                    if (sortedAttribute !== undefined && sortedAttribute.field === attributesKeys[i].field) {
                        sortedAttribute.renderer = this.rendererFor(sortedAttribute);
                        res[attributesKeys[i].field] = sortedAttribute;
                    } else {
                        attributesKeys[i].renderer = this.rendererFor(attributesKeys[i]);
                        res[attributesKeys[i].field] = attributesKeys[i];
                    }
                }

                return res;
            }.property('attributesKeys'),

            /**
             * @method redenrerFor
             * @argument attribute
             */
            rendererFor: function(attribute) {
                var type = get(attribute, 'type');
                var role = get(attribute, 'options.role');
                if(get(attribute, 'model.options.role')) {
                    role = get(attribute, 'model.options.role');
                }
                var subRole = get(attribute, 'options.items.role');
                if(role === 'array' && !isNone(subRole)) {
                    role = subRole;
                }

                var rendererName;
                if (role) {
                    rendererName = 'renderer-' + role;
                } else {
                    rendererName = 'renderer-' + type;
                }

                if (Ember.TEMPLATES[rendererName] === undefined) {
                    rendererName = undefined;
                }

                return rendererName;
            },

            /**
             * @property shown_columns
             */
            shown_columns: function() {
                console.log('compute shown_columns', get(this, 'sorted_columns'), get(this, 'sortedAttribute'));

                //user preference for displayed columns.
                if (this.get('user_show_columns') !== undefined) {
                    console.log('user columns selected', get(this, 'user_show_columns'));
                    return get(this, 'user_show_columns');
                }

                var shown_columns = [];
                var displayed_columns = get(this, 'displayed_columns') || get(this, 'model.columns') ;
                if (displayed_columns !== undefined && displayed_columns.length > 0) {

                    var attributesKeysDict = this.get('attributesKeysDict');

                    for (var i = 0, li = displayed_columns.length; i < li; i++) {
                        if (attributesKeysDict[displayed_columns[i]] !== undefined) {
                            set(attributesKeysDict[displayed_columns[i]].options, 'show', true);
                            shown_columns.push(attributesKeysDict[displayed_columns[i]]);
                        }
                    }
                } else {
                    console.log('no shown columns set, displaying everything');

                    shown_columns = this.get('attributesKeys');
                }

                var selected_columns = [];
                var sortedAttribute = get(this, 'sortedAttribute');
                var columnSort = get(this, 'default_column_sort');
                for(var column=0, l = shown_columns.length; column < l; column++) {
                    //reset previous if any in case list configuration is updated
                    if (shown_columns[column].options.canUseDisplayRecord) {
                        delete shown_columns[column].options.canUseDisplayRecord;
                    }

                    //set option display record field to true allow list line template to change renderer
                    //diusplay and if true, an action can be triggrered from trusted column.
                    if (shown_columns[column].field === get(this, 'display_record_field')) {
                        shown_columns[column].options.canUseDisplayRecord = true;
                    }

                    //Manage hidden colums from the list parameters information.
                    //If colname exists in hidden_column list, then it is not displayed.
                    if ($.inArray(shown_columns[column].field, get(this, 'hidden_columns')) === -1) {
                        selected_columns.pushObject(shown_columns[column]);
                    }

                    //Manage sort icon from default sort
                    if (!isNone(columnSort) &&
                        columnSort.property === shown_columns[column].field &&
                        !isNone(columnSort.direction) &&
                        (isNone(sortedAttribute) || sortedAttribute === {})) {
                        var headerClass = columnSort.direction === 'ASC' ? 'sorting_asc' : 'sorting_desc';
                        shown_columns[column].headerClassName = headerClass;
                    }

                    //select appropriate title for column headers from shema options.
                    var display_title = 'no label';
                    var label = get(shown_columns[column], 'options.label');
                    var title = get(shown_columns[column], 'options.title');
                    var field = get(shown_columns[column], 'field');
                    if (title) {
                        display_title = title;
                    } else if (label) {
                        display_title = label;
                    } else if (field) {
                        display_title = field;
                    }
                    set(shown_columns[column], 'display_title', display_title);

                }

                console.debug('selected cols', selected_columns);

                var hbs = this.generateListlineTemplate(selected_columns);

                if(hbs !== undefined) {
                    var tpl = Ember.Handlebars.compile(hbs);
                    set(this, 'hbsListline', hbs);

                    var oldTemplateName = get(this, 'dynamicTemplateName');

                    if(oldTemplateName) {
                        delete Ember.TEMPLATES[oldTemplateName];
                    }

                    var dynamicTemplateName = 'dynamic-list' + Math.floor(Math.random() * 10000);

                    Ember.TEMPLATES[dynamicTemplateName] = tpl;
                    set(this, 'dynamicTemplateName', dynamicTemplateName);
                }

                return selected_columns;

            }.property('attributesKeysDict', 'sorted_columns'),

            /**
             * @method computeFilterFragmentsList
             * @description Computes the list of different filter fragments used to create a proper query
             * @returns {Array} the list of fragments
             */
            computeFilterFragmentsList: function() {
                var list = Ember.A();
                var additionalFilterPart = get(this, 'model.volatile.forced_filter') || get(this, 'additional_filter');
                console.log('additionalFilterPart', get(this, 'model.volatile.forced_filter'), get(this, 'additional_filter'));

                list.pushObject(this.getTimeInterval());
                list.pushObject(additionalFilterPart);

                return list;
            },

            /**
             * @method computeFindParams
             */
            computeFindParams: function(){
                console.group('computeFindParams', get(this, 'model.selected_filter.filter'));

                var filterFragments = this.computeFilterFragmentsList();

                var filters = [];

                for (var i = 0, l = filterFragments.length; i < l; i++) {
                    if(typeof filterFragments[i] === 'string') {
                        //if json, parse json
                        try {
                            filterFragments[i] = JSON.parse(filterFragments[i]);
                        } catch (e) {
                            if(get(this, 'content.displayedErrors') === undefined) {
                                set(this, 'content.displayedErrors', Ember.A());
                            }

                            get(this, 'content.displayedErrors').pushObject('There seems to be an error with the currently selected filter.');
                            filterFragments[i] = {};
                        }
                    }
                    if (filterFragments[i] !== {} && !isNone(filterFragments[i])) {
                        //when defined filter then it is added to the filter list
                        filters.pushObject(filterFragments[i]);
                    }
                }

                var params = {};

                params.limit = get(this, 'itemsPerPagePropositionSelected');

                //TODO check if useless or not
                if(params.limit === 0) {
                    params.limit = 5;
                }

                params.start = get(this, 'paginationFirstItemIndex') - 1;


                if (filters.length) {
                    params._filter = JSON.stringify({ '$and': filters });
                }

                var userSortedAttribute = get(this, 'model.user_sortedAttribute');

                if (isNone(userSortedAttribute)) {
                    if(!isNone(get(this, 'model.default_column_sort'))) {
                        params.sort = JSON.stringify([ get(this, 'model.default_column_sort') ]);
                        console.log('defaultSortedAttribute', get(this, 'model.default_column_sort'));
                    }

                } else if (!isNone(userSortedAttribute)) {
                    var direction;
                    if(get(userSortedAttribute, 'headerClassName') === 'sorting_asc') {
                        direction = 'ASC';
                    } else {
                        direction = 'DESC';
                    }

                    params.sort = JSON.stringify([{property: userSortedAttribute.field, direction: direction}]);
                    console.log('userSortedAttribute', userSortedAttribute);
                }

                console.groupEnd();

                return params;
            },

            actions: {
                createMultiplePBehaviors: function() {
                    var crudController = this;

                     var selected = crudController.get('widgetData').filterBy('isSelected', true);
                     if(selected.length === 0)
                        return

                    // this.get('widgetData.content')[0]['crecord_type'] = 'ccontext'
                    var obj = Ember.Object.create({'crecord_type': 'pbehaviorform'})
                    var confirmform = formsUtils.showNew('modelform', obj, {
                        title: __('Put a pbehavior on these elements ?')
                    });

                    confirmform.submit.then(function(form) {
                        if(!Array.isArray(selected))
                            selected = [selected]

                        var listId = []

                        for (var i = 0, l = selected.length; i < l; i++)
                            listId.push(selected[i]['_data']['_id'])

                        var payload = {
                            'tstart':form.get('formContext.start'),
                            'tstop': form.get('formContext.type_') === "Pause" ? 2147483647 : form.get('formContext.end'),
                            'rrule':form.get('formContext.rrule'),
                            'name':form.get('formContext.name'),
                            'type_': form.get('formContext.type_'),
                            'reason': form.get('formContext.reason'),
                            'author': window.username,
                            'filter': {}
                        }

                        if(!payload.rrule)
                            delete(payload.rrule)

                        payload.filter = {
                            '_id': {
                                '$in': listId
                            }
                        }

                        //$.post(url)
                        return $.ajax({
                            type: 'POST',
                            url: '/api/v2/pbehavior',
                            data: JSON.stringify(payload),
                            contentType: 'application/json',
                            dataType: 'json',
                            success: function () {
                                console.log("Ticket sent");
                            },
                            error: function () { console.error("Failure to send downtime") }
                        });

                    });
                }
            }
        }, listOptions);

        application.register('widget:crudcontext', widget);

    }
});
