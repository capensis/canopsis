import os
import time

from canopsis.engines.core import TaskHandler, publish
from canopsis.event import forger
from canopsis.context_graph.import_ctx import ImportKey as Keys,\
    ContextGraphImport, Manager

MSG_SUCCEED = "Import {0} succeed."
MSG_FAILED = "Import {0} failed."
ST_OK = 0
ST_KO = 3

perf_data_array = [{
    'metric': 'importctx/execution_time',
    'value': None,
    'unit': 's'
}, {
    'metric': 'importctx/ent_updated',
    'value': None,
    'unit': 'ent'
}, {
    'metric': 'importctx/ent_deleted',
    'value': None,
    'unit': 'ent'
}]

# event =


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
    I_IMPORT_DONE = "Import {0} done."
    I_START_IMPORT = "Start import {0}."

    def send_perfdata(self, uuid, time, updated, deleted, state):

        perf_data_array[0]["value"] = time
        perf_data_array[1]["value"] = updated
        perf_data_array[2]["value"] = deleted

        output = "execution {0} sec, updated ent"\
                 " {1}, deleted ent {2}".format(time, updated, deleted)

        event = forger(
            connector="Engine",
            connector_name="engine",
            event_type="check",
            source_type=Keys.JOB_ID.format(uuid),
            resource=self.amqp_queue,
            state=state,
            state_type=1,
            output=output,
            perf_data_array=perf_data_array)

        publish(event, self.amqp)

    def handle_task(self, job):
        importer = ContextGraphImport(logger=self.logger)
        report_manager = Manager()

        self.logger.info(job)

        uuid = job[Keys.EVT_IMPORT_UUID]

        self.logger.info("Processing import {0}.".format(uuid))

        report_manager.update_status(uuid, {Keys.F_STATUS: Keys.ST_ONGOING})

        start = time.time()
        report = {}

        try:
            updated, deleted = importer.import_context(uuid)

        except Exception as e:
            report = {Keys.F_STATUS: Keys.ST_FAILED, Keys.F_INFO: repr(e)}

            self.logger.error(self.E_IMPORT_FAILED.format(uuid, repr(e)))
            self.logger.exception(e)
            msg = MSG_FAILED.format(uuid)
            state = ST_KO

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
            state = ST_OK

        end = time.time()
        delta = end - start
        report[Keys.F_EXECTIME] = human_exec_time(delta)
        report_manager.update_status(uuid, report)

        self.send_perfdata(uuid, delta, updated, deleted, 0)

        try:
            os.remove(Keys.IMPORT_FILE.format(uuid))
        except:
            pass


        del importer

        return (state, msg)
