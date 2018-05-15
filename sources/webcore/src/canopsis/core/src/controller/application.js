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
    name: 'ApplicationController',
    after: ['PartialslotAbleController',
            'UserprofilestatusmenuMixin',
            'RequirejsmocksmanagerMixin',
            'ScreentoolstatusmenuMixin',
            'DocumentationMixin',
            'SchemamanagerMixin',
            'ConsolemanagerMixin',
            'LoadingindicatorMixin',
            'PromisemanagerMixin',
            'NotificationsMixin',
            'FormsUtils',
            'DataUtils',
            'DebugUtils',
            'HashUtils',
            'NotificationUtils',
            'LoginController',
            'RecordinfopopupController',
            'CrudMixin'],
    initialize: function(container, application) {
        var PartialslotAbleController = container.lookupFactory('controller:partialslot-able');

        var UserprofilestatusmenuMixin = container.lookupFactory('mixin:userprofilestatusmenu');
        var RequirejsmocksmanagerMixin = container.lookupFactory('mixin:requirejsmocksmanager');
        var ScreentoolstatusmenuMixin = container.lookupFactory('mixin:screentoolstatusmenu');
        var DocumentationMixin = container.lookupFactory('mixin:documentation');
        var SchemamanagerMixin = container.lookupFactory('mixin:schemamanager');
        var ConsolemanagerMixin = container.lookupFactory('mixin:consolemanager');
        var LoadingindicatorMixin = container.lookupFactory('mixin:loadingindicator');
        var PromisemanagerMixin = container.lookupFactory('mixin:promisemanager');
        var NotificationsMixin = container.lookupFactory('mixin:notifications');

        var formsUtils = container.lookupFactory('utility:forms');
        var dataUtils = container.lookupFactory('utility:data');
        var debugUtils = container.lookupFactory('utility:debug');
        var hashUtils = container.lookupFactory('utility:hash');
        var notificationUtils = container.lookupFactory('utility:notification');

        var controller;

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc,
            isEmpty = Ember.isEmpty;


        var indexController = Ember.Controller.extend(Ember.Evented, {});

        /**
         * @class ApplicationController
         * @extends PartialslotAbleController
         * @constructor
         * @description Main application controller
         */
        var ApplicationControllerDict = {

            needs: ['login', 'recordinfopopup'],

            partials: {
                statusbar: []
            },

            /**
             * @property currentViewModel
             * @type DS.Model
             * @description The model of the userview currently opened
             */
            currentViewModel: undefined,

            /**
             * @property currentViewId
             * @type DS.Model
             * @description The id of the userview currently opened
             */
            currentViewId: undefined,

            /**
             * @property editMode
             * @type {boolean}
             * @description Whether the current view is editable or not
             */
            editMode: false,

            /**
             * @property runtimeConfiguration
             * @type object
             * @description the canopsis frontend configuration object
             */
            runtimeConfiguration: canopsisConfiguration,

            /**
             * @property debug
             * @type boolean
             */
            debug: Ember.computed.alias('runtimeConfiguration.DEBUG'),

            /**
             * @property defaultView
             * @type string
             */
            defaultView: function() {
                var userDefaultView = get(this, 'controllers.login.record.defaultview');
                var roleDefaultView = get(this, 'controllers.login.role.defaultview');
                var frontendDefaultView = get(this, 'frontendConfig.defaultview');

                if(!isEmpty(userDefaultView)) {
                    return userDefaultView;
                } else if(!isEmpty(roleDefaultView)) {
                    return roleDefaultView;
                } else {
                    return frontendDefaultView;
                }

            }.property('frontendConfig.defaultview', 'controllers.login.defaultview'),

            /**
             * @property enginesviews
             * @type Array
             * @default
             * @description used to feed the left menu "Engines"
             */
            enginesviews: [
                {label: __('Context'), value: 'view.context'},
                {label: __('Selectors'), value: 'view.selectors'},
                {label: __('Event Filter'), value: 'view.filters'},
                {label: __('Performance Data'), value: 'view.series'},
                {label: __('Scheduled Jobs'), value: 'view.jobs'},
                {label: __('Notifications'), value: 'view.notifications'},
                {label: __('Link list'), value: 'view.linklist'}
            ],

            /**
             * @method init
             */
            init: function() {
                console.group('app init');
                notificationUtils.setController(this);

                //Set page title
                var title = get(canopsisConfiguration, 'TITLE');
                if (!isNone(title)) {
                    $('title').html(title);
                }

                console.groupEnd();
                this.refreshPartialsList();
                this._super.apply(this, arguments);

                dataUtils.setLoggedUserController(get(this, 'controllers.login'));

                //close the init group
                console.groupEnd();
                console.tags.remove('init');
            },

            /**
             * @method editAuthConfig
             * @param {String} authType the type of the auth system to edit
             */
            editAuthConfig: function(authType) {
                console.log('edit ' + authType);

                var conf = get(this, authType);

                if(isNone(conf)) {
                    this.send(
                        'prompt',
                        'Unable to load ' + authType + ' configuration. No edition possible.'
                    );
                } else {

                    set(conf, 'crecord_type', authType);

                    var editForm = formsUtils.showNew(
                        'modelform',
                        conf,
                        { title: __('Edit ' + authType + ' configuration') }
                    );

                    editForm.submit.done(function() {
                        conf.save();
                    });
                }
            },

            /**
             * @method didSaveView
             * @description method triggered when an userview is saved
             */
            didSaveView: function(userview) {
                this.transitionToRoute("/userview/" + get(userview, 'id'));
            },

            actions: {
                /**
                 * @event inspect
                 * @param {object} object
                 * Show debug info in console and put widget var in window.$E
                 */
                inspect: function (object) {
                    debugUtils.inspectObject(object);
                },

                /**
                 * @event editConfig
                 *
                 * Shows a form to edit the frontend configuration
                 */
                editConfig: function() {
                    var frontendConfig = get(this, 'frontendConfig');
                    console.log('editConfig', formsUtils, frontendConfig);

                    var editForm = formsUtils.showNew('modelform', frontendConfig, { title: __("Edit settings"), inspectedItemType: "frontend" });
                    editForm.submit.done(function() {
                        frontendConfig.save();
                    });
                },

                /**
                 * @event editAuthConfiguration
                 * @param {String} authType
                 */
                editAuthConfiguration: function(authType) {
                    console.log('edit ' + authType);
                    this.editAuthConfig(authType);
                },

                /**
                 * @event prompt
                 * @description Shows a popup with a message
                 * @param {String} message
                 */
                prompt:function (message) {
                    formsUtils.showNew('confirmform', {}, {
                        title: __(message)
                    });
                },

                /**
                 * @event promptReloadApplication
                 * @description Shows a popup with a message
                 * @param {String} title
                 * @param {String} location the location to redirect the user when he accepts to get redirected. Defaults to "/index.html"
                 */
                promptReloadApplication: function(title, location) {
                    setTimeout(function (){
                        var recordWizard = formsUtils.showNew('confirmform', {}, {
                            title: __(title)
                        });

                        recordWizard.submit.then(function(form) {
                            void(form);

                            if (isNone(location)) {
                                window.location = '/index.html';
                            } else {
                                window.location = location;
                            }
                        });
                    }, 1500);
                },

                /**
                 * @event editEventSettings
                 * @description Shows a Modelform with event settings
                 */
                editEventSettings: function () {
                    formsUtils.editSchemaRecord('statusmanagement', get(this, "container"));
                },

                /**
                 * @event editDataclean
                 * @description Shows a Modelform to edit data cleaner options
                 */
                editDataclean: function () {
                    formsUtils.editSchemaRecord('datacleaner', get(this, "container"));
                },


                 /**
                 * @event editEnabledPlugins
                 * @description Shows a form to edit the canopsis UI enabled plugins
                 */
                editEnabledPlugins: function() {
                    formsUtils.editSchemaRecord('enabledmodules', get(this, "container"));
                },

                /**
                 * @event editTicketJob
                 * @description Shows a form to edit the tickets configuration
                 */
                editTicketJob: function() {
                    console.group('editTicketJob');

                    var ticketConfig = get(this, 'ticketConfig');
                    set(ticketConfig, 'crecord_type', 'ticket');

                    console.log('ticketConfig:', ticketConfig);

                    var job = get(ticketConfig, 'job');
                    var params = get(job, 'params');

                    console.log('job:', job, params);

                    if(params) {
                        console.log('param subject', get(params, 'subject'));
                    }

                    var editForm = formsUtils.showNew('jobform', job, {
                        scheduled: false
                    });

                    editForm.submit.done(function(form) {
                        var job = get(form, 'formContext');
                        var params = get(job, 'params');

                        set(ticketConfig, 'job', job);
                        set(ticketConfig, 'job.params', params);

                        ticketConfig.save();
                    });

                    console.groupEnd();
                },

                /**
                 * @event addNewView
                 * @description Shows a form to create a new userview
                 */
                addNewView: function () {
                    var type = 'userview';
                    var applicationController = this;
                    console.log('add', type);

                    var containerwidgetId = hashUtils.generateId('container');
                    var viewId = hashUtils.generateId('userview');

                    dataUtils.getStore().createRecord('widgetcontainer', {
                        xtype: 'widgetcontainer',
                        mixins : [{
                            name: 'verticallayout'
                        }],
                        id: containerwidgetId
                    });

                    var userId = get(this, 'controllers.login.record._id');

                    var userview = dataUtils.getStore().push(type, {
                        id: viewId,
                        crecord_type: 'view',
                        author: userId,
                        containerwidget: containerwidgetId,
                        containerwidgetType: 'widgetcontainer'
                    });

                    var formattedViewId = viewId.replace('.', '_');

                    console.log('temp record', userview);

                    var recordWizard = formsUtils.showNew('modelform', userview, { title: __("Add ") + type });

                    recordWizard.submit.done(function() {
                        userview.save().then(function() {
                            applicationController.didSaveView(userview);
                        });
                    });
                },

                /**
                 * @event importView
                 * @description Shows a file upload window, and then import the selected view
                 */
                importView: function () {
                    var applicationController = this;

                    dataUtils.uploadTextFilePopup(function(name, type, size, content) {
                        content = JSON.parse(content);

                        var res = applicationController.get('store').pushPayload('userview', {
                            userview: content
                        });

                        console.error('done', res, applicationController.get('store').getById('userview', content.id));
                        applicationController.get('store').getById('userview', content.id).save();
                    });
                },

                /**
                 * @event openTab
                 * @description Change frontend route to a new url
                 * @param {string} url the new url to go to
                 */
                openTab: function(url) {
                    this.transitionToRoute(url);
                },

                /**
                 * @event logout
                 * @description close the user session and go back to the login screen
                 */
                logout: function() {
                    get(this, 'controllers.login').setProperties({
                        'authkey': null,
                        'errors': []
                    });

                    window.location = '/logout';
                }
            }

        };

        if(canopsisConfiguration.DEBUG) {
            controller = PartialslotAbleController.extend(
                UserprofilestatusmenuMixin,
                SchemamanagerMixin,
                PromisemanagerMixin,
                ConsolemanagerMixin,
                NotificationsMixin,
                LoadingindicatorMixin,
                RequirejsmocksmanagerMixin,
                ScreentoolstatusmenuMixin,
                DocumentationMixin,
                ApplicationControllerDict);
        } else {
            controller = PartialslotAbleController.extend(
                UserprofilestatusmenuMixin,
                NotificationsMixin,
                LoadingindicatorMixin,
                DocumentationMixin,
                ApplicationControllerDict);
        }


        application.register('controller:application', controller);
        application.register('controller:index', indexController);
    }
});
