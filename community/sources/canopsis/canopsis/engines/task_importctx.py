#!/usr/bin/env python2
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import os
import time

from canopsis.engines.core import TaskHandler
from canopsis.event import forger
from canopsis.context_graph.import_ctx import ImportKey as Keys,\
    ContextGraphImport, Manager
from canopsis.confng import Configuration, Ini
from canopsis.logger import Logger, OutputFile

MSG_SUCCEED = "Import {} succeed."
MSG_FAILED = "Import {} failed."
ST_INFO = 0
ST_MINOR = 1
ST_WARNING = 2
ST_CRITICAL = 3

perf_data_array = [{
    'metric': 'execution_time',
    'value': None,
    'unit': 'GAUGE'
}, {
    'metric': 'ent_updated',
    'value': None,
    'unit': 'GAUGE'
}, {
    'metric': 'ent_deleted',
    'value': None,
    'unit': 'GAUGE'
}]


def human_exec_time(exec_time):
    """Return from time a human readable string that represent the execution
    time in a human readable format."""
    exec_time = int(exec_time)  # we do not care of everything under the second
    hours = exec_time / 3600
    minutes = (exec_time - 3600 * hours) / 60
    seconds = exec_time - (hours * 3600) - (minutes * 60)
    return "{0}:{1}:{2}".format(
        str(hours).zfill(2), str(minutes).zfill(2), str(seconds).zfill(2))


class engine(TaskHandler):

    etype = "task_importctx"

    E_IMPORT_FAILED = "Error during the import of id {0} : {1}."
    I_IMPORT_DONE = "Import {} done."
    I_START_IMPORT = "Start import {}."
    LOG_NAME = "Task_importctx"
    LOG_PATH = "var/log/engines/task_importctx.log"

    TASK_CONF = "TASK_IMPORTCTX"
    CONF_PATH = "etc/context_graph/manager.conf"
    THD_WARN_S = "thd_warn_min_per_import"
    THD_CRIT_S = "thd_crit_min_per_import"

    def __init__(self, config=None, *args, **kwargs):

        super(engine, self).__init__(*args, **kwargs)

        if config is None:
            config = Configuration.load(self.CONF_PATH, Ini)

        section = config.get(self.TASK_CONF)

        self._thd_warn_s = section.get(self.THD_WARN_S) * 60
        self._thd_crit_s = section.get(self.THD_CRIT_S) * 60

        self.logger = Logger.get(self.LOG_NAME,
                                 self.LOG_PATH,
                                 output_cls=OutputFile)
        # self.importer = ContextGraphImport(logger=self.logger)
        self.report_manager = Manager()

    def send_perfdata(self, uuid, time, updated, deleted):
        """Send stat about the import through a perfdata event.

        :param str uuid: the import uuid
        :param float time: the execution time of the import
        :param int updated: the number of updated entities during the import
        :param int deleted: the number of deleted entities during the import
        """
        # define the state according to the duration of the import
        if time > self._thd_crit_s:
            state = ST_WARNING
        elif time > self._thd_warn_s:
            state = ST_MINOR
        else:
            state = ST_INFO

        perf_data_array[0]["value"] = time
        perf_data_array[1]["value"] = updated
        perf_data_array[2]["value"] = deleted

        output = "execution : {0} sec, updated ent :"\
                 " {1}, deleted ent : {2}".format(time, updated, deleted)

        self.logger.critical("AMQP queue = {0}".format(self.amqp_queue))

        # create a perfdata event
        event = forger(
            connector="Taskhandler",
            connector_name=self.etype,
            component=uuid,
            event_type="check",
            source_type="resource",
            resource="task_importctx/report",
            state=state,
            state_type=1,
            output=output,
            perf_data_array=perf_data_array
        )

        try:
            self.work_amqp_publisher.canopsis_event(event)
        except Exception as e:
            self.logger.exception("Unable to send event")

    def handle_task(self, job):
        """
        Handle the import.

        :param dict job: the event
        :returns:
        """
        self.logger.info(job)

        uuid = job[Keys.EVT_IMPORT_UUID]

        self.logger.info("Processing import {0}.".format(uuid))

        self.report_manager.update_status(
            uuid, {Keys.F_STATUS: Keys.ST_ONGOING})

        start = time.time()
        report = {}
        updated = 0
        deleted = 0

        importer = ContextGraphImport(logger=self.logger)

        try:
            updated, deleted = importer.import_context(uuid)

        except Exception as ex:
            report = {Keys.F_STATUS: Keys.ST_FAILED, Keys.F_INFO: repr(ex)}

            self.logger.error(self.E_IMPORT_FAILED.format(uuid, repr(ex)))
            self.logger.exception(ex)
            msg = MSG_FAILED.format(uuid)
            state = ST_CRITICAL

        else:
            report = {
                Keys.F_STATUS: Keys.ST_DONE,
                Keys.F_STATS: {
                    Keys.F_DELETED: deleted,
                    Keys.F_UPDATED: updated
                }
            }
            self.logger.info(self.I_IMPORT_DONE.format(uuid))
            msg = MSG_SUCCEED.format(uuid)
            state = ST_INFO

        end = time.time()
        delta = end - start
        report[Keys.F_EXECTIME] = human_exec_time(delta)
        self.report_manager.update_status(uuid, report)

        self.send_perfdata(uuid, delta, updated, deleted)
        del importer

        try:
            os.remove(Keys.IMPORT_FILE.format(uuid))
        except Exception as ex:
            pass

        return (state, msg)
