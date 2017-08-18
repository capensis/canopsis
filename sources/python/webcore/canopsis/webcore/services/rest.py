# -*- coding: utf-8 -*-
# --------------------------------
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
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------
from __future__ import unicode_literals

from bottle import HTTPError, response
from canopsis.common.ws import route

from canopsis.common.utils import ensure_iterable
from canopsis.context_graph.manager import ContextGraph
from canopsis.old.record import Record

from base64 import b64decode
from json import loads

def get_records(ws, namespace, ctype=None, _id=None, **params):
    options = {
        'limit': 20,
        'start': 0,
        'search': None,
        'filter': None,
        'sort': None,
        'query': None,
        'onlyWritable': False,
        'noInternal': False,
        'ids': [],
        'multi': None,
        'fields': {}
    }

    for key in options.keys():
        options[key] = params.get(key, options[key])

    # Ensure sort always evaluates to list
    sort = options['sort']

    if not sort:
        sort = []

    else:
        sort = ensure_iterable(sort)

    if isinstance(sort, basestring):  # NOQA
        try:
            sort = loads(sort)
        except ValueError as json_error:
            ws.logger.warning('Unable to parse sort field : {} {}'.format(
                sort, json_error
            ))
            sort = []

    # Generate MongoDB sorting query
    msort = [
        (
            item['property'],
            1 if item['direction'] == 'DESC' else -1
        )
        for item in sort if item.get('property', None) is not None
    ]

    # Generate MongoDB filter
    mfilter = {}

    if isinstance(options['filter'], list):
        for item in options['filter']:
            mfilter[item['property']] = item['value']

    elif isinstance(options['filter'], dict):
        mfilter = options['filter']

    if options['multi']:
        mfilter['crecord_type'] = {
            '$in': options['multi'].split(',')
        }

    elif ctype:
        mfilter['crecord_type'] = ctype

    if options['query']:
        # FIXME: bad query can't be indexed
        mfilter['crecord_name'] = {
            '$regex': '.*{0}.*'.format(options['query']),
            '$options': 'i'
        }

    if options['search']:
        # FIXME: bad query can't be indexed
        mfilter['_id'] = {
            '$regex': '.*{0}.*'.format(options['search']),
            '$options': 'i'
        }

    ids = options['ids'] if not _id else _id.split(',')

    # Perform query
    total = 0
    records = []

    if len(ids) > 0:
        try:
            records = ws.db.get(ids, namespace=namespace)

        except KeyError:
            records = []

        if isinstance(records, Record):
            records = [records]
            total = 1

        elif isinstance(records, list):
            total = len(records)

        else:
            total = 0

        if total == 0:
            return HTTPError(404, 'IDs not found: {0}'.format(ids))

    else:
        records, total = ws.db.find(
            mfilter,
            sort=msort,
            limit=options['limit'],
            offset=options['start'],
            with_total=True,
            namespace=namespace
        )

    # Generate output
    output = []
    noInternal = options['noInternal']

    for record in records:
        if record:
            # TODO: make use of onlyWritable
            # This can be done with canopsis.old.account, but the goal is to
            # use the new permissions/rights system to do it.

            dump = record.data.get('internal', False) if noInternal else True

            if dump:
                data = record.dump(json=True)
                data['id'] = data['_id']

                if 'next_run_time' in data:
                    data['next_run_time'] = str(data['next_run_time'])

                # TODO: Handle projection in ws.db.find()
                if options['fields']:
                    for item in data.keys():
                        if item not in options['fields']:
                            del data[item]

                output.append(data)

    return output, total


def save_records(ws, namespace, ctype, _id, items):

    records = []

    for data in items:
        m_id = data.pop('_id', None)
        mid = data.pop('id', None)
        _id = m_id or mid or _id

        record = None

        # Try to fetch existing record for update
        if _id:
            try:
                record = ws.db.get(_id, namespace=namespace)

            except KeyError:
                pass  # record is None here

        if record:
            for key in data.keys():
                record.data[key] = data[key]

            record.name = data.get('crecord_name', record.name)

        else:
            cname = data.pop('crecord_name', 'noname')
            record = Record(_id=_id, data=data, name=cname, _type=ctype)

        try:
            _id = ws.db.put(record, namespace=namespace)

            drecord = record.dump()
            drecord['_id'] = str(_id)
            drecord['id'] = drecord['_id']
            records.append(drecord)

        except Exception as err:
            ws.logger.error(u'Impossible to save record: {0}'.format(
                err
            ))

    return records


