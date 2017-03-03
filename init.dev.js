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
        'components/component-alarmactions': 'canopsis/brick-listalarm/src/components/alarmactions/template',
        'components/component-alarmraw': 'canopsis/brick-listalarm/src/components/alarmraw/template',
        'components/component-alarmstate': 'canopsis/brick-listalarm/src/components/alarmstate/template',
        'components/component-alarmtable': 'canopsis/brick-listalarm/src/components/alarmtable/template',
        'components/component-alarmtd': 'canopsis/brick-listalarm/src/components/alarmtd/template',
        'components/component-columntemplate': 'canopsis/brick-listalarm/src/components/columntemplate/template',
        'components/component-customtimeline': 'canopsis/brick-listalarm/src/components/customtimeline/template',
        'components/component-popupinfo': 'canopsis/brick-listalarm/src/components/popupinfo/template',
        'components/component-rendererack': 'canopsis/brick-listalarm/src/components/rendererack/template',
        'components/component-rendererpbehaviors': 'canopsis/brick-listalarm/src/components/rendererpbehaviors/template',
        'components/component-rendererstate': 'canopsis/brick-listalarm/src/components/rendererstate/template',
        'components/component-rendererstatetimestamp': 'canopsis/brick-listalarm/src/components/rendererstatetimestamp/template',
        'components/component-rendererstatusval': 'canopsis/brick-listalarm/src/components/rendererstatusval/template',
        'components/component-search': 'canopsis/brick-listalarm/src/components/search/template',
        'components/component-selectionactions': 'canopsis/brick-listalarm/src/components/selectionactions/template',
        'components/component-selectioncheckbox': 'canopsis/brick-listalarm/src/components/selectioncheckbox/template',
        'editor-pair': 'canopsis/brick-listalarm/src/editors/editor-pair',
        'editor-radio': 'canopsis/brick-listalarm/src/editors/editor-radio',
        '_v_ack': 'canopsis/brick-listalarm/src/partials/_v_ack',
        '_v_pbehaviors': 'canopsis/brick-listalarm/src/partials/_v_pbehaviors',
        '_v_state_t': 'canopsis/brick-listalarm/src/partials/_v_state_t',
        '_v_state_val': 'canopsis/brick-listalarm/src/partials/_v_state_val',
        '_v_status_val': 'canopsis/brick-listalarm/src/partials/_v_status_val',
        'listalarm': 'canopsis/brick-listalarm/src/widgets/listalarm/listalarm',

    }
});

 define([
    'canopsis/brick-listalarm/src/adapters/alertexpression',
    'canopsis/brick-listalarm/src/adapters/alerts',
    'canopsis/brick-listalarm/src/components/alarmactions/component',
    'ehbs!components/component-alarmactions',
    'canopsis/brick-listalarm/src/components/alarmraw/component',
    'ehbs!components/component-alarmraw',
    'canopsis/brick-listalarm/src/components/alarmstate/component',
    'ehbs!components/component-alarmstate',
    'canopsis/brick-listalarm/src/components/alarmtable/component',
    'ehbs!components/component-alarmtable',
    'canopsis/brick-listalarm/src/components/alarmtd/component',
    'ehbs!components/component-alarmtd',
    'canopsis/brick-listalarm/src/components/columntemplate/component',
    'ehbs!components/component-columntemplate',
    'canopsis/brick-listalarm/src/components/customtimeline/component',
    'ehbs!components/component-customtimeline',
    'canopsis/brick-listalarm/src/components/popupinfo/component',
    'ehbs!components/component-popupinfo',
    'canopsis/brick-listalarm/src/components/rendererack/component',
    'ehbs!components/component-rendererack',
    'canopsis/brick-listalarm/src/components/rendererpbehaviors/component',
    'ehbs!components/component-rendererpbehaviors',
    'canopsis/brick-listalarm/src/components/rendererstate/component',
    'ehbs!components/component-rendererstate',
    'canopsis/brick-listalarm/src/components/rendererstatetimestamp/component',
    'ehbs!components/component-rendererstatetimestamp',
    'canopsis/brick-listalarm/src/components/rendererstatusval/component',
    'ehbs!components/component-rendererstatusval',
    'canopsis/brick-listalarm/src/components/search/component',
    'ehbs!components/component-search',
    'canopsis/brick-listalarm/src/components/selectionactions/component',
    'ehbs!components/component-selectionactions',
    'canopsis/brick-listalarm/src/components/selectioncheckbox/component',
    'ehbs!components/component-selectioncheckbox',
    'ehbs!editor-pair',
    'ehbs!editor-radio',
    'canopsis/brick-listalarm/src/mixins/customsendevent',
    'canopsis/brick-listalarm/src/mixins/rinfopop',
    'ehbs!_v_ack',
    'ehbs!_v_pbehaviors',
    'ehbs!_v_state_t',
    'ehbs!_v_state_val',
    'ehbs!_v_status_val',
    'canopsis/brick-listalarm/src/serializers/alertexpression',
    'canopsis/brick-listalarm/src/serializers/alerts',
    'link!canopsis/brick-listalarm/src/style.css',
    'canopsis/brick-listalarm/src/widgets/listalarm/controller',
    'ehbs!listalarm',
    'canopsis/brick-listalarm/requirejs-modules/externals.conf'
], function () {
    
});

