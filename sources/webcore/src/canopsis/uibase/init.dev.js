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
        'components/component-actionbutton': 'canopsis/uibase/src/components/actionbutton/template',
        'components/component-actionfilter': 'canopsis/uibase/src/components/actionfilter/template',
        'components/component-arrayclassifiedcrecordselector': 'canopsis/uibase/src/components/arrayclassifiedcrecordselector/template',
        'components/component-arrayeditor': 'canopsis/uibase/src/components/arrayeditor/template',
        'components/component-classifiedcrecordselector': 'canopsis/uibase/src/components/classifiedcrecordselector/template',
        'components/component-classifieditemselector': 'canopsis/uibase/src/components/classifieditemselector/template',
        'components/component-colorchooser': 'canopsis/uibase/src/components/colorchooser/template',
        'components/component-colpick': 'canopsis/uibase/src/components/colpick/template',
        'components/component-contextselector': 'canopsis/uibase/src/components/contextselector/template',
        'components/component-crudentitydetails': 'canopsis/uibase/src/components/crudentitydetails/template',
        'components/component-dateinterval': 'canopsis/uibase/src/components/dateinterval/template',
        'components/component-dictclassifiedcrecordselector': 'canopsis/uibase/src/components/dictclassifiedcrecordselector/template',
        'components/component-dropdownbutton': 'canopsis/uibase/src/components/dropdownbutton/template',
        'components/component-dropdownbuttoncontent': 'canopsis/uibase/src/components/dropdownbuttoncontent/template',
        'components/component-dropdownbuttonheader': 'canopsis/uibase/src/components/dropdownbuttonheader/template',
        'components/component-dropdownbuttonoverview': 'canopsis/uibase/src/components/dropdownbuttonoverview/template',
        'components/component-dropdownbuttontitle': 'canopsis/uibase/src/components/dropdownbuttontitle/template',
        'components/component-elementidselectorwithoptions': 'canopsis/uibase/src/components/elementidselectorwithoptions/template',
        'components/component-eventkey': 'canopsis/uibase/src/components/eventkey/template',
        'components/component-filefield': 'canopsis/uibase/src/components/filefield/template',
        'components/component-filterclause': 'canopsis/uibase/src/components/filterclause/template',
        'components/component-formulaeditor': 'canopsis/uibase/src/components/formulaeditor/template',
        'components/component-ical': 'canopsis/uibase/src/components/ical/template',
        'components/component-if_equal_component': 'canopsis/uibase/src/components/if_equal_component/template',
        'components/component-labelledlink': 'canopsis/uibase/src/components/labelledlink/template',
        'components/component-linklist': 'canopsis/uibase/src/components/linklist/template',
        'components/component-mixinselector': 'canopsis/uibase/src/components/mixinselector/template',
        'components/component-modelselect': 'canopsis/uibase/src/components/modelselect/template',
        'components/component-nestedarrayeditor': 'canopsis/uibase/src/components/nestedarrayeditor/template',
        'components/component-password': 'canopsis/uibase/src/components/password/template',
        'components/component-propertiestopopup': 'canopsis/uibase/src/components/propertiestopopup/template',
        'components/component-restobjectcombo': 'canopsis/uibase/src/components/restobjectcombo/template',
        'components/component-searchbar': 'canopsis/uibase/src/components/searchbar/template',
        'components/component-simpledicteditor': 'canopsis/uibase/src/components/simpledicteditor/template',
        'components/component-slider': 'canopsis/uibase/src/components/slider/template',
        'components/component-stringclassifiedcrecordselector': 'canopsis/uibase/src/components/stringclassifiedcrecordselector/template',
        'components/component-tabcontent': 'canopsis/uibase/src/components/tabcontent/template',
        'components/component-tabheader': 'canopsis/uibase/src/components/tabheader/template',
        'components/component-table': 'canopsis/uibase/src/components/table/template',
        'components/component-tabs': 'canopsis/uibase/src/components/tabs/template',
        'components/component-tabscontentgroup': 'canopsis/uibase/src/components/tabscontentgroup/template',
        'components/component-tabsheadergroup': 'canopsis/uibase/src/components/tabsheadergroup/template',
        'components/component-textwithsortoption': 'canopsis/uibase/src/components/textwithsortoption/template',
        'components/component-timestamptooltiped': 'canopsis/uibase/src/components/timestamptooltiped/template',
        'components/component-typedvalue': 'canopsis/uibase/src/components/typedvalue/template',
        'editor-actionfilter': 'canopsis/uibase/src/editors/editor-actionfilter',
        'editor-array': 'canopsis/uibase/src/editors/editor-array',
        'editor-arrayclassifiedcrecordselector': 'canopsis/uibase/src/editors/editor-arrayclassifiedcrecordselector',
        'editor-boolean': 'canopsis/uibase/src/editors/editor-boolean',
        'editor-color': 'canopsis/uibase/src/editors/editor-color',
        'editor-contextselector': 'canopsis/uibase/src/editors/editor-contextselector',
        'editor-dateinterval': 'canopsis/uibase/src/editors/editor-dateinterval',
        'editor-defaultpropertyeditor': 'canopsis/uibase/src/editors/editor-defaultpropertyeditor',
        'editor-dictclassifiedcrecordselector': 'canopsis/uibase/src/editors/editor-dictclassifiedcrecordselector',
        'editor-duration': 'canopsis/uibase/src/editors/editor-duration',
        'editor-durationWithUnits': 'canopsis/uibase/src/editors/editor-durationWithUnits',
        'editor-elementidselectorwithoptions': 'canopsis/uibase/src/editors/editor-elementidselectorwithoptions',
        'editor-error': 'canopsis/uibase/src/editors/editor-error',
        'editor-eventkey': 'canopsis/uibase/src/editors/editor-eventkey',
        'editor-integer': 'canopsis/uibase/src/editors/editor-integer',
        'editor-json': 'canopsis/uibase/src/editors/editor-json',
        'editor-labelandviewid': 'canopsis/uibase/src/editors/editor-labelandviewid',
        'editor-labelledlink': 'canopsis/uibase/src/editors/editor-labelledlink',
        'editor-mail': 'canopsis/uibase/src/editors/editor-mail',
        'editor-mixins': 'canopsis/uibase/src/editors/editor-mixins',
        'editor-modelselect': 'canopsis/uibase/src/editors/editor-modelselect',
        'editor-nestedarray': 'canopsis/uibase/src/editors/editor-nestedarray',
        'editor-password': 'canopsis/uibase/src/editors/editor-password',
        'editor-passwordmd5': 'canopsis/uibase/src/editors/editor-passwordmd5',
        'editor-passwordsha1': 'canopsis/uibase/src/editors/editor-passwordsha1',
        'editor-restobject': 'canopsis/uibase/src/editors/editor-restobject',
        'editor-richtext': 'canopsis/uibase/src/editors/editor-richtext',
        'editor-separator': 'canopsis/uibase/src/editors/editor-separator',
        'editor-serieformula': 'canopsis/uibase/src/editors/editor-serieformula',
        'editor-simpledict': 'canopsis/uibase/src/editors/editor-simpledict',
        'editor-simplelist': 'canopsis/uibase/src/editors/editor-simplelist',
        'editor-slider': 'canopsis/uibase/src/editors/editor-slider',
        'editor-sortable': 'canopsis/uibase/src/editors/editor-sortable',
        'editor-source': 'canopsis/uibase/src/editors/editor-source',
        'editor-state': 'canopsis/uibase/src/editors/editor-state',
        'editor-stringclassifiedcrecordselector': 'canopsis/uibase/src/editors/editor-stringclassifiedcrecordselector',
        'editor-stringpair': 'canopsis/uibase/src/editors/editor-stringpair',
        'editor-tags': 'canopsis/uibase/src/editors/editor-tags',
        'editor-textarea': 'canopsis/uibase/src/editors/editor-textarea',
        'editor-timestamp': 'canopsis/uibase/src/editors/editor-timestamp',
        'editor-typedvalue': 'canopsis/uibase/src/editors/editor-typedvalue',
        'editor-userpreference': 'canopsis/uibase/src/editors/editor-userpreference',
        'pbehaviorform': 'canopsis/uibase/src/forms/pbehavior/pbehaviorform',
        'renderer-actionfilter': 'canopsis/uibase/src/renderers/renderer-actionfilter',
        'renderer-boolean': 'canopsis/uibase/src/renderers/renderer-boolean',
        'renderer-color': 'canopsis/uibase/src/renderers/renderer-color',
        'renderer-conf': 'canopsis/uibase/src/renderers/renderer-conf',
        'renderer-ellipsis': 'canopsis/uibase/src/renderers/renderer-ellipsis',
        'renderer-labelledlink': 'canopsis/uibase/src/renderers/renderer-labelledlink',
        'renderer-mail': 'canopsis/uibase/src/renderers/renderer-mail',
        'renderer-object': 'canopsis/uibase/src/renderers/renderer-object',
        'renderer-percent': 'canopsis/uibase/src/renderers/renderer-percent',
        'renderer-recordinfopopup': 'canopsis/uibase/src/renderers/renderer-recordinfopopup',
        'renderer-richtext': 'canopsis/uibase/src/renderers/renderer-richtext',
        'renderer-source': 'canopsis/uibase/src/renderers/renderer-source',
        'renderer-subprocess': 'canopsis/uibase/src/renderers/renderer-subprocess',
        'renderer-tags': 'canopsis/uibase/src/renderers/renderer-tags',
        'renderer-timestamp': 'canopsis/uibase/src/renderers/renderer-timestamp',
        'renderer-translator': 'canopsis/uibase/src/renderers/renderer-translator',
        'actionbutton-ack': 'canopsis/uibase/src/templates/actionbutton-ack',
        'actionbutton-ackselection': 'canopsis/uibase/src/templates/actionbutton-ackselection',
        'actionbutton-cancel': 'canopsis/uibase/src/templates/actionbutton-cancel',
        'actionbutton-cancelselection': 'canopsis/uibase/src/templates/actionbutton-cancelselection',
        'actionbutton-changestate': 'canopsis/uibase/src/templates/actionbutton-changestate',
        'actionbutton-create': 'canopsis/uibase/src/templates/actionbutton-create',
        'actionbutton-createpbehavior': 'canopsis/uibase/src/templates/actionbutton-createpbehavior',
        'actionbutton-duplicate': 'canopsis/uibase/src/templates/actionbutton-duplicate',
        'actionbutton-edit': 'canopsis/uibase/src/templates/actionbutton-edit',
        'actionbutton-eventnavigation': 'canopsis/uibase/src/templates/actionbutton-eventnavigation',
        'actionbutton-foldable': 'canopsis/uibase/src/templates/actionbutton-foldable',
        'actionbutton-history': 'canopsis/uibase/src/templates/actionbutton-history',
        'actionbutton-incident': 'canopsis/uibase/src/templates/actionbutton-incident',
        'actionbutton-info': 'canopsis/uibase/src/templates/actionbutton-info',
        'actionbutton-remove': 'canopsis/uibase/src/templates/actionbutton-remove',
        'actionbutton-removeselection': 'canopsis/uibase/src/templates/actionbutton-removeselection',
        'actionbutton-show': 'canopsis/uibase/src/templates/actionbutton-show',
        'actionbutton-ticketnumber': 'canopsis/uibase/src/templates/actionbutton-ticketnumber',
        'column-unfold': 'canopsis/uibase/src/templates/column-unfold',
        'consolemanagerstatusmenu': 'canopsis/uibase/src/templates/consolemanagerstatusmenu',
        'crecordform': 'canopsis/uibase/src/templates/crecordform',
        'customfilters': 'canopsis/uibase/src/templates/customfilters',
        'documentation': 'canopsis/uibase/src/templates/documentation',
        'draggableheaders': 'canopsis/uibase/src/templates/draggableheaders',
        'formbutton-ack': 'canopsis/uibase/src/templates/formbutton-ack',
        'formbutton-ackandproblem': 'canopsis/uibase/src/templates/formbutton-ackandproblem',
        'formbutton-cancel': 'canopsis/uibase/src/templates/formbutton-cancel',
        'formbutton-delete': 'canopsis/uibase/src/templates/formbutton-delete',
        'formbutton-incident': 'canopsis/uibase/src/templates/formbutton-incident',
        'formbutton-inspectform': 'canopsis/uibase/src/templates/formbutton-inspectform',
        'formbutton-next': 'canopsis/uibase/src/templates/formbutton-next',
        'formbutton-previous': 'canopsis/uibase/src/templates/formbutton-previous',
        'formbutton-submit': 'canopsis/uibase/src/templates/formbutton-submit',
        'formwrapper': 'canopsis/uibase/src/templates/formwrapper',
        'gridlayout': 'canopsis/uibase/src/templates/gridlayout',
        'groupedrowslistlayout': 'canopsis/uibase/src/templates/groupedrowslistlayout',
        'groupedrowslistlinelayout': 'canopsis/uibase/src/templates/groupedrowslistlinelayout',
        'horizontallayout': 'canopsis/uibase/src/templates/horizontallayout',
        'index': 'canopsis/uibase/src/templates/index',
        'itemsperpage': 'canopsis/uibase/src/templates/itemsperpage',
        'lightlayout': 'canopsis/uibase/src/templates/lightlayout',
        'listline': 'canopsis/uibase/src/templates/listline',
        'loading': 'canopsis/uibase/src/templates/loading',
        'loadingindicator': 'canopsis/uibase/src/templates/loadingindicator',
        'menu': 'canopsis/uibase/src/templates/menu',
        'mixineditdropdown': 'canopsis/uibase/src/templates/mixineditdropdown',
        'mixinselector-itempartial': 'canopsis/uibase/src/templates/mixinselector-itempartial',
        'notificationsstatusmenu': 'canopsis/uibase/src/templates/notificationsstatusmenu',
        'pagination-infos': 'canopsis/uibase/src/templates/pagination-infos',
        'pagination': 'canopsis/uibase/src/templates/pagination',
        'partialslot': 'canopsis/uibase/src/templates/partialslot',
        'presettoolbar': 'canopsis/uibase/src/templates/presettoolbar',
        'promisemanagerstatusmenu': 'canopsis/uibase/src/templates/promisemanagerstatusmenu',
        'recordinfopopup': 'canopsis/uibase/src/templates/recordinfopopup',
        'requirejsmockingstatusmenu': 'canopsis/uibase/src/templates/requirejsmockingstatusmenu',
        'schemamanagerstatusmenu': 'canopsis/uibase/src/templates/schemamanagerstatusmenu',
        'screentoolstatusmenu': 'canopsis/uibase/src/templates/screentoolstatusmenu',
        'search': 'canopsis/uibase/src/templates/search',
        'stackedcolumns': 'canopsis/uibase/src/templates/stackedcolumns',
        'tablayout': 'canopsis/uibase/src/templates/tablayout',
        'tabledraggableth': 'canopsis/uibase/src/templates/tabledraggableth',
        'titlebarbutton-duplicate': 'canopsis/uibase/src/templates/titlebarbutton-duplicate',
        'titlebarbutton-minimize': 'canopsis/uibase/src/templates/titlebarbutton-minimize',
        'titlebarbutton-movedown': 'canopsis/uibase/src/templates/titlebarbutton-movedown',
        'titlebarbutton-moveleft': 'canopsis/uibase/src/templates/titlebarbutton-moveleft',
        'titlebarbutton-moveright': 'canopsis/uibase/src/templates/titlebarbutton-moveright',
        'titlebarbutton-moveup': 'canopsis/uibase/src/templates/titlebarbutton-moveup',
        'titlebarbutton-widgeterrors': 'canopsis/uibase/src/templates/titlebarbutton-widgeterrors',
        'titlebarbutton-widgetfullscreen': 'canopsis/uibase/src/templates/titlebarbutton-widgetfullscreen',
        'titlebarbutton-widgetrefresh': 'canopsis/uibase/src/templates/titlebarbutton-widgetrefresh',
        'userstatusmenu': 'canopsis/uibase/src/templates/userstatusmenu',
        'userview': 'canopsis/uibase/src/templates/userview',
        'verticallayout': 'canopsis/uibase/src/templates/verticallayout',
        'widget': 'canopsis/uibase/src/templates/widget',
        'widgetslot-default': 'canopsis/uibase/src/templates/widgetslot-default',
        'widgetslot-grey': 'canopsis/uibase/src/templates/widgetslot-grey',
        'widgetslot-light': 'canopsis/uibase/src/templates/widgetslot-light',
        'widgettitlebar': 'canopsis/uibase/src/templates/widgettitlebar',
        'crudcontext': 'canopsis/uibase/src/widgets/crudcontext/crudcontext',
        'list': 'canopsis/uibase/src/widgets/list/list',
        'textwidget': 'canopsis/uibase/src/widgets/text/textwidget',
        'topology': 'canopsis/uibase/src/widgets/topology/topology',
        'uimaintabcollection': 'canopsis/uibase/src/widgets/uimaintabcollection/uimaintabcollection',
        'widgetcontainer': 'canopsis/uibase/src/widgets/widgetcontainer/widgetcontainer',

    }
});

 define([
    'canopsis/uibase/src/components/actionbutton/component',
    'ehbs!components/component-actionbutton',
    'canopsis/uibase/src/components/actionfilter/component',
    'ehbs!components/component-actionfilter',
    'canopsis/uibase/src/components/arrayclassifiedcrecordselector/component',
    'ehbs!components/component-arrayclassifiedcrecordselector',
    'canopsis/uibase/src/components/arrayeditor/component',
    'ehbs!components/component-arrayeditor',
    'canopsis/uibase/src/components/classifiedcrecordselector/component',
    'ehbs!components/component-classifiedcrecordselector',
    'canopsis/uibase/src/components/classifieditemselector/component',
    'ehbs!components/component-classifieditemselector',
    'canopsis/uibase/src/components/colorchooser/component',
    'ehbs!components/component-colorchooser',
    'canopsis/uibase/src/components/colpick/component',
    'ehbs!components/component-colpick',
    'canopsis/uibase/src/components/contextselector/component',
    'ehbs!components/component-contextselector',
    'canopsis/uibase/src/components/crudentitydetails/component',
    'ehbs!components/component-crudentitydetails',
    'canopsis/uibase/src/components/dateinterval/component',
    'ehbs!components/component-dateinterval',
    'canopsis/uibase/src/components/dictclassifiedcrecordselector/component',
    'ehbs!components/component-dictclassifiedcrecordselector',
    'canopsis/uibase/src/components/dropdownbutton/component',
    'ehbs!components/component-dropdownbutton',
    'canopsis/uibase/src/components/dropdownbuttoncontent/component',
    'ehbs!components/component-dropdownbuttoncontent',
    'canopsis/uibase/src/components/dropdownbuttonheader/component',
    'ehbs!components/component-dropdownbuttonheader',
    'canopsis/uibase/src/components/dropdownbuttonoverview/component',
    'ehbs!components/component-dropdownbuttonoverview',
    'canopsis/uibase/src/components/dropdownbuttontitle/component',
    'ehbs!components/component-dropdownbuttontitle',
    'canopsis/uibase/src/components/elementidselectorwithoptions/component',
    'ehbs!components/component-elementidselectorwithoptions',
    'canopsis/uibase/src/components/eventkey/component',
    'ehbs!components/component-eventkey',
    'canopsis/uibase/src/components/filefield/component',
    'ehbs!components/component-filefield',
    'canopsis/uibase/src/components/filterclause/component',
    'ehbs!components/component-filterclause',
    'canopsis/uibase/src/components/formulaeditor/component',
    'ehbs!components/component-formulaeditor',
    'canopsis/uibase/src/components/ical/component',
    'ehbs!components/component-ical',
    'canopsis/uibase/src/components/if_equal_component/component',
    'ehbs!components/component-if_equal_component',
    'canopsis/uibase/src/components/labelledlink/component',
    'ehbs!components/component-labelledlink',
    'canopsis/uibase/src/components/linklist/component',
    'ehbs!components/component-linklist',
    'canopsis/uibase/src/components/mixinselector/component',
    'ehbs!components/component-mixinselector',
    'canopsis/uibase/src/components/modelselect/component',
    'ehbs!components/component-modelselect',
    'canopsis/uibase/src/components/nestedarrayeditor/component',
    'ehbs!components/component-nestedarrayeditor',
    'canopsis/uibase/src/components/password/component',
    'ehbs!components/component-password',
    'canopsis/uibase/src/components/propertiestopopup/component',
    'ehbs!components/component-propertiestopopup',
    'canopsis/uibase/src/components/restobjectcombo/component',
    'ehbs!components/component-restobjectcombo',
    'canopsis/uibase/src/components/searchbar/component',
    'ehbs!components/component-searchbar',
    'canopsis/uibase/src/components/simpledicteditor/component',
    'ehbs!components/component-simpledicteditor',
    'canopsis/uibase/src/components/slider/component',
    'ehbs!components/component-slider',
    'canopsis/uibase/src/components/stringclassifiedcrecordselector/component',
    'ehbs!components/component-stringclassifiedcrecordselector',
    'canopsis/uibase/src/components/tabcontent/component',
    'ehbs!components/component-tabcontent',
    'canopsis/uibase/src/components/tabheader/component',
    'ehbs!components/component-tabheader',
    'canopsis/uibase/src/components/table/component',
    'ehbs!components/component-table',
    'canopsis/uibase/src/components/tabs/component',
    'ehbs!components/component-tabs',
    'canopsis/uibase/src/components/tabscontentgroup/component',
    'ehbs!components/component-tabscontentgroup',
    'canopsis/uibase/src/components/tabsheadergroup/component',
    'ehbs!components/component-tabsheadergroup',
    'canopsis/uibase/src/components/textwithsortoption/component',
    'ehbs!components/component-textwithsortoption',
    'canopsis/uibase/src/components/timestamptooltiped/component',
    'ehbs!components/component-timestamptooltiped',
    'canopsis/uibase/src/components/typedvalue/component',
    'ehbs!components/component-typedvalue',
    'ehbs!editor-actionfilter',
    'ehbs!editor-array',
    'ehbs!editor-arrayclassifiedcrecordselector',
    'ehbs!editor-boolean',
    'ehbs!editor-color',
    'ehbs!editor-contextselector',
    'ehbs!editor-dateinterval',
    'ehbs!editor-defaultpropertyeditor',
    'ehbs!editor-dictclassifiedcrecordselector',
    'ehbs!editor-duration',
    'ehbs!editor-durationWithUnits',
    'ehbs!editor-elementidselectorwithoptions',
    'ehbs!editor-error',
    'ehbs!editor-eventkey',
    'ehbs!editor-integer',
    'ehbs!editor-json',
    'ehbs!editor-labelandviewid',
    'ehbs!editor-labelledlink',
    'ehbs!editor-mail',
    'ehbs!editor-mixins',
    'ehbs!editor-modelselect',
    'ehbs!editor-nestedarray',
    'ehbs!editor-password',
    'ehbs!editor-passwordmd5',
    'ehbs!editor-passwordsha1',
    'ehbs!editor-restobject',
    'ehbs!editor-richtext',
    'ehbs!editor-separator',
    'ehbs!editor-serieformula',
    'ehbs!editor-simpledict',
    'ehbs!editor-simplelist',
    'ehbs!editor-slider',
    'ehbs!editor-sortable',
    'ehbs!editor-source',
    'ehbs!editor-state',
    'ehbs!editor-stringclassifiedcrecordselector',
    'ehbs!editor-stringpair',
    'ehbs!editor-tags',
    'ehbs!editor-textarea',
    'ehbs!editor-timestamp',
    'ehbs!editor-typedvalue',
    'ehbs!editor-userpreference',
    'canopsis/uibase/src/forms/pbehavior/controller',
    'ehbs!pbehaviorform',
    'canopsis/uibase/src/helpers/color',
    'canopsis/uibase/src/helpers/compare',
    'canopsis/uibase/src/helpers/ellipsis',
    'canopsis/uibase/src/helpers/eventtype',
    'canopsis/uibase/src/helpers/glyphicon',
    'canopsis/uibase/src/helpers/humanreadable',
    'canopsis/uibase/src/helpers/infosdetails',
    'canopsis/uibase/src/helpers/interval2html',
    'canopsis/uibase/src/helpers/json2html',
    'canopsis/uibase/src/helpers/linksdetails',
    'canopsis/uibase/src/helpers/logo',
    'canopsis/uibase/src/helpers/percent',
    'canopsis/uibase/src/helpers/rights',
    'canopsis/uibase/src/helpers/sorticon',
    'canopsis/uibase/src/helpers/timeSince',
    'canopsis/uibase/src/helpers/timestamp',
    'canopsis/uibase/src/helpers/timestampfulldate',
    'canopsis/uibase/src/mixins/arraysearch',
    'canopsis/uibase/src/mixins/background',
    'canopsis/uibase/src/mixins/contextarraysearch',
    'canopsis/uibase/src/mixins/crud',
    'canopsis/uibase/src/mixins/customfilterlist',
    'canopsis/uibase/src/mixins/customsendevent',
    'canopsis/uibase/src/mixins/draggablecolumns',
    'canopsis/uibase/src/mixins/gridlayout',
    'canopsis/uibase/src/mixins/horizontallayout',
    'canopsis/uibase/src/mixins/lightlayout',
    'canopsis/uibase/src/mixins/listlinedetail',
    'canopsis/uibase/src/mixins/pagination',
    'canopsis/uibase/src/mixins/periodicrefresh',
    'canopsis/uibase/src/mixins/responsivelist',
    'canopsis/uibase/src/mixins/showviewbutton',
    'canopsis/uibase/src/mixins/sortablearray',
    'canopsis/uibase/src/mixins/tablayout',
    'canopsis/uibase/src/mixins/verticallayout',
    'canopsis/uibase/src/mixins/widgetfullscreen',
    'canopsis/uibase/src/mixins/widgetrefresh',
    'ehbs!renderer-actionfilter',
    'ehbs!renderer-boolean',
    'ehbs!renderer-color',
    'ehbs!renderer-conf',
    'ehbs!renderer-ellipsis',
    'ehbs!renderer-labelledlink',
    'ehbs!renderer-mail',
    'ehbs!renderer-object',
    'ehbs!renderer-percent',
    'ehbs!renderer-recordinfopopup',
    'ehbs!renderer-richtext',
    'ehbs!renderer-source',
    'ehbs!renderer-subprocess',
    'ehbs!renderer-tags',
    'ehbs!renderer-timestamp',
    'ehbs!renderer-translator',
    'link!canopsis/uibase/src/style.css',
    'ehbs!actionbutton-ack',
    'ehbs!actionbutton-ackselection',
    'ehbs!actionbutton-cancel',
    'ehbs!actionbutton-cancelselection',
    'ehbs!actionbutton-changestate',
    'ehbs!actionbutton-create',
    'ehbs!actionbutton-createpbehavior',
    'ehbs!actionbutton-duplicate',
    'ehbs!actionbutton-edit',
    'ehbs!actionbutton-eventnavigation',
    'ehbs!actionbutton-foldable',
    'ehbs!actionbutton-history',
    'ehbs!actionbutton-incident',
    'ehbs!actionbutton-info',
    'ehbs!actionbutton-remove',
    'ehbs!actionbutton-removeselection',
    'ehbs!actionbutton-show',
    'ehbs!actionbutton-ticketnumber',
    'ehbs!column-unfold',
    'ehbs!consolemanagerstatusmenu',
    'ehbs!crecordform',
    'ehbs!customfilters',
    'ehbs!documentation',
    'ehbs!draggableheaders',
    'ehbs!formbutton-ack',
    'ehbs!formbutton-ackandproblem',
    'ehbs!formbutton-cancel',
    'ehbs!formbutton-delete',
    'ehbs!formbutton-incident',
    'ehbs!formbutton-inspectform',
    'ehbs!formbutton-next',
    'ehbs!formbutton-previous',
    'ehbs!formbutton-submit',
    'ehbs!formwrapper',
    'ehbs!gridlayout',
    'ehbs!groupedrowslistlayout',
    'ehbs!groupedrowslistlinelayout',
    'ehbs!horizontallayout',
    'ehbs!index',
    'ehbs!itemsperpage',
    'ehbs!lightlayout',
    'ehbs!listline',
    'ehbs!loading',
    'ehbs!loadingindicator',
    'ehbs!menu',
    'ehbs!mixineditdropdown',
    'ehbs!mixinselector-itempartial',
    'ehbs!notificationsstatusmenu',
    'ehbs!pagination-infos',
    'ehbs!pagination',
    'ehbs!partialslot',
    'ehbs!presettoolbar',
    'ehbs!promisemanagerstatusmenu',
    'ehbs!recordinfopopup',
    'ehbs!requirejsmockingstatusmenu',
    'ehbs!schemamanagerstatusmenu',
    'ehbs!screentoolstatusmenu',
    'ehbs!search',
    'ehbs!stackedcolumns',
    'ehbs!tablayout',
    'ehbs!tabledraggableth',
    'ehbs!titlebarbutton-duplicate',
    'ehbs!titlebarbutton-minimize',
    'ehbs!titlebarbutton-movedown',
    'ehbs!titlebarbutton-moveleft',
    'ehbs!titlebarbutton-moveright',
    'ehbs!titlebarbutton-moveup',
    'ehbs!titlebarbutton-widgeterrors',
    'ehbs!titlebarbutton-widgetfullscreen',
    'ehbs!titlebarbutton-widgetrefresh',
    'ehbs!userstatusmenu',
    'ehbs!userview',
    'ehbs!verticallayout',
    'ehbs!widget',
    'ehbs!widgetslot-default',
    'ehbs!widgetslot-grey',
    'ehbs!widgetslot-light',
    'ehbs!widgettitlebar',
    'canopsis/uibase/src/widgets/crudcontext/controller',
    'ehbs!crudcontext',
    'canopsis/uibase/src/widgets/list/controller',
    'ehbs!list',
    'canopsis/uibase/src/widgets/text/controller',
    'ehbs!textwidget',
    'canopsis/uibase/src/widgets/topology/adapter',
    'canopsis/uibase/src/widgets/topology/controller',
    'canopsis/uibase/src/widgets/topology/layout/cluster',
    'canopsis/uibase/src/widgets/topology/layout/force',
    'canopsis/uibase/src/widgets/topology/layout/pack',
    'canopsis/uibase/src/widgets/topology/layout/partition',
    'canopsis/uibase/src/widgets/topology/layout/tree',
    'canopsis/uibase/src/widgets/topology/layout',
    'link!canopsis/uibase/src/widgets/topology/style.css',
    'ehbs!topology',
    'canopsis/uibase/src/widgets/topology/view',
    'canopsis/uibase/src/widgets/uimaintabcollection/controller',
    'ehbs!uimaintabcollection',
    'canopsis/uibase/src/widgets/widgetcontainer/controller',
    'ehbs!widgetcontainer',
    'canopsis/uibase/requirejs-modules/externals.conf'
], function (templates) {
    templates = $(templates).filter('script');
for (var i = 0, l = templates.length; i < l; i++) {
var tpl = $(templates[i]);
Ember.TEMPLATES[tpl.attr('data-template-name')] = Ember.Handlebars.compile(tpl.text());
};
});