def delete_records(ws, namespace, ctype, _id, data):
    if data:
        if isinstance(data, list):
            ids = []

            for item in data:
                if isinstance(item, basestring):  # NOQA
                    ids.append(item)

                elif isinstance(item, dict):
                    item_id = item.get('_id', item.get('id', None))

                    if item_id:
                        ids.append(item_id)

        elif isinstance(data, str):
            ids = data

        elif isinstance(data, dict):
            ids = data.get('_id', data.get('id', None))

    elif _id:
        ids = [_id]

    else:
        return HTTPError(400, 'Missing ids in request')

    try:
        ws.db.remove(ids, namespace=namespace)

    except Exception as err:
        return HTTPError(500, 'Impossible to remove documents: {0}'.format(
            err
        ))

    return ids


def exports(ws):
    ctxmgr = ContextGraph()

    @route(ws.application.get, name='rest/indexes', response=lambda r, a: r)
    def indexes(collection):
        storage = ws.db.get_backend(collection)
        indexes = storage.index_information()

        return {'collection': collection, 'indexes': indexes}

    @route(ws.application.get, name='rest/media', response=lambda r, a: r)
    def media(namespace, _id):
        try:
            raw = ws.db.get(
                _id,
                namespace=namespace,
                mfields=[
                    'media_bin',
                    'media_name',
                    'media_type'
                ],
                ignore_bin=False
            )

        except KeyError as err:
            return HTTPError(404, str(err))

        try:
            media_type = raw['media_type']
            media_name = raw['media_name']
            media_bin = raw['media_bin']

        except KeyError as err:
            return HTTPError(500, str(err))

        cdisp = 'attachment; filename="{0}"'.format(media_name)
        response.headers['Content-Disposition'] = cdisp
        response.headers['Content-Type'] = media_type

        return b64decode(media_bin)

    @route(
        ws.application.get,
        payload=[
            'limit',
            'start',
            'search',
            'filter',
            'sort',
            'query',
            'onlyWritable',
            'noInternal',
            'ids',
            'multi',
            'fields'
        ],
        adapt=False
    )
    def rest(namespace, ctype=None, _id=None, **params):
        records, nrecords = get_records(
            ws, namespace,
            ctype=ctype, _id=_id,
            **params
        )

        for record in records:
            if record['crecord_type'] == 'event':
                eid = ''
                if 'resource' in record.keys():
                    eid = '/{0}/{1}/{2}/{3}/{4}'.format(
                        record['source_type'],
                        record['connector'],
                        record['connector_name'],
                        record['component'],
                        record['resource']
                    )
                else:
                    eid = '/{0}/{1}/{2}/{3}'.format(
                        record['source_type'],
                        record['connector'],
                        record['connector_name'],
                        record['component']
                    )
                record['entity_id'] = eid

        return records, nrecords

    @route(ws.application.put, raw_body=True, adapt=False)  # NOQA
    def rest(namespace, ctype, _id=None, body='[]', **kwargs):
        try:
            items = ensure_iterable(loads(body))

        except ValueError as err:
            return HTTPError(500, 'Impossible to parse body: {0}'.format(err))

        return save_records(ws, namespace, ctype, _id, items)

    @route(ws.application.post, raw_body=True, adapt=False)  # NOQA
    def rest(namespace, ctype, _id=None, body='[]', **kwargs):
        try:
            items = ensure_iterable(loads(body))

        except ValueError as err:
            return HTTPError(500, 'Impossible to parse body: {0}'.format(err))

        return save_records(ws, namespace, ctype, _id, items)

    @route(ws.application.delete, raw_body=True, adapt=False)  # NOQA
    def rest(namespace, ctype, _id=None, body='[]', **kwargs):
        try:
            data = loads(body)

        except ValueError:
            data = None
        
        return delete_records(ws, namespace, ctype, _id, data)
