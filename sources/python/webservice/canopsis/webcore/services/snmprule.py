# -*- coding: utf-8 -*-

import os
from subprocess import Popen, PIPE
from canopsis.common.ws import route
from canopsis.snmp.rulesmanager import RulesManager
from canopsis.snmp.mibs import MibsManager
from bottle import request, HTTPError
from json import loads, dumps

manager = RulesManager()
mibmanager = MibsManager()


def exports(ws):

    rest = ws.require('rest')

    @route(ws.application.delete, payload=['ids'])
    def snmprule(ids):
        manager.remove(ids)
        ws.logger.info('Delete : {}'.format(ids))
        return True

    @route(
        ws.application.post,
        payload=['document'],
        name='snmprule/put'
    )
    def snmprule(document):
        ws.logger.debug(document)
        manager.put(document)
        return True

    @route(ws.application.post, payload=['limit', 'start', 'sort', 'filter'])
    def snmprule(limit=0, start=0, sort=None, filter={}):
        result = manager.find(
            limit=limit,
            skip=start,
            query=filter,
            sort=sort,
            with_count=True
        )
        return result

    @route(ws.application.post, payload=['limit', 'query', 'projection'])
    def snmpmib(limit=None, query={}, projection=None):
        result = mibmanager.get(
            limit=limit,
            query=query,
            projection=projection
        )
        return result

    @route(ws.application.post, payload=['filecontent'])
    def uploadmib(filecontent):

        mib_path = '~/tmp/{}'

        filenames = []

        for fileinfo in filecontent:
            filename = fileinfo['filename']
            filenames.append(filename)
            filepath = os.path.expanduser(
                mib_path.format(filename)
            )
            with open(filepath, 'w') as f:
                f.write(fileinfo['data'].encode('utf-8'))

        validation = []
        for filename in filenames:
            filepath = os.path.expanduser(
                mib_path.format(filename)
            )
            check = mibmanager.check_mib(filepath)
            validation.append((check, filename))

        if not all([x[0] for x in validation]):
            return {
                'msg': 'Uploaded mib file validation failed.',
                'data': validation
            }
        else:

            import_all = True
            notif_count = object_count = 0
            errors = []
            for filename in filenames:
                filepath = os.path.expanduser(
                    mib_path.format(filename)
                )
            try:
                counts = mibmanager.import_mib(filepath)
                notif_count += counts[0]
                object_count += counts[1]
            except Exception as e:
                errors.append((filename, e))
                import_all = False

            message = 'Updated {} objects, {} notification'.format(
                object_count,
                notif_count
            )
            return {
                'msg': message,
                'data': errors
            }

    """
    @route(ws.application.post, payload=[])
    def uploadmib():

        upload = request.files.get('file')
        filepath = os.path.expanduser('~/tmp/mibimport.mib')
        ws.logger.info('writting mib file {}'.format(filepath))

        with open(filepath, 'w') as f:
            f.write(upload.file.read())

        # Try import mib
        process = Popen(
            [
                'python',
                '-m',
                'canopsis.snmp.mibs',
                '-k',
                filepath
            ],
            stdout=PIPE
        )
        stdout, _ = process.communicate()
        ws.logger.info(stdout)

        try:
            # Depends on what does the mib parser produces
            lines = stdout.split('----------')[-1].strip().split('\n')
            # Filter output mibimport error line if any (nice summary in ui)
            lines = [line for line in lines if 'mibimport.mib' not in line]
            output = '\n'.join(lines)
        except Exception as e:
            output = 'Could not get import summary'

        return {'message': output}
    """

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
    def trap(**params):

        events, length = rest.get_records(ws, 'events_log', **params)

        event_vars = {}
        oids = []
        module_oids = []

        # Gets all mib object id information embed in events
        for event in events:

            # Snmp vars information
            snmp_vars = event.get('snmp_vars', '{}')
            snmp_vars = loads(snmp_vars)
            event_vars[event['_id']] = snmp_vars
            oids += snmp_vars.keys()

            # Snmp module oid
            module_oid = event.get('snmp_trap_oid', None)
            if module_oid:
                module_oids.append(module_oid)

        oids = list(set(oids))

        # Find information in db for oids set
        query = {
            'oid': {'$in': oids},
            'nodetype': {'$exists': 0}
        }
        mibs_infos = mibmanager.get(query=query)

        # Generate a translation
        oid_to_object = {}
        for mib_info in mibs_infos:
            oid_to_object[mib_info['oid']] = mib_info['_id'].split('::')[-1]

        # Find module information from oid if any
        query = {
            'nodetype': 'notification'
        }
        projection = {
            'moduleName': 1,
            'name': 1
        }
        module_infos = mibmanager.get(
            oids=module_oids,
            projection=projection
        )

        # Generate dict with module name information
        snmp_trap_oid_to_object = {}
        for module_info in module_infos:
            module_name = '{}::{}'.format(
                module_info['moduleName'],
                module_info['name']
            ).strip()
            oid = module_info['_id']
            if not module_name:
                module_name = oid
            snmp_trap_oid_to_object[oid] = module_name

        # Replace fetched information in event query
        for event in events:
            if event.get('snmp_vars', None):
                new_snmp_vars = {}
                for oid in event_vars[event['_id']]:
                    replace_oid = oid_to_object.get(oid, oid)
                    # put back original value with
                    # maybe object translated value
                    new_snmp_vars[replace_oid] = event_vars[event['_id']][oid]
                event['snmp_vars'] = dumps(new_snmp_vars)

                # replace snmpvarsoid byt module name information if possible
                oid = event.get('snmp_trap_oid', None)
                if oid:
                    event['snmp_trap_oid'] = snmp_trap_oid_to_object[oid]

        return (events, length)
