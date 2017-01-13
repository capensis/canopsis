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

window.bricks['brick-listalarm'].schemasArray = ["text!canopsis/brick-listalarm/schemas/crecord.alerts.json","text!canopsis/brick-listalarm/schemas/widget.listalarm.json"];

define(window.bricks['brick-listalarm'].schemasArray, function () {
    for (var i = 0; i < arguments.length; i++) {
        var schemaName = window.bricks['brick-listalarm'].schemasArray[i];
        var urlPrefix = 'canopsis/brick-listalarm/schemas/';

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
