#!/usr/bin/env python
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

from time import time
from logging import getLogger
from canopsis.old.event import forger
from json import dumps, loads
import pprint

pp = pprint.PrettyPrinter(indent=2)


class Sla(object):

    def __init__(
        self,
        selector_record,
        event,
        storage,
        record_id,
        logger=None
    ):

        self.type = 'sla'

        if logger:
            self.logger = logger
        else:
            self.logger = getLogger('Sla')

        self.states = [0, 1, 2, 3]

        # This template should be always set
        template = selector_record.get_sla_output_tpl()

        # Retrieve sla information from selector record
        sla_information = self.get_sla_information(selector_record)
        self.logger.debug('Sla information is {}'.format(sla_information))

        # Timewindow computation duration
        timewindow = selector_record.get_sla_timewindow()
        self.logger.debug('Timewindow is {}'.format(timewindow))

        previous_selector_state = selector_record.get_previous_selector_state()
        prev_state_tw_start = selector_record.get_state_at_timewindow_start()

        current_state = event['state']

        previous_sla_information = sla_information.copy()

        # sla_information is updated in this method
        state_at_timewindow_start = self.update_sla_information(
            timewindow,
            current_state,
            previous_selector_state,
            sla_information
        )

        sla_information_changed = previous_sla_information != sla_information
        # When new sla information computed
        # save new selector record new information
        if current_state != previous_selector_state or sla_information_changed:
            self.update_selector_record(
                current_state,
                previous_selector_state,
                state_at_timewindow_start,
                sla_information,
                storage,
                record_id
            )

        # Compute effective sla dict to be able to fill the ouput template
        sla_measures = self.compute_sla(
            timewindow,
            sla_information,
            previous_selector_state,
            prev_state_tw_start
        )
        self.logger.debug('Sla measures is {}'.format(sla_measures))

        # Compute template from sla measures
        output = self.compute_output(template, sla_measures)

        state = self.compute_state(sla_measures, selector_record)
        self.logger.debug('Sla computed state is {}'.format(state))
        self.logger.debug('thresholds : warning {}, critical {}'.format(
            selector_record.get_sla_warning(),
            selector_record.get_sla_critical()
        ))

        self.event = self.prepare_event(
            selector_record,
            sla_measures,
            output,
            state
        )

    def compute_state(self, sla_measures, selector_record):

        warning = selector_record.get_sla_warning()
        critical = selector_record.get_sla_critical()
        alerts_percent = sla_measures[1] + sla_measures[2] + sla_measures[3]

        CRITICAL = 3
        MINOR = 1
        INFO = 0

        if alerts_percent > float(critical) / 100.0:
            return CRITICAL

        if alerts_percent > float(warning) / 100.0:
            return MINOR

        return INFO

    def update_selector_record(
        self,
        current_state,
        previous_selector_state,
        state_at_timewindow_start,
        sla_information,
        storage,
        record_id
    ):
        self.logger.debug(
            'Change found!, Updating selector record with id {},' +
            ' previous state {} and sla information {}'.format(
                record_id,
                current_state,
                sla_information
            ))

        update = {
            'sla_information': dumps(sla_information)
        }
        if current_state != previous_selector_state:
            update['previous_selector_state'] = current_state

        if state_at_timewindow_start is not None:
            update['state_at_timewindow_start'] = state_at_timewindow_start

        storage.update(record_id, update)

    def get_event(self):
        # may have any modifiers here
        return self.event

    def get_sla_information(self, selector_record):
        """
        Sla information is a dict that contains a list for each possible state.
        These lists are made of dict that contains a start and a stop date and
        looks like {'start': XXX, 'stop': YYY}
        """
        # When serialized, key are converted to string,
        # so on init they are declared as string
        info = {
            '0': [],
            '1': [],
            '2': [],
            '3': [],
        }

        if 'sla_information' in selector_record.data:
            info = selector_record.data['sla_information']
            if isinstance(selector_record.data['sla_information'], basestring):
                info = loads(info)

        return info

    def update_sla_information(
        self,
        timewindow,
        current_state,
        previous_selector_state,
        sla_information
    ):
        """
        This method aims to clean a sla information dict by
        removing entries that are not part of the given timewindow
        It also adds new information to this dict when the current
        selector state changed.
        It modifies the sla information dict given as param and returns
        the state where the selector was at start of timewindow.
        """
        now = time()
        start_date = now - timewindow

        # Allow computing last state whose sla missing time is attribued
        latest_date_clear = 0
        # This state is the one in witch selector was before the timewindow.
        state_at_timewindow_start = None

        # Clear timewindow that are out of sla scope.
        for state in self.states:
            cleaned_state_info = []

            # Keep only information that remains in the current timewindow
            for window in sla_information[str(state)]:
                start = window.get('start', None)
                stop = window.get('stop', None)
                if start >= start_date or stop >= start_date:
                    cleaned_state_info.append(window)
                else:
                    # latest date clear just help retrieving the good state
                    if start > latest_date_clear:
                        latest_date_clear = start
                        state_at_timewindow_start = state
                    # stop date may not exist
                    if stop > latest_date_clear:
                        latest_date_clear = stop
                        state_at_timewindow_start = state

            # New information computed
            sla_information[str(state)] = cleaned_state_info

        # Add timewindow infomation to the sla information when state changed
        if current_state != previous_selector_state:

            # Append a new window for current_state
            sla_information[str(current_state)].append({
                'start': now
            })

            # Ends a timewindow for previous state information
            # when previsous state exists
            if (previous_selector_state is not None and
                    len(sla_information[str(previous_selector_state)])):
                sla_information[str(previous_selector_state)][-1]['stop'] = now

        self.logger.debug('computed state_at_timewindow_start is {}'.format(
            state_at_timewindow_start
        ))
        return state_at_timewindow_start

    def compute_sla(
        self,
        timewindow,
        sla_information,
        previous_state,
        prev_state_tw_start
    ):

        """
        From a sla informatio dict, a new dict is computed and looks like
        {state: percent} where state can be one of [0, 1, 2, 3] values and
        percent value is computed depending on the time a selector remained
        in a state or another for the timewindow whole duration.
        """
        results = {}

        # Lowest date allow compute sla difference between the first
        # state change and current date for sla to get at end sla
        # sum equal to 100%. This algorithm consider the first sla
        # range as ok state.
        now = lowest_date = time()
        timewindow_date_start = now - timewindow

        previous_state_missing_time = 0

        # Clear timewindow that are out of sla scope.
        for state in self.states:

            total_duration = 0.0
            # Keep only information that remains in the current timewindow
            for window in sla_information[str(state)]:

                if window['start'] < lowest_date:
                    lowest_date = window['start']

                # case where sla information exists but
                # starts before timewindow
                if window['start'] < timewindow_date_start:
                    start = timewindow_date_start
                else:
                    start = window['start']

                # Stop can be not already defined ,
                # and then computation is done until now
                if 'stop' not in window:
                    # stop is set same as start then difference equals 0
                    stop = now
                else:
                    stop = window['stop']

                # Check duration is positive value
                duration = float(stop) - float(start)
                if duration <= 0:
                    self.logger.error('Sla error when computing duration.')

                total_duration += duration

            # Avoids division by 0, may introduce precision errors
            if timewindow == 0:
                timewindow = 1
                self.logger.warning(
                    'timewindow for sla computation is 0,' +
                    ' this may not be normal'
                )

            # Keys are not string
            results[state] = float(total_duration) / float(timewindow)

        # When state at start of time window is not defined,
        # then take the previous state instead
        if prev_state_tw_start is None:

            # Consider that missing time is for previous state.
            if previous_state is None:
                previous_state = 0
            missing_time_state_target = previous_state
        else:
            missing_time_state_target = prev_state_tw_start

        # Add difference between first date and now - timewindow
        missing_time = lowest_date - timewindow_date_start
        missing_time_percent = float(missing_time) / float(timewindow)

        self.logger.debug('missing time is {}, represents {} %'.format(
            missing_time,
            missing_time_percent * 100
        ))

        results[missing_time_state_target] += missing_time_percent

        return results

    def compute_output(self, template, sla_measures):

        def to_percent(value):
            value = value * 100
            return ("%0.2f" % value)

        if isinstance(template, basestring):
            output = template.replace('[OFF]', to_percent(sla_measures[0]))
            output = output.replace('[MINOR]', to_percent(sla_measures[1]))
            output = output.replace('[MAJOR]', to_percent(sla_measures[2]))
            output = output.replace('[CRITICAL]', to_percent(sla_measures[3]))
            output = output.replace(
                '[ALERTS]',
                to_percent(sla_measures[1] + sla_measures[2] + sla_measures[3])
            )
            self.logger.debug('SLA computed output is : {}'.format(output))
            return output
        else:
            self.logger.warning('Sla template is not a string, nothing done.')
            return ''

    def prepare_event(self, selector_record, sla_measures, output, sla_state):

        perf_data_array = []

        # Compute metrics to publish
        for state in self.states:

            state_name = {
                0: 'off',
                1: 'minor',
                2: 'major',
                3: 'critical',
            }[state]

            perf_data_array.append({
                'metric': 'cps_sla_{}'.format(state_name),
                'value': sla_measures[state]
            })

        event = forger(
            connector="sla",
            connector_name="engine",
            event_type="sla",
            source_type="resource",
            component=selector_record.display_name,
            resource='sla',
            state=sla_state,
            output=output,
            perf_data_array=perf_data_array,
            display_name=selector_record.display_name
        )

        self.logger.info('publishing sla {}, states {}'.format(
            selector_record.display_name,
            sla_measures
        ))
        self.logger.debug('event: {}'.format(event))

        return event

