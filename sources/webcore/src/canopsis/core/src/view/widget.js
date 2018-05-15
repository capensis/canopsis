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
    name: 'WidgetView',
    after: ['WidgetsRegistry', 'MixinFactory', 'WidgetController', 'MixinsRegistry'],
    initialize: function(container, application) {
        var widgetsregistry = container.lookupFactory('registry:widgets');
        var MixinFactory = container.lookupFactory('factory:mixin');
        var WidgetController = container.lookupFactory('controller:widget');
        var mixinsregistry = container.lookupFactory('registry:mixins');
        var schemasregistry = window.schemasRegistry;

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            widgetsregistry;

        function computeMixinsArray(view, widget) {
            var mixinsNames = get(widget, 'mixins');

            var mixinArray = [];

            console.log('computeMixinsArray', mixinsNames, widget);

            var mixinOptions = {};

            if(mixinsNames) {
                for (var i = 0, l = mixinsNames.length; i < l; i++) {
                    var currentName = mixinsNames[i];

                    //DEPRECATE handle progressive deprecation of mixins as strings
                    if(typeof currentName === 'string') {
                        Ember.deprecate('Defining mixins as strings is deprecated. The new format is : \'{ name: "mixinName" }\'. This is required by the mixin options system.');
                    } else {
                        currentName = currentName.name.camelize();
                    }

                    mixinOptions[currentName] = mixinsNames[i];

                    var currentClass = mixinsregistry.getByName(currentName);

                    console.log('find mixin', currentName, currentClass);

                    //merge mixin's userpreferences into the userpref model
                    var mixinModel = schemasregistry.getByName(currentName);
                    if(mixinModel !== undefined) {
                        mixinModel = mixinModel.EmberModel;

                        var mixinUserPreferenceModel = mixinModel.proto().userPreferencesModel;

                        console.log('mixinModel', mixinUserPreferenceModel);
                        var mixinUserPreferenceModelAttributes = get(mixinUserPreferenceModel, 'attributes');
                        console.log('mixinModelAttributes', mixinUserPreferenceModelAttributes);

                        mixinUserPreferenceModelAttributes.forEach(function(item) {
                            widget.userPreferencesModel[item.name] = mixinUserPreferenceModel[item.name];
                            widget.userPreferencesModel.attributes.add(item);
                        });
                    }

                    if(currentClass) {
                        mixinArray.pushObject(currentClass.EmberClass);
                    } else {
                        get(view, 'displayedErrors').pushObject('mixin not found : ' + currentName);
                        console.error('mixin not found', currentName);
                    }
                }
                var controller = view.get('controller');

                if(controller.onMixinsApplied) {
                    controller.onMixinsApplied();
                }
            }

            mixinArray.pushObject(Ember.Evented);

            return {array: mixinArray, mixinOptions: mixinOptions};
        }

        /**
         * @class WidgetView
         * @extends Ember.View
         * @constructor
         */
        var view = Ember.View.extend({
            templateName:'widget',
            classNames: ['widget'],

            /**
             * Used to visually display error messages to the user (in the widget template)
             *
             * @property widgetController
             * @type Array
             */
            errorMessages : Ember.A(),

            /**
             * @property widgetController
             */
            widgetController: undefined,

            /**
             * @method init
             */
            init: function() {
                console.group('widget initialisation :', get(this, "widget.xtype"), this.widget, get(this, 'widget.tagName'));
                set(this, 'target', get(this, 'controller'));

                this._super();

                set(this, 'displayedErrors', Ember.A());
                if (!! get(this, 'widget')) {
                    this.intializeController(this.widget);
                    this.applyAllViewMixins();
                } else {
                    console.error("No correct widget found for view", this);
                    this.errorMessages.pushObject('No correct widget found');
                }
                if(get(this, 'widget.tagName')) {
                    console.log('custom tagName', get(this, 'widget.tagName'));
                    set(this, 'tagName', get(this, 'widget.tagName'));
                }

                var cssClasses = get(this, 'widget.cssClass');
                if(cssClasses) {
                    console.log('custom tagName', get(this, 'widget.tagName'));
                    set(this, 'classNames', cssClasses.split(','));
                }

                console.groupEnd();
            },

            /**
             * @method applyAllViewMixins
             */
            applyAllViewMixins: function(){
                var controller = get(this, 'controller');
                console.group('apply widget view mixins', controller.viewMixins);
                if(controller.viewMixins !== undefined) {
                    for (var i = 0, mixinsLength = controller.viewMixins.length; i < mixinsLength; i++) {
                        var mixinToApply = controller.viewMixins[i];

                        console.log('mixinToApply', mixinToApply);

                        if(mixinToApply.mixins[0].properties.init !== undefined) {
                            console.warn('The mixin', mixinToApply, 'have a init method. This practice is not encouraged for view mixin as they are applied at runtime and the init method will not be triggerred');
                        }

                        mixinToApply.apply(this);
                    }
                }
                console.groupEnd();
            },

            /**
             * @method intializeController
             */
            intializeController: function(widget) {
                console.group('set controller for widget', widget);

                var controller = this.instantiateCorrectController(widget);

                var widgetTemplate = get(widget, "xtype");

                if(widgetTemplate === 'text') widgetTemplate = 'textwidget';

                this.setProperties({
                    'controller': controller,
                    'templateName': widgetTemplate
                });

                widget.set('controller', controller);

                this.registerHooks();
                console.groupEnd();
            },

            /**
             * @method instantiateCorrectController
             * @param {DS.Model} widget
             * @return WidgetController
             */
            instantiateCorrectController: function(widget) {
                //for a widget that have xtype=widget, controllerName=WidgetController
                console.log('instantiateCorrectController', arguments);
                var xtype = get(widget, "xtype");
                if(xtype === undefined || xtype === null) {
                    console.error('no xtype for widget', widget, this);
                }

                var mixins = computeMixinsArray(this, widget);

                mixins.array.pushObject({
                    model: widget,
                    target: get(this, 'target')
                });

                var widgetControllerInstance;

                var widgetClass = widgetsregistry.getByName(get(widget, "xtype"));

                if(widgetClass !== undefined) {
                    widgetClass = widgetClass.EmberClass;
                } else {
                    widgetClass = WidgetController;
                }

                widgetControllerInstance = widgetClass.createWithMixins.apply(widgetClass, mixins.array);
                widgetControllerInstance.refreshPartialsList();

                Ember.setProperties(widgetControllerInstance, {
                    'model.displayedErrors': get(this, 'displayedErrors'),
                    'mixinOptions': mixins.mixinOptions
                });

                widgetControllerInstance.mixinsOptionsReady();

                var mixinsName = get(widget, 'mixins');

                if (mixinsName) {
                    for (var i = 0, l = mixinsName.length; i < l ; i++ ){
                        var currentName =  mixinsName[i];
                        var currentMixin = mixinsregistry.all[currentName];

                        if (currentMixin) {
                            currentMixin.apply(widgetControllerInstance);
                        }
                    }
                }

                return widgetControllerInstance;
            },

            /**
             * @method didInsertElement
             */
            didInsertElement : function() {
                console.log("inserted widget, view:", this);

                this.registerHooks();

                return this._super.apply(this, arguments);
            },

            /**
             * @method willDestroyElement
             */
            willDestroyElement: function () {
                clearInterval(get(this, 'widgetRefreshInterval'));

                this.unregisterHooks();

                return this._super.apply(this, arguments);
            },

            onWidgetRefresh: function() {},

            /**
             * @method registerHooks
             */
            registerHooks: function() {
                console.log("registerHooks", get(this, "controller"), get(this, "controller").on);
                get(this, "controller").on('refresh', this, this.rerender);
                return this._super();
            },

            /**
             * @method unregisterHooks
             */
            unregisterHooks: function() {
                get(this, "controller").off('refresh', this, this.rerender);
                return this._super();
            }
        });
        widgetsregistry = container.lookupFactory('registry:widgets');

        Ember.Handlebars.helper('widgethelper', view);

        application.register('view:widget', view);
    }
});

