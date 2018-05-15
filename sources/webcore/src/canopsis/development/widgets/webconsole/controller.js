/*
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
*/

define([
    'app/lib/factories/widget'
], function(WidgetFactory) {
    var get = Ember.get,
        set = Ember.set;

    var consoleRenderers = {
        widget: function (widget) {
            return get(widget, 'xtype') + '(' + get(widget, 'title') + ')';
        }
    };

    guessObjectType = function(element) {
        if(element && element.abstractType) {
            return element.abstractType || 'object';
        }
    };

    prettyLoggedType = function(element) {
        var elType = typeof element;
        if(elType === 'string') {
            return '"' + element + '"';
        } else if(elType === 'undefined') {
            return '<i style="color:red">' + element + '</i>';
        } else if(elType === "object") {
            var objectAbstractType = guessObjectType(element);
            var consoleRenderer = consoleRenderers[objectAbstractType];

            var text;
            if(consoleRenderer) {
                text = consoleRenderer(element);
            } else {
                text = objectAbstractType;
            }

            return '<i class="block ' + objectAbstractType + '">'+ text +'</i>';
        } else return '<i style="color:red">(unknown:'+ element +')</i>';
    };

    var WebconsoleViewMixin = Ember.Mixin.create({
        classNames: ['webconsole'],

        didInsertElement: function() {
            this._super.apply(this, arguments);
            this.logger = this.$('.logger');
            var textareaEditor = this.$('textarea').get(0);

            this.controller.loggerDiv = this.$('.logger');
            console.backends.add("webconsole", get(this,'controller'));
        },
        willDestroyElement: function () {
            console.backends.remove("webconsole", get(this,'controller'));
            this._super.apply(this, arguments);
        }
    });

    var widget = WidgetFactory('webconsole', {
        messages: [],
        viewMixins: [
            WebconsoleViewMixin
        ],

        init: function() {
            this._super.apply(this, arguments);
            this.buffer = "";
        },

        send: function(function_name, args) {
            if(function_name !== 'info') {
                if(this.loggerDiv) {
                    var logline = '<tr class="logline"><td class="logauthor">' + args[0] + '</td>';
                    var buffer = '' + args[1];
                    for (var i = 2, l = args.length; i < l; i++) {
                        buffer += prettyLoggedType(args[i]) + ' ';
                    }
                    logline += '<td class="logline">' + buffer + '</td></tr>';
                    this.loggerDiv.append($(logline));
                }
            }
        }
    });

    return widget;
});
