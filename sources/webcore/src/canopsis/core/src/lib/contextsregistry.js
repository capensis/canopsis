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
    name: 'ContextRegistry',
    after: 'AbstractClassRegistry',
    initialize: function(container, application) {

        var Abstractclassregistry = container.lookupFactory('registry:abstractclass');

        var contexts = [{
                name: 'component',
                modelName: 'ctxcomponent'
            },{
                name: 'resource',
                modelName: 'ctxresource'
            },{
                name: 'selector',
                modelName: 'ctxselector'
            },{
                name: 'topology',
                modelName: 'ctxtopology'
            }];

        var registry = Abstractclassregistry.create({
            name: 'contexts',

            all: [],
            byClass: {},
            tableColumns: [{
                title: 'name',
                name: 'name'
            },{
                title: 'modelName',
                name: 'modelName'
            }]
        });

        for (var i = 0, l = contexts.length; i < l; i++) {
            registry.all.pushObject({
                name: contexts[i].name,
                modelName: contexts[i].modelName
            });
        }

        application.register('registry:context', registry);
    }
});
