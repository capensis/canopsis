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
        'components/component-ack': 'canopsis/monitoring/src/components/ack/template',
        'components/component-cfiltereditor': 'canopsis/monitoring/src/components/cfiltereditor/template',
        'components/component-eventselector': 'canopsis/monitoring/src/components/eventselector/template',
        'components/component-stateeditor': 'canopsis/monitoring/src/components/stateeditor/template',
        'components/component-statemapping': 'canopsis/monitoring/src/components/statemapping/template',
        'editor-cfilter': 'canopsis/monitoring/src/editors/editor-cfilter',
        'editor-cfilterwithproperties': 'canopsis/monitoring/src/editors/editor-cfilterwithproperties',
        'editor-cmetric': 'canopsis/monitoring/src/editors/editor-cmetric',
        'editor-eventselector': 'canopsis/monitoring/src/editors/editor-eventselector',
        'editor-metricselector': 'canopsis/monitoring/src/editors/editor-metricselector',
        'renderer-ack': 'canopsis/monitoring/src/renderers/renderer-ack',
        'renderer-cfilter': 'canopsis/monitoring/src/renderers/renderer-cfilter',
        'renderer-cfilterwithproperties': 'canopsis/monitoring/src/renderers/renderer-cfilterwithproperties',
        'renderer-crecord-type': 'canopsis/monitoring/src/renderers/renderer-crecord-type',
        'renderer-criticity': 'canopsis/monitoring/src/renderers/renderer-criticity',
        'renderer-eventselector': 'canopsis/monitoring/src/renderers/renderer-eventselector',
        'renderer-eventtimestamp': 'canopsis/monitoring/src/renderers/renderer-eventtimestamp',
        'renderer-eventtype': 'canopsis/monitoring/src/renderers/renderer-eventtype',
        'renderer-state': 'canopsis/monitoring/src/renderers/renderer-state',
        'renderer-stateConnector': 'canopsis/monitoring/src/renderers/renderer-stateConnector',
        'renderer-status': 'canopsis/monitoring/src/renderers/renderer-status',
        'actionbutton-editurlfield': 'canopsis/monitoring/src/templates/actionbutton-editurlfield',
        'weather': 'canopsis/monitoring/src/widgets/weather/weather',

    }
});

 define([
    'canopsis/monitoring/src/components/ack/component',
    'ehbs!components/component-ack',
    'canopsis/monitoring/src/components/cfiltereditor/component',
    'ehbs!components/component-cfiltereditor',
    'canopsis/monitoring/src/components/eventselector/component',
    'ehbs!components/component-eventselector',
    'canopsis/monitoring/src/components/stateeditor/component',
    'ehbs!components/component-stateeditor',
    'canopsis/monitoring/src/components/statemapping/component',
    'ehbs!components/component-statemapping',
    'ehbs!editor-cfilter',
    'ehbs!editor-cfilterwithproperties',
    'ehbs!editor-cmetric',
    'ehbs!editor-eventselector',
    'ehbs!editor-metricselector',
    'canopsis/monitoring/src/forms/ack/controller',
    'canopsis/monitoring/src/forms/ticket/controller',
    'canopsis/monitoring/src/helpers/criticity',
    'canopsis/monitoring/src/helpers/recordcanbeack',
    'canopsis/monitoring/src/helpers/stateview',
    'canopsis/monitoring/src/helpers/statusview',
    'canopsis/monitoring/src/mixins/downtime',
    'canopsis/monitoring/src/mixins/editurlfield',
    'canopsis/monitoring/src/mixins/eventconsumer',
    'canopsis/monitoring/src/mixins/eventhistory',
    'canopsis/monitoring/src/mixins/eventnavigation',
    'canopsis/monitoring/src/mixins/history',
    'canopsis/monitoring/src/mixins/infobutton',
    'canopsis/monitoring/src/mixins/metricconsumer',
    'canopsis/monitoring/src/mixins/metricfilterable',
    'canopsis/monitoring/src/mixins/recordinfopopup',
    'canopsis/monitoring/src/mixins/sendevent',
    'ehbs!renderer-ack',
    'ehbs!renderer-cfilter',
    'ehbs!renderer-cfilterwithproperties',
    'ehbs!renderer-crecord-type',
    'ehbs!renderer-criticity',
    'ehbs!renderer-eventselector',
    'ehbs!renderer-eventtimestamp',
    'ehbs!renderer-eventtype',
    'ehbs!renderer-state',
    'ehbs!renderer-stateConnector',
    'ehbs!renderer-status',
    'canopsis/monitoring/src/reopens/routes/application',
    'link!canopsis/monitoring/src/style.css',
    'ehbs!actionbutton-editurlfield',
    'canopsis/monitoring/src/widgets/weather/controller',
    'ehbs!weather'
], function (templates) {
    templates = $(templates).filter('script');
for (var i = 0, l = templates.length; i < l; i++) {
var tpl = $(templates[i]);
Ember.TEMPLATES[tpl.attr('data-template-name')] = Ember.Handlebars.compile(tpl.text());
};
});

