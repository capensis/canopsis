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
    name: 'WidgetController',
    after: ['PartialslotAbleController', 'WidgetsUtils', 'RoutesUtils', 'FormsUtils', 'DebugUtils', 'DataUtils', 'HashUtils', 'WidgetsRegistry'],
    initialize: function(container, application) {
        var PartialslotAbleController = container.lookupFactory('controller:partialslot-able');

        var widgetUtils = container.lookupFactory('utility:widgets');
        var routesUtils = container.lookupFactory('utility:routes');
        var formsUtils = container.lookupFactory('utility:forms');
        var debugUtils = container.lookupFactory('utility:debug');
        var dataUtils = container.lookupFactory('utility:data');
        var hashUtils = container.lookupFactory('utility:hash');
        var schemasregistry = window.schemasRegistry;
        var WidgetsRegistry = container.lookupFactory('registry:widgets');
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @class WidgetController
         * @extends PartialslotAbleController
         * @constructor
         */
        var controller = PartialslotAbleController.extend({
            needs: ['application', 'login'],

            partials: {
                titlebarsbuttons : ['titlebarbutton-duplicate', 'titlebarbutton-moveup','titlebarbutton-movedown', 'titlebarbutton-widgeterrors']
            },

            /**
             * This is useful mostly for debug, to know that a printend object is a widget
             *
             * @property abstractType
             * @type string
             */
            abstractType: 'widget',

            /**
             * @property canopsisConfiguration
             * @type Object
             */
            canopsisConfiguration: canopsisConfiguration,

            /**
             * true is the Frontend is in debug mode
             *
             * @property debug
             * @type boolean
             */
            debug: Ember.computed.alias('canopsisConfiguration.DEBUG'),

            /**
             * @property editMode
             * @type boolean
             */
            editMode : Ember.computed.alias('controllers.application.editMode'),

            /**
             * Alias for content
             *
             * @property config
             * @deprecated
             * @type DS.Model
             */
            config: Ember.computed.alias('content'),

            /**
             * @method init
             */
            init: function () {

                set(this, 'userParams', {});

                console.log('widget init');

                var viewId = get(widgetUtils.getParentViewForWidget(this), 'content.id');
                var container = routesUtils.getCurrentRouteController().container;
                var store = DS.Store.create({
                    container: container
                });

                console.debug('View id for current widget is ', viewId);
                this.setProperties({
                    'model.controllerInstance': this,
                    'viewId': viewId,
                    'viewController': widgetUtils.getParentViewForWidget(this),
                    'isOnMainView': get(widgetUtils.getParentViewForWidget(this), 'isMainView'),
                    'container': container,
                    'widgetDataStore': store
                });

                //User preference are called just before the refresh to ensure
                //refresh takes care of user information and widget general preference is overriden
                //All widget may not have this mixin, so it's existance is tested
                if (!isNone(this.loadUserConfiguration)) {
                    this.loadUserConfiguration();
                }

                this._super();

                console.debug('user configuration loaded for widget ' + get(this, 'title'));
                this.refreshContent();
            },


            /**
             * @method mixinsOptionsReady
             */
            mixinsOptionsReady: function () {
                //can be overriden to trigger action when mixins options ready.
            },

            /**
             * Adds mixins view to the current widget controller
             *
             * @method addMixinView
             * @param viewMixin
             */
            addMixinView: function (viewMixin) {
                var viewMixins = get(this, 'viewMixins');
                if (isNone(viewMixins)) {
                    viewMixins = [];
                    set(this, 'viewMixins', viewMixins);
                }
                viewMixins.push(viewMixin);
            },

            /**
             * @method updateInterval
             * @param interval
             */
            updateInterval: function (interval) {
                console.warn('This method should be overriden for current widget', get(this, 'id'), interval);
            },

            /**
             * @method stopRefresh
             */
            stopRefresh: function () {
                set(this, 'isRefreshable', false);
            },

            /**
             * @method startRefresh
             */
            startRefresh: function () {
                this.setProperties({
                    'isRefreshable': true,
                    'lastRefresh': null
                });
            },

            /**
             * @method isRollbackable
             */
            isRollbackable: function() {
                if(get(this, 'isDirty') && get(this, 'dirtyType') === 'updated' && get(this, 'rollbackable') === true) {
                    return true;
                }

                return false;
            }.property('isDirty', 'dirtyType', 'rollbackable'),


            actions: {

                /**
                 * Show debug info in console and put widget var in window.$E
                 *
                 * @event inspect
                 * @param object
                 */
                inspect: function (object) {
                    debugUtils.inspectObject(object);
                },

                /**
                 * @event do
                 * @param action
                 */
                do: function(action) {
                    var params = [];

                    for (var i = 1, l = arguments.length; i < l; i++) {
                        params.push(arguments[i]);
                    }

                    this.send(action, params);
                },

                /**
                 * @event creationForm
                 * @param itemType
                 */
                creationForm: function(itemType) {
                    formsUtils.addRecord(itemType);
                },

                /**
                 * @event rollback
                 * @param widget
                 */
                rollback: function(widget){
                    console.log('rollback changes', arguments);
                    set(widget, 'volatile', {});
                    widget.rollback();
                    set(widget, 'rollbackable', false);
                },

                /**
                 * @event editWidget
                 * @param widget
                 */
                editWidget: function (widget) {
                    console.info('edit widget', widget);

                    var widgetTitle = get(widget, 'title') || '';
                    var widgetType = get(widget, 'xtype') || '';

                    var widgetWizard = formsUtils.showNew(
                        'modelform',
                        widget,
                        { title: __('Edit widget') + ' ' + widgetType + ' ' + widgetTitle}
                    );
                    console.log('widgetWizard', widgetWizard);

                    var widgetController = this;

                    widgetWizard.submit.done(function() {
                        console.log('record going to be saved', widget);

                        if(!get(widget, 'widgetId')) {
                            set(widget, 'widgetId', get(widget,'id'));
                        }

                        var userview = get(widget, 'controller.viewController').get('content');

                        userview.save().then(function(){
                            get(widgetController, 'viewController').send('refresh');
                            console.log('triggering refresh', userview);
                        });
                    });
                },

                /**
                 * @event editMixin
                 * @param widget
                 * @param mixinName
                 */
                editMixin: function (widget, mixinName) {
                    console.info('edit mixin', widget, mixinName);

                    var mixinDict = get(widget, 'mixins').findBy('name', mixinName);

                    if(!Ember.isNone(mixinDict)) {
                        mixinDict.id = hashUtils.generateId('mixin');
                    }

                    var mixinModelInstance = dataUtils.getStore().createRecord(mixinName, mixinDict);

                    var mixinForm = formsUtils.showNew('modelform', mixinModelInstance, { title: __('Edit mixin'), inspectedItemType: mixinName });

                    var mixinObject = get(widget, 'mixins').findBy('name', mixinName);

                    if(isNone(mixinObject)) {
                        mixinObject = get(widget, 'mixins').pushObject({name: mixinName});
                    }

                    console.log('mixinObject', mixinObject);

                    var widgetController = this;

                    mixinForm.submit.done(function() {
                        var referenceModel = schemasregistry.getByName(mixinName).EmberModel;
                        var modelAttributes = get(referenceModel, 'attributes');

                        console.log('attributes', modelAttributes);

                        modelAttributes.forEach(function(property) {
                            console.log('each', arguments);
                            var propertyValue = get(mixinModelInstance, property.name);
                            console.log('mixinObject', mixinObject);

                            set(mixinObject, property.name, propertyValue);

                            var userview = get(widgetController, 'viewController').get('content');
                            userview.save().then(function(){
                                get(widgetController, 'viewController').send('refresh');
                                console.log('triggering refresh', userview);
                            });
                        });
                    });
                },

                /**
                 * @event removeWidget
                 * @param widget
                 */
                removeWidget: function (widget) {

                    var widgetController = this;

                    var confirmform = formsUtils.showNew('confirmform', {}, {
                        title: __('Delete this widget ?')
                    });

                    confirmform.submit.then(function(form) {


                        console.group('remove widget', widget);
                        console.log('parent container', widgetController);

                        var itemsContent = get(widgetController, 'content.items.content');

                        for (var i = 0, l = itemsContent.length; i < l; i++) {
                            console.log(get(widgetController, 'content.items.content')[i]);

                            if (get(itemsContent[i], 'widget') === widget) {
                                itemsContent.removeAt(i);
                                console.log('deleteRecord ok');
                                break;
                            }
                        }

                        var userview = get(widgetController, 'viewController.content');

                        userview.save();

                        console.groupEnd();
                    });
                },

                /**
                 * Moves the widget under the next one, if any
                 *
                 * @event movedown
                 * @param widgetwrapper
                 */
                movedown: function(widgetwrapper) {
                    console.group('movedown', widgetwrapper);

                    try{
                        console.log('context', this);

                        var foundElementIndex,
                            nextElementIndex,
                            itemsContent = get(this, 'content.items.content');

                        for (var i = 0, l = itemsContent.length; i < l; i++) {

                            if (foundElementIndex !== undefined && nextElementIndex === undefined) {
                                nextElementIndex = i;
                                console.log('next element found');
                            }

                            if (itemsContent[i] === widgetwrapper) {
                                foundElementIndex = i;
                                console.log('searched element found');
                            }
                        }

                        if (foundElementIndex !== undefined && nextElementIndex !== undefined) {
                            //swap objects
                            var array = itemsContent;
                            console.log('swap objects', array);

                            var tempObject = array.objectAt(foundElementIndex);

                            array.insertAt(foundElementIndex, array.objectAt(nextElementIndex));
                            array.insertAt(nextElementIndex, tempObject);
                            array.replace(foundElementIndex + 2, 2);

                            console.log('new array', array);
                            set(this, 'content.items.content', array);

                            var widgetController = this,
                                userview = get(this, 'viewController.content');

                            userview.save();
                        }
                    } catch (e) {
                        console.error(e.stack, e.message);
                    }
                    console.groupEnd();
                },

                /**
                 * Moves the widget above the previous one, if any
                 *
                 * @event moveup
                 * @param widgetwrapper
                 */
                moveup: function(widgetwrapper) {
                    console.group('moveup', widgetwrapper);

                    try{
                        console.log('context', this);

                        var foundElementIndex,
                            nextElementIndex,
                            itemsContent = get(this, 'content.items.content');

                        for (var i = itemsContent.length; i >= 0 ; i--) {

                            if (foundElementIndex !== undefined && nextElementIndex === undefined) {
                                nextElementIndex = i;
                                console.log('next element found');
                            }

                            if (itemsContent[i] === widgetwrapper) {
                                foundElementIndex = i;
                                console.log('searched element found');
                            }
                        }

                        console.log('indexes to swap', foundElementIndex, nextElementIndex);

                        if (foundElementIndex !== undefined && nextElementIndex !== undefined) {
                            //swap objects
                            var array = get(this, 'content.items.content');
                            console.log('swap objects', array);

                            var tempObject = array.objectAt(foundElementIndex);

                            array.insertAt(foundElementIndex, array.objectAt(nextElementIndex));
                            array.insertAt(nextElementIndex, tempObject);
                            array.replace(nextElementIndex + 2, 2);

                            console.log('new array', array);
                            set(this, 'content.items.content', array);

                            var widgetController = this,
                                userview = get(widgetUtils.getParentViewForWidget(this), 'content');

                            userview.save();
                        }
                    } catch (e) {
                        console.error(e.stack, e.message);
                    }
                    console.groupEnd();
                }
            },


            /**
             * @property itemController
             * @type string
             */
            itemController: function() {
                if(get(this, 'itemType')) {
                    return get(this, 'itemType').capitalize();
                }
            }.property('itemType'),

            /**
             * @method refreshContent
             */
            refreshContent: function() {

                console.log('refreshContent', get(this, 'xtype'));

                this._super();
                this.findItems();

                this.setProperties({
                    'lastRefresh': new Date().getTime(),
                    'lastRefreshControlDelay': true
                });
            },

            /**
             * @method findItems
             */
            findItems: function() {
                console.warn('findItems not implemented', this);
            },

            /**
             * @method extractItems
             * @param queryResult
             */
            extractItems: function(queryResult) {
                console.log('extractItems', queryResult);

                this._super(queryResult);
                set(this, 'widgetData', queryResult);
            },

            /**
             * @property availableTitlebarButtons
             * @type array
             */
            availableTitlebarButtons: function(){
                var buttons = get(this, '_partials.titlebarsbuttons');

                if(buttons === undefined) {
                    return Ember.A();
                }

                var res = Ember.A();

                for (var i = 0, l = buttons.length; i < l; i++) {
                    var currentButton = buttons[i];

                    if(Ember.TEMPLATES[currentButton] !== undefined) {
                        res.pushObject(currentButton);
                    } else {
                        //TODO manage this with utils.problems
                        console.warn('template not found', currentButton);
                    }
                }

                return res;
            }.property()
        });

        application.reopen({
            register: function (name, object) {

                if(name.split(':')[0] === 'widget') {
                    var widgetName = name.split(':')[1];
                    var initializerName = widgetName.capitalize() + 'Serializer';
                    var widgetSerializerName = name.split(':')[1];
                    if(schemasregistry.getByName(widgetSerializerName) === undefined) debugger;
                    var widgetModel = schemasregistry.getByName(widgetSerializerName).EmberModel;

                    Ember.Application.initializer({
                        name: initializerName,
                        after: 'WidgetSerializer',
                        initialize: function(container, application) {
                            var WidgetSerializer = container.lookupFactory('serializer:widget');
                            application.register('serializer:' + widgetSerializerName, WidgetSerializer.extend());
                        }
                    });

                    if (isNone(widgetModel)) {
                        notificationUtils.error('No model found for the widget ' + widgetName + '. There might be no schema concerning this widget on the database');
                    } else {

                        var capitalizedWidgetName = widgetName.camelize().capitalize();
                        var metadataDict = widgetModel.proto().metadata;

                        var registryEntry = Ember.Object.create({
                            name: widgetName,
                            EmberClass: object
                        });

                        if(!isNone(metadataDict)) {
                            if(metadataDict.icon) {
                                registryEntry.set('icon', metadataDict.icon);
                            }
                            if(metadataDict.classes) {
                                var classes = metadataDict.classes;
                                for (var j = 0, lj = classes.length; j < lj; j++) {
                                    var currentClass = classes[j];
                                    if(!Ember.isArray(get( WidgetsRegistry.byClass, currentClass))) {
                                        set(WidgetsRegistry.byClass, currentClass, Ember.A());
                                    }

                                    get(WidgetsRegistry.byClass, currentClass).pushObject(registryEntry);
                                }
                            }
                        }

                        WidgetsRegistry.all.pushObject(registryEntry);
                    }
                }

                return this._super.apply(this, arguments);
            }
        });

        application.register('controller:widget', controller);
    }
});
