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

window.bricks['uibase'].schemasArray = ["text!canopsis/uibase/schemas/crecord.pbehaviorform.json","text!canopsis/uibase/schemas/crecord.rangecolor.json","text!canopsis/uibase/schemas/customfilter.json","text!canopsis/uibase/schemas/mixin.arraysearch.json","text!canopsis/uibase/schemas/mixin.background.json","text!canopsis/uibase/schemas/mixin.ccontext.json","text!canopsis/uibase/schemas/mixin.context.json","text!canopsis/uibase/schemas/mixin.crud.json","text!canopsis/uibase/schemas/mixin.customfilterlist.json","text!canopsis/uibase/schemas/mixin.draggablecolumns.json","text!canopsis/uibase/schemas/mixin.gridlayout.json","text!canopsis/uibase/schemas/mixin.horizontallayout.json","text!canopsis/uibase/schemas/mixin.listlinedetail.json","text!canopsis/uibase/schemas/mixin.pagination.json","text!canopsis/uibase/schemas/mixin.periodicrefresh.json","text!canopsis/uibase/schemas/mixin.responsivelist.json","text!canopsis/uibase/schemas/mixin.sortablearray.json","text!canopsis/uibase/schemas/widget.canvas.json","text!canopsis/uibase/schemas/widget.containerwidget.widgetcontainer.json","text!canopsis/uibase/schemas/widget.context.json","text!canopsis/uibase/schemas/widget.crudcontext.json","text!canopsis/uibase/schemas/widget.euewi.json","text!canopsis/uibase/schemas/widget.list.jobmanager.json","text!canopsis/uibase/schemas/widget.list.json","text!canopsis/uibase/schemas/widget.text.json","text!canopsis/uibase/schemas/widget.uimaintabcollection.json","text!canopsis/uibase/schemas/widget.wgraph.topology.json"];

define(window.bricks['uibase'].schemasArray, function () {
    for (var i = 0; i < arguments.length; i++) {
        var schemaName = window.bricks['uibase'].schemasArray[i];
        var urlPrefix = 'canopsis/uibase/schemas/';

        //remove "text!" and the brick schema folder prefix
        schemaName = schemaName.slice(5 + urlPrefix.length);
        //remove ".json at the end"
        schemaName = schemaName.slice(0, -5);
        schema = JSON.parse(arguments[i]);
        record = {
            id: schemaName,
            _id: schemaName,
            crecord_name: schemaName.split('.'),
            schema: schema
        };
        record.crecord_name = record.crecord_name[record.crecord_name.length -1];

        window.schemasToLoad.push(record);
    }
 });
