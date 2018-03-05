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

 require.config({
    paths: {
        'components/component-editor': 'canopsis/core/src/components/editor/template',
        'components/component-renderer': 'canopsis/core/src/components/renderer/template',
        'confirmform': 'canopsis/core/src/forms/confirmform/confirmform',
        'modelform': 'canopsis/core/src/forms/modelform/modelform',
        'widgetform': 'canopsis/core/src/forms/widgetform/widgetform',

    }
});

 define([
    'canopsis/core/src/components/editor/component',
    'ehbs!components/component-editor',
    'canopsis/core/src/components/renderer/component',
    'ehbs!components/component-renderer',
    'canopsis/core/src/controller/application',
    'canopsis/core/src/controller/form',
    'canopsis/core/src/controller/formwrapper',
    'canopsis/core/src/controller/login',
    'canopsis/core/src/controller/partialslotable',
    'canopsis/core/src/controller/perfdata',
    'canopsis/core/src/controller/recordinfopopup',
    'canopsis/core/src/controller/serie',
    'canopsis/core/src/controller/userview',
    'canopsis/core/src/controller/widget',
    'ehbs!confirmform',
    'canopsis/core/src/forms/confirmform/controller',
    'canopsis/core/src/forms/modelform/controller',
    'ehbs!modelform',
    'canopsis/core/src/forms/widgetform/controller',
    'ehbs!widgetform',
    'canopsis/core/src/lib/abstractclassregistry',
    'canopsis/core/src/lib/actionsregistry',
    'canopsis/core/src/lib/contextsregistry',
    'canopsis/core/src/lib/editorregistry',
    'canopsis/core/src/lib/factories/form',
    'canopsis/core/src/lib/factories/mixin',
    'canopsis/core/src/lib/factories/widget',
    'canopsis/core/src/lib/formsregistry',
    'canopsis/core/src/lib/helpers/eachindexed',
    'canopsis/core/src/lib/helpers/formview',
    'canopsis/core/src/lib/helpers/getfield',
    'canopsis/core/src/lib/helpers/i18n',
    'canopsis/core/src/lib/helpers/log',
    'canopsis/core/src/lib/helpers/partialslot',
    'canopsis/core/src/lib/helpers/renderwidget',
    'canopsis/core/src/lib/helpers/validationtextarea',
    'canopsis/core/src/lib/helpers/validationtextfield',
    'canopsis/core/src/lib/helpers/widgetslot',
    'canopsis/core/src/lib/indexesregistry',
    'canopsis/core/src/lib/inflections',
    'canopsis/core/src/lib/mixinsregistry',
    'canopsis/core/src/lib/promisesmanager',
    'canopsis/core/src/lib/rendererregistry',
    'canopsis/core/src/lib/requirejsmocksmanager',
    'canopsis/core/src/lib/searchmethodsregistry',
    'canopsis/core/src/lib/templateregistry',
    'canopsis/core/src/lib/utilityclass',
    'canopsis/core/src/lib/utils/actions',
    'canopsis/core/src/lib/utils/data',
    'canopsis/core/src/lib/utils/dates',
    'canopsis/core/src/lib/utils/debug',
    'canopsis/core/src/lib/utils/dom',
    'canopsis/core/src/lib/utils/drag',
    'canopsis/core/src/lib/utils/event',
    'canopsis/core/src/lib/utils/filterObject',
    'canopsis/core/src/lib/utils/forms',
    'canopsis/core/src/lib/utils/hash',
    'canopsis/core/src/lib/utils/indexes',
    'canopsis/core/src/lib/utils/modelsolve',
    'canopsis/core/src/lib/utils/notification',
    'canopsis/core/src/lib/utils/routes',
    'canopsis/core/src/lib/utils/slug',
    'canopsis/core/src/lib/utils/timewindow',
    'canopsis/core/src/lib/utils/values',
    'canopsis/core/src/lib/utils/widgetSelectors',
    'canopsis/core/src/lib/utils/widgets',
    'canopsis/core/src/lib/widgetsregistry',
    'canopsis/core/src/lib/wrappers/console',
    'canopsis/core/src/mixins/consolemanager',
    'canopsis/core/src/mixins/criticitylevels',
    'canopsis/core/src/mixins/documentation',
    'canopsis/core/src/mixins/embeddedrecordserializer',
    'canopsis/core/src/mixins/hashserializer',
    'canopsis/core/src/mixins/inspectablearray',
    'canopsis/core/src/mixins/inspectableitem',
    'canopsis/core/src/mixins/loadingindicator',
    'canopsis/core/src/mixins/metaserializer',
    'canopsis/core/src/mixins/notifications',
    'canopsis/core/src/mixins/promisemanager',
    'canopsis/core/src/mixins/requirejsmocksmanager',
    'canopsis/core/src/mixins/schemamanager',
    'canopsis/core/src/mixins/screentoolstatusmenu',
    'canopsis/core/src/mixins/userprofilestatusmenu',
    'canopsis/core/src/mixins/validation',
    'canopsis/core/src/mixins/validationfield',
    'canopsis/core/src/routes/application',
    'canopsis/core/src/routes/authenticated',
    'canopsis/core/src/routes/index',
    'canopsis/core/src/routes/userview',
    'canopsis/core/src/serializers/application',
    'canopsis/core/src/serializers/loggedaccount',
    'canopsis/core/src/serializers/userview',
    'canopsis/core/src/serializers/widget',
    'canopsis/core/src/serializers/widgetwrapper',
    'canopsis/core/src/validators/mail/validator',
    'canopsis/core/src/validators/maxItems/validator',
    'canopsis/core/src/validators/minItems/validator',
    'canopsis/core/src/validators/number/validator',
    'canopsis/core/src/validators/required/validator',
    'canopsis/core/src/validators/rights/validator',
    'canopsis/core/src/validators/validate/validator',
    'canopsis/core/src/view/application',
    'canopsis/core/src/view/editor',
    'canopsis/core/src/view/formwrapper',
    'canopsis/core/src/view/listline',
    'canopsis/core/src/view/mixineditdropdown',
    'canopsis/core/src/view/partialslot',
    'canopsis/core/src/view/recordinfopopup',
    'canopsis/core/src/view/tabledraggableth',
    'canopsis/core/src/view/userview',
    'canopsis/core/src/view/validationtextarea',
    'canopsis/core/src/view/validationtextfield',
    'canopsis/core/src/view/widget',
    'canopsis/core/src/view/widgetslot',
    'canopsis/core/requirejs-modules/externals.conf',
    'canopsis/core/requirejs-modules/i18n'
], function (templates) {
    templates = $(templates).filter('script');
for (var i = 0, l = templates.length; i < l; i++) {
var tpl = $(templates[i]);
Ember.TEMPLATES[tpl.attr('data-template-name')] = Ember.Handlebars.compile(tpl.text());
};
});

