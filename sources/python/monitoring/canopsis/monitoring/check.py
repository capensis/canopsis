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

from canopsis.monitoring.parser import CheckParser, PerfDataParser
from canopsis import schema as cschema

from os.path import join, expanduser, exists
from subprocess import Popen, PIPE
import json
import sys


class CheckRunner(object):
    commanddir = join(sys.prefix, 'etc', 'monitoring', 'commands')

    class CheckError(Exception):
        pass

    def __init__(self, context, command, *args, **kwargs):
        super(CheckRunner, self).__init__(*args, **kwargs)

        self.context = context
        self.command = command

    def load_config(self, command):
        """
        Check that config file exists

        :param command: Command's name
        :type command: str

        :return: schema, configuration
        """

        conffile = join(self.commanddir, '{0}.json'.format(command))
        schema_id = 'monitoringplugin.{0}'.format(command)

        if not exists(conffile):
            raise self.CheckError(
                'Impossible to find command: {0}'.format(command)
            )

        # Load configuration and validate it
        with open(conffile) as f:
            try:
                conf = json.load(f)
                schema = cschema.get(schema_id)

                if not cschema.validate(conf, schema_id):
                    raise ValueError('Impossible to validate configuration')

            except ValueError as err:
                raise self.CheckError(
                    'Invalid command {0}: {1}'.format(command, err)
                )

            except cschema.NoSchemaError:
                raise self.CheckError(
                    'Cannot validate command {0}: {1}'.format(command, err)
                )

        return schema, conf

    def build_command(self, schema, conf):
        """
        Build command from configuration and schema

        :param schema: Command's schema
        :type schema: dict

        :param conf: Command's configuration
        :type conf: dict
        """

        builder = schema['meta']['command']
        cmd = expanduser(builder['binpath'])
        cmdargs = []

        for opt in builder['args']:
            cmdopt = '-{0}'.format(opt)

            fieldName = builder['args'][opt]
            field = schema['properties'].get(fieldName, None)

            if field is not None:
                if fieldName in conf:
                    value = conf[fieldName]

                elif 'default' in field:
                    value = field['default']

                else:
                    value = None

                if value is not None:
                    if field['type'] == 'boolean':
                        if value:
                            cmdargs.append(cmdopt)

                    elif field['type'] == 'array':
                        for item in value:
                            cmdargs.append(cmdopt)
                            cmdargs.append(str(item))

                    else:
                        cmdargs.append(cmdopt)
                        cmdargs.append(str(value))

        return cmd, cmdargs

    def gen_event(self, errcode, output):
        check = CheckParser(errcode, output)
        perfdata = PerfDataParser(check.perfdata)

        event = {
            'connector': self.context['connector'],
            'connector_name': self.context['connector_name'],
            'event_type': 'check',
            'source_type': self.context['source_type'],
            'component': self.context['component'],

            'state': check.status,
            'state_type': 1,

            'output': check.text,
            'long_output': check.long_output,

            'perf_data_array': perfdata.perf_data_array
        }

        if 'resource' in self.context:
            event['resource'] = self.context['resource']

        return event

    def __call__(self):
        schema, conf = self.load_config(self.command)
        cmd, cmdargs = self.build_command(schema, conf)
        args = [cmd] + cmdargs

        p = Popen(' '.join(args), stdout=PIPE, shell=True)
        output = p.communicate()[0]
        errcode = p.returncode

        return self.gen_event(errcode, output)
