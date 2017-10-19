#!/usr/bin/env python
# -*- coding: utf-8  -*-

"""
Series halepers for influxdb interactions.
"""

from __future__ import unicode_literals

from influxdb import SeriesHelper


class SystemSeriesHelper(SeriesHelper):
    class Meta:
        # The series name must be a string. Add dependent fields/tags in curly brackets.
        #series_name = 'events.stats.{serie_name}'
        series_name = 'system'
        # Defines all the fields in this time series.
        fields = ['min']
        # Defines all the tags for the series.
        tags = ['serie_name']


        #bulk_size = 5
        #autocommit = True

# measurements: _key, component, connector, connector_name, eid, min, resource, type
# series: name, ack_size, alarm_opened_count, alerts_by_host, alerts_count, buffered, cache_size, cached, cps_evt_per_sec, cps_sec_per_evt, cps_state, delay, disk_merged-read, disk_merged-write, disk_octets-read, disk_octets-write, disk_ops-read, disk_ops-write, disk_time-read, disk_time-write, drop_event, entities_size, events_log_size, events_size, files_size, free, idle, if_errors-rx, if_errors-tx, if_octets-rx, if_octets-tx, if_packets-rx, if_packets-tx, in, interrupt, load-longterm, load-midterm, load-shortterm, nice, object_size, out, pass_event, perfdata2_daily_size, perfdata2_size, perfdata3_size, reserved, session_duration, size, softirq, steal, system, used, user, wait
