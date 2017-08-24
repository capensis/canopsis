import os
import time

from canopsis.engines.core import TaskHandler, publish
from canopsis.event import forger
from canopsis.context_graph.import_ctx import ImportKey as Keys,\
    ContextGraphImport, Manager
from canopsis.configuration.configurable.decorator import conf_paths,\
    add_category
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.model import Parameter

TASK_CONF = "TASK_IMPORTCTX"
CONTENT = {
    Parameter("thd_warn_min_per_import", parser=int), Parameter(
        "thd_crit_min_per_import", parser=int)
}

MSG_SUCCEED = "Import {0} succeed."
MSG_FAILED = "Import {0} failed."
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


@conf_paths("context_graph/manager.conf")
@add_category(TASK_CONF, content=CONTENT)
class engine(TaskHandler, MiddlewareRegistry):

    etype = "task_importctx"

    E_IMPORT_FAILED = "Error during the import of id {0} : {1}."
    I_IMPORT_DONE = "Import {0} done."
    I_START_IMPORT = "Start import {0}."

    def send_perfdata(self, uuid, time, updated, deleted):
        """Send stat about the import through a perfdata event.
        :param uuid: the import uuid
        :type uuid: a string
        :param time: the execution time of the import
        :type time: a float
        :param updated: the number of updated entities during the import
        :type updated: an integer
        :param deleted: the number of deleted entities during the import
        :type deleted: an integer
        """

        if not hasattr(self, "thd_warn_s"):
            values = values = self.conf.get(TASK_CONF)
            self._thd_warn_s = values.get("thd_warn_min_per_import").value * 60

        if not hasattr(self, "thd_crit_s"):
            values = values = self.conf.get(TASK_CONF)
            self._thd_crit_s = values.get("thd_crit_min_per_import").value * 60

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
            source_type="task_importctx/report",
            resource=self.amqp_queue,
            state=state,
            state_type=1,
            output=output,
            perf_data_array=perf_data_array)

        publish(event, self.amqp)

    def handle_task(self, job):
        """Handlt the import.
        :param job: the event.
        :type job: a dict"""

        importer = ContextGraphImport(logger=self.logger)
        report_manager = Manager()

        self.logger.info(job)

        uuid = job[Keys.EVT_IMPORT_UUID]

        self.logger.info("Processing import {0}.".format(uuid))

        report_manager.update_status(uuid, {Keys.F_STATUS: Keys.ST_ONGOING})

        start = time.time()
        report = {}
        updated = 0
        deleted = 0

        try:
            updated, deleted = importer.import_context(uuid)

        except Exception as e:
            report = {Keys.F_STATUS: Keys.ST_FAILED, Keys.F_INFO: repr(e)}

            self.logger.error(self.E_IMPORT_FAILED.format(uuid, repr(e)))
            self.logger.exception(e)
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
        report_manager.update_status(uuid, report)

        self.send_perfdata(uuid, delta, updated, deleted)

        try:
            os.remove(Keys.IMPORT_FILE.format(uuid))
        except Exception as e:
            pass

        del importer

        return (state, msg)
