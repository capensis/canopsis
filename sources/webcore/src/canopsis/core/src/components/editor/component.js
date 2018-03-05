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
    name: 'component-editor',
    after: 'DebugUtils',
    initialize: function(container, application) {
        var debugUtils = container.lookupFactory('utility:debug');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @component editor
         *
         * @description Component displaying the correct editor for an attribute.
         * It is possible to specify the editor type to use. If not specified, it will try to get the correct type on its own.
         *
         * @example {{component-editor attr=item editorType=itemEditorType}}
         */
        var component = Ember.Component.extend({
            /**
             * @property tagName
             * @type string
             * @default
             */
            tagName: 'span',

            /**
             * @property canopsisConfiguration
             * @type object
             * @description the canopsis frontend configuration object
             */
            canopsisConfiguration: canopsisConfiguration,

            /**
             * @property debug
             * @description whether the UI is in debug mode or not
             * @type boolean
             * @default Ember.computed.alias('canopsisConfiguration.DEBUG')
             */
            debug: Ember.computed.alias('canopsisConfiguration.DEBUG'),

            /**
             * @property content
             * @type object
             * @description the attribute which information is displayed and that can be edited with the editor.
             */
            content: undefined,

            /**
             * @method init
             * @description binds the "onValueChange" observer
             */
            init: function() {
                console.log('init editor compo');

                //FIXME auto-detect if we need standalone mode or not, stop using a variable, for a better comprehension
                if(get(this, 'mode') === 'standalone') {
                    set(this, 'attr', { value: undefined });
                    Ember.addObserver(this, 'attr.value', this, this.onValueChange);
                }

                this._super();
            },

            //TODO check if still used
            /**
             * @method onValueChange
             *
             * @description Triggered automatically when "attr.value" change
             */
            onValueChange: function () {
                set(this, 'value', get(this, 'attr.value'));
            },

            actions: {
                /**
                 * @method actions_inspect
                 * @description inspects the object in the console (see debugUtils for more info)
                 */
                inspect: function() {
                    debugUtils.inspectObject(get(this, 'content'));
                }
            },

            /**
             * @property description
             * @description Computed property, binded to the description of the attribute
             * @type string
             */
            description: function () {
                return get(this, 'content.model.options.description');
            }.property(),

            /**
             * @property editorType
             * @description finds the correct template name for the editor
             * @return {string} the template name found, or "editor-defaultpropertyeditor" if nothing is suitable for the attribute
             * @type string
             */
            editorType: function() {
                console.group('editorType');

                var overrides = get(this, 'editorOverrides');

                var type = get(this, 'content.model.type');
                var role = get(this, 'content.model.options.role');
                var field = get(this, 'content.field');


                console.log('content:', get(this, 'content'));
                console.log('type:', get(this, 'content.field'));
                console.log('type:', type);
                console.log('role:', role);

                if (field == 'tabs') {
                    return 'editor-nestedarray'                
                }

                var editorName;

                if(!isNone(overrides) && !isNone(field) && get(overrides, field)) {
                    editorName = 'editor-' + get(overrides, field);
                } else {
                    if (role) {
                        if(!isNone(overrides) && get(overrides, role)) {
                            editorName = 'editor-' + get(overrides, role);
                        } else {
                            editorName = 'editor-' + role;
                        }
                    } else {
                        if(!isNone(overrides) && get(overrides, type)) {
                            editorName = 'editor-' + get(overrides, type);
                        } else {
                            editorName = 'editor-' + type;
                        }
                    }
                }

                if (Ember.TEMPLATES[editorName] === undefined) {
                    editorName = 'editor-defaultpropertyeditor';
                }

                console.groupEnd();
                
                return editorName;
            }.property('content.type', 'content.role'),

            /**
             * @property attr
             * @default Ember.computed.alias('content')
             */
            attr: Ember.computed.alias('content')
        });

        application.register('component:component-editor', component);
    }
});
