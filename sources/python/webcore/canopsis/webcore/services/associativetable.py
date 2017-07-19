# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from __future__ import unicode_literals

from bottle import request

from canopsis.common.associative_table.manager import AssociativeTableManager
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR


def exports(ws):

    atmanager = AssociativeTableManager(logger=ws.logger)

    @ws.application.get(
        '/api/v2/associativetable/<name>'
    )
    def get_associativetable(name):
        """
        Get this particular associative table.

        :param str name: name of the associative table
        :returns: <AssociativeTable>
        """
        content = atmanager.get(name).get_all()

        if content is None:
            return gen_json({})

        return gen_json(content.get_all())

    @ws.application.post('/api/v2/associativetable/<name>')
    def insert_associativetable(name):
        """
        Create an associative table.

        :param str name: name of the associative table
        :returns: mongo result dict
        """
        # element is a full AssociativeTable (dict) to upsert
        element = request.json

        if element is None or not isinstance(element, dict):
            return gen_json_error(
                {'description': 'nothing to insert'}, HTTP_ERROR)

        assoctable = atmanager.get(name)
        if assoctable is not None:
            return gen_json_error(
                {'description': 'already exist'}, HTTP_ERROR)

        assoctable = atmanager.create(name)

        for key, val in element.items():
            assoctable.set(key, val)

        result = atmanager.save(assoctable)

        return gen_json(result)

    @ws.application.put('/api/v2/associativetable/<name>')
    def update_associativetable(name):
        """
        Update an associative table.

        :param str name: name of the associative table
        :returns: mongo result dict
        """
        # element is a full AssociativeTable (dict) to upsert
        element = request.json

        if element is None or not isinstance(element, dict):
            return gen_json_error(
                {'description': 'nothing to update'},
                HTTP_ERROR)

        assoctable = atmanager.get(name)
        if assoctable is None:
            return gen_json_error(
                {'description': 'cannot find object'}, HTTP_ERROR)

        for key, val in element.items():
            assoctable.set(key, val)

        result = atmanager.save(assoctable)

        return gen_json(result)

    @ws.application.delete(
        '/api/v2/associativetable/<name>'
    )
    def delete_associativetable(name):
        """
        Delete a associative table, based on his id.

        :param str name: name of the associative table
        :returns: mongo result dict of the deletion
        """
        ws.logger.info('Deleting associative table: {}'.format(name))

        deletion_dict = atmanager.delete(name)

        return gen_json(deletion_dict)
