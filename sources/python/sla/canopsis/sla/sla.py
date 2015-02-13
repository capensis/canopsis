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

__version__ = '0.1'


from time import time
from logging import getLogger
from canopsis.old.event import forger
from json import dumps, loads
import datetime
import pprint
pp = pprint.PrettyPrinter(indent=2)


class Sla(object):

    """
        parameters:
        :param storage: topology ids from where find elts.
        :type storage: canopsis storage
        :param rk: the rk to build sla from in the event_log collection
        :type rk: string
        :param template: the sla event template to fill with alert computation
        :type template: string
        :param timewindow: sla timewindow to consider for given rk in seconds
        :type timewindow: int
        :param sla_warning: thresholds beyond witch the alerts percent
            pass the event in warning state
        :type sla_warning: int
        :param sla_critical: thresholds beyond witch the alerts percent
            pass the event in critical state
        :type sla_critical: int
        :param alert_level: defining what is minimum level to consider state
            in alert total time.
        :type alert_level: string one value between minor, major or critical
        :param display_name: used as the sla component name
        :type display_name: string
        :param logger: a logger instance to where the module can write output
        :type logger: logger
    """
    def __init__(
        self,
        storage,
        rk,
        template,
        timewindow,
        sla_warning,
        sla_critical,
        alert_level,
        display_name,
        logger=None
    ):

        self.storage = storage

        self.type = 'sla'

        if logger:
            self.logger = logger
        else:
            self.logger = getLogger('Sla')

        self.states = [0, 1, 2, 3]

        now = time()

        timewindow_date_start = now - timewindow
        self.logger.debug('Timewindow is {}, timestamp is {}'.format(
            timewindow,
            timewindow_date_start
        ))

        self.logger.debug('Computing sla for selector {}'.format(rk))
        # Retrieve sla information from selector record
        sla_information = self.get_sla_information(
            rk,
            timewindow_date_start,
            now
        )

        self.logger.debug('Sla length information is {}'.format(
            len(sla_information)
        ))

        # Compute effective sla dict to be able to fill the ouput template
        sla_measures, first_timestamp = self.compute_sla(sla_information, now)
        self.logger.debug('Sla measures is {}'.format(sla_measures))

        self.logger.debug('Alert level is {}'.format(alert_level))
        # Compute alerts precent depending on user algorithm
        alerts_percent = self.get_alert_percent(sla_measures, alert_level)

        # Compute template from sla measures
        output = self.compute_output(
            template,
            sla_measures,
            alerts_percent,
            first_timestamp
        )

        state = self.compute_state(
            alerts_percent,
            sla_warning,
            sla_critical
        )

        self.logger.debug('Sla computed state is {}'.format(state))

        self.logger.debug('thresholds : warning {}, critical {}'.format(
            sla_warning,
            sla_critical
        ))

        self.event = self.prepare_event(
            display_name,
            sla_measures,
            output,
            state,
            alerts_percent
        )

    def get_alert_percent(self, sla_measures, alert_level):

        # alert_level should never be something else than minor,major,critical

        if alert_level == 'minor':
            return sla_measures[1] + sla_measures[2] + sla_measures[3]

        if alert_level == 'major':
            return sla_measures[2] + sla_measures[3]

        if alert_level == 'critical':
            return sla_measures[3]

    def compute_state(self, alerts_percent, warning, critical):

        CRITICAL = 3
        MINOR = 1
        INFO = 0

        availability = 1.0 - alerts_percent

        self.logger.debug('availability {} warning {}, critical {}'.format(
            availability,
            warning,
            critical
        ))

        if availability < float(critical) / 100.0:
            return CRITICAL

        if availability < float(warning) / 100.0:
            return MINOR

        return INFO

    def get_event(self):
        # may have any modifiers here
        return self.event

    def get_sla_information(
        self,
        selector_rk,
        timewindow_date_start,
        now
    ):

        """Sla information is a list containing all state in the timewindow for
        current sla event."""

        sla = []

        events_log = self.storage.get_backend('events_log')

        # Fetch previous state
        state_before = events_log.find_one({
            'rk': selector_rk,
            'timestamp': {'$lte': timewindow_date_start}
        }, {
            'state': 1,
            'timestamp': 1,
            '_id': 0
        }, sort=[('timestamp', -1)])

        self.logger.debug('state_before {}'.format(state_before))

        if state_before:
            self.logger.debug('State before found ! {}'.format(state_before))
            state_before['timestamp'] = timewindow_date_start
            sla.append(state_before)

        # Fetch all state between before and now
        sla_infos = events_log.find({
            'rk': selector_rk,
            'timestamp': {'$gte': timewindow_date_start}
        }, {
            'state': 1,
            'timestamp': 1,
            '_id': 0
        })

        for sla_info in sla_infos:
            sla.append(sla_info)

        # Add last delta time because state may remain until now
        if len(sla) and sla[-1]['timestamp'] < now:
            delta_time = now - sla[-1]['timestamp']
            sla.append({
                'timestamp': now,
                'state': sla[-1]['state']
            })
            self.logger.debug('Add time until now for last state')

        self.logger.debug('Sla information from events_log : {}'.format(
            pp.pformat(sla)
        ))

        return sla

    def compute_sla(self, sla_information, now):

        """
        Allow computing percents time portion where the
        selector state were. sla_information is a list of
        state and timestamp dict representing selector state
        evolution
        """

        sla_times = {
            0: 0.0,
            1: 0.0,
            2: 0.0,
            3: 0.0,
        }
        sla_measures = {
            0: 0.0,
            1: 0.0,
            2: 0.0,
            3: 0.0,
        }

        # Compute duration between eache state change
        # default value for first timestamp
        first_timestamp = now
        if len(sla_information):

            # Allow computing the percentage inside the timewindow
            first_timestamp = date_start = sla_information[0]['timestamp']
            previous_state = sla_information[0]['state']

            self.logger.debug('Compute since {}, state were {}'.format(
                date_start,
                previous_state
            ))

            # compute what proportion of time the event
            # remained in the same state
            for sla_info in sla_information:
                delta_time = sla_info['timestamp'] - date_start
                date_start = sla_info['timestamp']
                sla_times[previous_state] += delta_time
                previous_state = sla_info['state']

                self.logger.debug('Add time {} to state {}'.format(
                    delta_time,
                    sla_info['state'],
                ))

            self.logger.debug('Computed sla times are {}'.format(
                sla_times
            ))

            total_time = now - first_timestamp
            self.logger.debug('total_time {}, now {}, date_start {}'.format(
                total_time,
                now,
                first_timestamp
            ))

            if total_time == 0:
                # Avoids divid by 0 error
                self.logger.warning('Tried to divide by 0 in compute sla')
                total_time = 1

            for state in sla_times:
                percent = float(sla_times[state]) / float(total_time)
                sla_measures[state] = percent

        return sla_measures, first_timestamp

    def compute_output(
        self,
        template,
        sla_measures,
        alerts_percent,
        first_timestamp
    ):

        def to_percent(value):
            value = value * 100
            return ("%0.2f" % value)

        # Timestamp to date string
        TSTART = datetime.datetime.fromtimestamp(
            first_timestamp
        ).strftime('%Y-%m-%d %H:%M:%S')

        if isinstance(template, basestring):
            output = template.replace('[OFF]', to_percent(sla_measures[0]))
            output = output.replace('[MINOR]', to_percent(sla_measures[1]))
            output = output.replace('[MAJOR]', to_percent(sla_measures[2]))
            output = output.replace('[CRITICAL]', to_percent(sla_measures[3]))
            output = output.replace('[ALERTS]', to_percent(alerts_percent))
            output = output.replace('[TSTART]', TSTART)
            self.logger.info('SLA computed output is : {}'.format(output))
            return output
        else:
            self.logger.warning('Sla template is not a string, nothing done.')
            return ''

    def prepare_event(
        self,
        display_name,
        sla_measures,
        output,
        sla_state,
        alerts_percent
    ):
        perf_data_array = []

        # Compute metrics to publish
        for state in self.states:

            perf_data_array.append({
                'metric': 'cps_pct_by_{}'.format(state),
                'value': round(sla_measures[state] * 100.0, 2),
                'max': 100
            })

        availability = (1.0 - alerts_percent) * 100.0
        perf_data_array.append({
            'metric': 'cps_avail',
            'value': round(availability, 2),
            'max': 100
        })

        event = forger(
            connector='sla',
            connector_name='engine',
            event_type='sla',
            source_type='resource',
            component=display_name,
            resource='sla',
            state=sla_state,
            output=output,
            perf_data_array=perf_data_array,
            display_name=display_name
        )

        self.logger.info('publishing sla {}, states {}'.format(
            display_name,
            sla_measures
        ))

        self.logger.debug('event : {}'.format(pp.pformat(event)))

        return event
