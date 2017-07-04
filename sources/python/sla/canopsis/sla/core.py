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

from canopsis.common.init import basestring
from time import time
from logging import getLogger
from canopsis.event import forger
from datetime import datetime
from pprint import PrettyPrinter
from canopsis.timeserie.timewindow import Period
from canopsis.event.manager import Event

from canopsis.perfdata.manager import SLIDING_TIME

pp = PrettyPrinter(indent=2)


__version__ = '0.1'


class Sla(object):

    """Enable Sla computation based on event's states history

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
        timewindow_dict,
        sla_warning,
        sla_critical,
        alert_level,
        display_name,
        logger=None
    ):

        self.storage = storage
        self.eventmanager = Event()
        self.type = 'sla'

        if logger:
            self.logger = logger
        else:
            self.logger = getLogger('Sla')

        self.states = [0, 1, 2, 3]

        now = time()

        timewindow = timewindow_dict['seconds']

        timewindow_date_start = now - timewindow
        self.logger.debug(u'Timewindow is {}, timestamp is {}'.format(
            timewindow,
            timewindow_date_start
        ))

        self.logger.debug(u'Computing sla for selector {}'.format(rk))
        # Retrieve sla information from selector record
        sla_information = self.get_sla_information(
            rk,
            timewindow_date_start,
            now
        )

        self.logger.debug(u'Sla length information is {}'.format(
            len(sla_information)
        ))

        # Compute effective sla dict to be able to fill the ouput template
        sla_measures, sla_times = self.compute_sla(
            sla_information,
            now
        )
        self.logger.debug(u'Sla measures is {}'.format(sla_measures))
        self.logger.debug(u'Sla times is {}'.format(sla_times))

        self.logger.debug(u'Alert level is {}'.format(alert_level))
        # Compute alerts precent depending on user algorithm
        alerts_percent, alerts_duration = \
            self.get_alert_percent(
                sla_measures,
                sla_times,
                alert_level
            )

        self.logger.debug(u'Alert percent : {} , Alert duration : {}'.format(
            alerts_percent,
            alerts_duration
        ))

        avail_duration = timewindow - alerts_duration

        # Compute template from sla measures
        output = self.compute_output(
            template,
            sla_measures,
            alerts_percent,
            alerts_duration,
            avail_duration,
            timewindow_date_start,
        )

        state = self.compute_state(
            alerts_percent,
            sla_warning,
            sla_critical
        )

        self.logger.debug(u'Sla computed state is {}'.format(state))

        self.logger.debug(u'thresholds : warning {}, critical {}'.format(
            sla_warning,
            sla_critical
        ))

        self.event = self.prepare_event(
            display_name,
            sla_measures,
            output,
            state,
            alerts_percent,
            alerts_duration,
            avail_duration,
            timewindow_dict,
            now
        )

    def get_alert_percent(self, sla_measures, sla_times, alert_level):

        # alert_level should never be something else than minor,major,critical
        alerts_percent = 0.0
        alerts_duration = 0.0

        if alert_level == 'minor':
            alerts_percent = (sla_measures[1] + sla_measures[2] +
                              sla_measures[3])

            alerts_duration = sla_times[1] + sla_times[2] + sla_times[3]

        if alert_level == 'major':
            alerts_percent = sla_measures[2] + sla_measures[3]
            alerts_duration = sla_times[2] + sla_times[3]

        if alert_level == 'critical':
            alerts_percent = sla_measures[3]
            alerts_duration = sla_times[3]

        return (alerts_percent, alerts_duration)

    def compute_state(self, alerts_percent, warning, critical):

        CRITICAL = 3
        MINOR = 1
        INFO = 0

        availability = 1.0 - alerts_percent

        self.logger.debug(u'availability {} warning {}, critical {}'.format(
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

        events_log = self.storage.get_backend('events_log')
        sla_information = []
        projection = {
            'state': 1,
            'timestamp': 1,
            '_id': 0
        }

        # Try to find last state
        previous_event_log = events_log.find_one({
            'rk': selector_rk,
            'timestamp': {'$lt': timewindow_date_start}
        }, projection, sort=[('timestamp', -1)])

        self.logger.debug(u'previous event log {}'.format(previous_event_log))

        # Default value is set as no previous information
        if previous_event_log is None:
            previous_event_log = {
                'state': 0,
                'timestamp': timewindow_date_start
            }
        else:
            previous_event_log['timestamp'] = timewindow_date_start

        sla_information.append(previous_event_log)

        # Fetch all event log data in timewindow
        sla_information += list(events_log.find(
            {
                'rk': selector_rk,
                'timestamp': {'$gt': timewindow_date_start}
            },
            projection,
            sort=[('timestamp', 1)]
        ))

        self.logger.debug(u' # sla information {}'.format(sla_information))

        return sla_information

    def compute_sla(self, sla_information, now):

        """Allow computing percents time portion where the
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

        total_time = float(now - sla_information[0]['timestamp'])
        previous_state = sla_information[0]['state']
        previous_timestamp = sla_information[0]['timestamp']
        duration = 0.0

        for step in sla_information[1:]:
            duration += step['timestamp'] - previous_timestamp
            previous_timestamp = step['timestamp']
            if step['state'] != previous_state:
                sla_times[previous_state] += duration
                previous_state = step['state']
                duration = 0.0

        sla_times[previous_state] += now - previous_timestamp

        sla_measures = {
            0: 0.0,
            1: 0.0,
            2: 0.0,
            3: 0.0,
        }

        for state in sla_times:
            sla_time = float(sla_times[state])
            self.logger.debug(
                u'state {} time {} total {} now {}, start {}'.format(
                    state, sla_time,
                    total_time, now,
                    sla_information[0]['state']
                )
            )
            percent = sla_time / total_time
            sla_measures[state] = percent

        return sla_measures, sla_times

    def compute_output(
        self,
        template,
        sla_measures,
        alerts_percent,
        alerts_duration,
        avail_duration,
        first_timestamp,
    ):

        def to_percent(value):
            value = value * 100
            return ("%0.2f" % value)

        # Timestamp to date string
        TSTART = datetime.fromtimestamp(
            first_timestamp
        ).strftime('%Y-%m-%d %H:%M:%S')

        def duration_to_time(seconds):
            # duration in seconds
            m, s = divmod(seconds, 60)
            h, m = divmod(m, 60)
            return "%dh%02dm%02ds" % (h, m, s)

        if isinstance(template, basestring):
            # Embed sla measures percents in the output
            output = template.replace('[OFF]', to_percent(sla_measures[0]))
            output = output.replace('[MINOR]', to_percent(sla_measures[1]))
            output = output.replace('[MAJOR]', to_percent(sla_measures[2]))
            output = output.replace('[CRITICAL]', to_percent(sla_measures[3]))
            output = output.replace('[ALERTS]', to_percent(alerts_percent))

            # Embed sla measures available total percentage in the output
            output = output.replace('[P_AVAIL]', to_percent(
                1.0 - alerts_percent)
            )

            # Embed sla measures durations in the output
            output = output.replace('[T_AVAIL]', duration_to_time(
                avail_duration)
            )
            output = output.replace('[T_ALERT]', duration_to_time(
                alerts_duration)
            )

            # Embed sla measures first date in the output
            output = output.replace('[TSTART]', TSTART)
            self.logger.info(u'SLA computed output is : {}'.format(output))
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
        alerts_percent,
        alerts_duration,
        avail_duration,
        timewindow_dict,
        now
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
            'max': 100,
            SLIDING_TIME: True
        })
        perf_data_array.append({
            'metric': 'cps_avail_duration',
            'value': avail_duration,
            SLIDING_TIME: True
        })
        perf_data_array.append({
            'metric': 'cps_alerts_duration',
            'value': alerts_duration,
            SLIDING_TIME: True
        })

        period_options = {
            timewindow_dict['durationType']: timewindow_dict['value']
        }
        self.logger.debug(u'period options {}, now {}'.format(
            period_options,
            now
        ))

        period = Period(**period_options)

        periodic_timestamp = period.round_timestamp(now, next_period=True)

        self.logger.debug(u'periodic timestamp {}'.format(periodic_timestamp))

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
            display_name=display_name,
            timestamp=periodic_timestamp
        )

        self.logger.info(u'publishing sla {}, states {}'.format(
            display_name,
            sla_measures
        ))

        self.logger.debug(u'event : {}'.format(pp.pformat(event)))

        return event
