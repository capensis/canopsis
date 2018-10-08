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
        'components/component-pbehaviorslist': 'canopsis/brick-listalarm/src/components/pbehaviorslist/template',
        'components/component-popupinfo': 'canopsis/brick-listalarm/src/components/popupinfo/template',
        'components/component-radio': 'canopsis/brick-listalarm/src/components/radio/template',
        'components/component-rendererack': 'canopsis/brick-listalarm/src/components/rendererack/template',
        'components/component-rendererextradetails': 'canopsis/brick-listalarm/src/components/rendererextradetails/template',
        'components/component-rendererlinks': 'canopsis/brick-listalarm/src/components/rendererlinks/template',
        'components/component-rendererpbehaviors': 'canopsis/brick-listalarm/src/components/rendererpbehaviors/template',
        'components/component-rendererstate': 'canopsis/brick-listalarm/src/components/rendererstate/template',
        'components/component-rendererstatetimestamp': 'canopsis/brick-listalarm/src/components/rendererstatetimestamp/template',
        'components/component-rendererstatus': 'canopsis/brick-listalarm/src/components/rendererstatus/template',
        'components/component-rendererstatusval': 'canopsis/brick-listalarm/src/components/rendererstatusval/template',
        'components/component-search': 'canopsis/brick-listalarm/src/components/search/template',
        'components/component-selectionactions': 'canopsis/brick-listalarm/src/components/selectionactions/template',
        'components/component-selectioncheckbox': 'canopsis/brick-listalarm/src/components/selectioncheckbox/template',
        'editor-pair': 'canopsis/brick-listalarm/src/editors/editor-pair',
        'editor-radio': 'canopsis/brick-listalarm/src/editors/editor-radio',
        'snoozeform': 'canopsis/brick-listalarm/src/forms/snooze/snoozeform',
        '_infos_hostgroups_value': 'canopsis/brick-listalarm/src/partials/_infos_hostgroups_value',
        '_v_ack': 'canopsis/brick-listalarm/src/partials/_v_ack',
        '_v_creation_date': 'canopsis/brick-listalarm/src/partials/_v_creation_date',
        '_v_current_state_duration': 'canopsis/brick-listalarm/src/partials/_v_current_state_duration',
        '_v_duration': 'canopsis/brick-listalarm/src/partials/_v_duration',
        '_v_extra_details': 'canopsis/brick-listalarm/src/partials/_v_extra_details',
        '_v_last_event_date': 'canopsis/brick-listalarm/src/partials/_v_last_event_date',
        '_v_last_update_date': 'canopsis/brick-listalarm/src/partials/_v_last_update_date',
        '_v_links': 'canopsis/brick-listalarm/src/partials/_v_links',
        '_v_pbehaviors': 'canopsis/brick-listalarm/src/partials/_v_pbehaviors',
        '_v_resolved': 'canopsis/brick-listalarm/src/partials/_v_resolved',
        '_v_state': 'canopsis/brick-listalarm/src/partials/_v_state',
        '_v_state_m': 'canopsis/brick-listalarm/src/partials/_v_state_m',
        '_v_state_t': 'canopsis/brick-listalarm/src/partials/_v_state_t',
        '_v_status': 'canopsis/brick-listalarm/src/partials/_v_status',
        '_v_status_val': 'canopsis/brick-listalarm/src/partials/_v_status_val',
        'renderer-listalarm_ellipsis': 'canopsis/brick-listalarm/src/renderers/renderer-listalarm_ellipsis',
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
    'canopsis/brick-listalarm/src/components/pbehaviorslist/component',
    'ehbs!components/component-pbehaviorslist',
    'canopsis/brick-listalarm/src/components/popupinfo/component',
    'ehbs!components/component-popupinfo',
    'canopsis/brick-listalarm/src/components/radio/component',
    'ehbs!components/component-radio',
    'canopsis/brick-listalarm/src/components/rendererack/component',
    'ehbs!components/component-rendererack',
    'canopsis/brick-listalarm/src/components/rendererextradetails/component',
    'ehbs!components/component-rendererextradetails',
    'canopsis/brick-listalarm/src/components/rendererlinks/component',
    'ehbs!components/component-rendererlinks',
    'canopsis/brick-listalarm/src/components/rendererpbehaviors/component',
    'ehbs!components/component-rendererpbehaviors',
    'canopsis/brick-listalarm/src/components/rendererstate/component',
    'ehbs!components/component-rendererstate',
    'canopsis/brick-listalarm/src/components/rendererstatetimestamp/component',
    'ehbs!components/component-rendererstatetimestamp',
    'canopsis/brick-listalarm/src/components/rendererstatus/component',
    'ehbs!components/component-rendererstatus',
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
    'canopsis/brick-listalarm/src/forms/snooze/controller',
    'ehbs!snoozeform',
    'canopsis/brick-listalarm/src/helpers/absoluteTimeSince',
    'canopsis/brick-listalarm/src/helpers/dateFormat',
    'canopsis/brick-listalarm/src/helpers/durationFromTimestamp',
    'canopsis/brick-listalarm/src/helpers/listalarm_ellipsis',
    'canopsis/brick-listalarm/src/mixins/rinfopop',
    'ehbs!_infos_hostgroups_value',
    'ehbs!_v_ack',
    'ehbs!_v_creation_date',
    'ehbs!_v_current_state_duration',
    'ehbs!_v_duration',
    'ehbs!_v_extra_details',
    'ehbs!_v_last_event_date',
    'ehbs!_v_last_update_date',
    'ehbs!_v_links',
    'ehbs!_v_pbehaviors',
    'ehbs!_v_resolved',
    'ehbs!_v_state',
    'ehbs!_v_state_m',
    'ehbs!_v_state_t',
    'ehbs!_v_status',
    'ehbs!_v_status_val',
    'ehbs!renderer-listalarm_ellipsis',
    'canopsis/brick-listalarm/src/serializers/alertexpression',
    'canopsis/brick-listalarm/src/serializers/alerts',
    'link!canopsis/brick-listalarm/src/style.css',
    'canopsis/brick-listalarm/src/widgets/listalarm/controller',
    'ehbs!listalarm',
    'canopsis/brick-listalarm/requirejs-modules/externals.conf',
    'canopsis/brick-listalarm/requirejs-modules/i18n'
], function () {
    
});
