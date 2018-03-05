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
    name:'RecordinfopopupMixin',
    after: 'MixinFactory',
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');

        var get = Ember.get;

        /**
         * @mixin recordinfopopup
         */
        var mixin = Mixin('recordinfopopup', {
            actions: {
                /**
                 * @method actions-sendDisplayRecord
                 */
                sendDisplayRecord: function (record) {
                    //This method is not ugly TODO refactor, it would be better if event bubble until application directly
                    // but at the moment, event doen t bubble properly
                    console.log('sendDisplayRecord action called with params', record);

                    var template = get(this, 'mixinOptions.recordinfopopup.popup_template');
                    if (Ember.isNone(template)) {
                        template = '';
                    }

                    console.log('Template is ', template);

                    var recordinfopopupController = get(this, 'controllers.recordinfopopup');

                    recordinfopopupController.send('show', record, template);
                }
            },

            /**
             * @method rendererFor
             * @argument attribute
             * @return {string} the renderer name
             */
            rendererFor: function(attribute) {
                var clickableColumn = get(this, 'mixinOptions.recordinfopopup.clickable_column');
                console.log('recordinfopopup rendererFor', attribute, clickableColumn);
                if(attribute.field === clickableColumn){
                    return 'renderer-recordinfopopup';
                }
                return this._super(attribute);
            }
        });

        application.register('mixin:recordinfopopup', mixin);
    }
});
