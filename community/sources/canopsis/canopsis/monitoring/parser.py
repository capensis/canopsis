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

from re import compile as re_compile


class CheckParser(object):
    """
    This class is used to parse the output of a Monitoring check plugin.

    It's an RAII class, the result is available via its instance properties.
    """

    def __init__(self, errcode, output, *args, **kwargs):
        """
        :param errcode: Exit code of monitoring check.
        :type errcode: int

        :param output: Full output of monitoring check.
        :type output: str
        """

        super(CheckParser, self).__init__(*args, **kwargs)

        self._status = errcode
        self._text = ''
        self._perfdata = ''
        self._long_output = ''

        lines = output.splitlines()
        perfdata = False

        # Extract text output and perfdata from first line
        if len(lines) > 0:
            parts = lines[0].split('|')
            self._text = parts[0].strip()

            self._perfdata = parts[-1].strip() if len(parts) > 1 else ''

        # Extract long output and extras perfdatas
        for line in lines[1:]:
            parts = line.split('|')

            # If there is perfdata on this line, then it's no longer a long output
            if len(parts) > 1:
                perfdata = True

                self._long_output += '{0}\n'.format(parts[0].strip())

            # If we are not adding perfdata, append to long output
            if not perfdata:
                self._long_output += '{0}\n'.format(parts[0].strip())

            # Else, if perfdata is already defined, append to it with correct syntax
            elif self._perfdata:
                self._perfdata += ',{0}'.format(parts[-1].strip())

            # Else, just define it
            else:
                self._perfdata = parts[-1].strip()

        self._long_output = self._long_output.strip('\n')
        self._text = self._text.strip('\n')
        self._perfdata = self._perfdata.strip('\n')

    @property
    def status(self):
        return self._status

    @property
    def text(self):
        return self._text

    @property
    def long_output(self):
        return self._long_output

    @property
    def perfdata(self):
        return self._perfdata


class PerfDataParser(object):
    """
    Parse Monitoring performance data string.

    It's a RAII class, the result is available via its instance properties.
    """

    def __init__(self, perfdata, *args, **kwargs):
        """
        :param perfdata: Nagios performance data string
        :type perfdata: str
        """

        super(PerfDataParser, self).__init__(*args, **kwargs)

        regex = re_compile(
            "('?([0-9A-Za-z/\\\:\.%%\-{}\?\[\]\(\)_ ]*)'?=(\-?[0-9.,]*)(([A-Za-z%%/]*))(;@?(\-?[0-9.,]*):?)?(;@?(\-?[0-9.,]*):?)?(;@?(\-?[0-9.,]*):?)?(;@?(\-?[0-9.,]*):?)?(;? ?))"
        )

        perfs = regex.split(perfdata)

        pfmap = {
            2: ('metric', lambda nfo: nfo),
            3: ('value', lambda nfo: float(nfo.replace(',', '.'))),
            4: ('unit', lambda nfo: nfo),
            7: ('warn', lambda nfo: float(nfo.replace(',', '.'))),
            9: ('crit', lambda nfo: float(nfo.replace(',', '.'))),
            11: ('min', lambda nfo: float(nfo.replace(',', '.'))),
            13: ('max', lambda nfo: float(nfo.replace(',', '.'))),
        }

        perf_data_array = []
        perf_data = {}
        i = 0

        for info in perfs:
            if info and info != '':
                if i in pfmap:
                    key, getval = pfmap[i]

                    perf_data[key] = getval(info)

            i += 1

            if i == 15:
                # Try to guess unit if not defined
                if not perf_data.get('unit', None):
                    metric = perf_data['metric']

                    if metric[-1] == ']':
                        metric = metric[:-1]

                    parts = metric.split('[', 1)

                    if len(parts) > 1:
                        metric, unit = parts

                        perf_data['metric'] = metric
                        perf_data['unit'] = unit

                # Make sure perfdata entry is valid
                if 'value' in perf_data and 'metric' in perf_data:
                    perf_data_array.append(perf_data)

                perf_data = {}
                i = 0

        self._pf_array = perf_data_array

    @property
    def perf_data_array(self):
        return self._pf_array
