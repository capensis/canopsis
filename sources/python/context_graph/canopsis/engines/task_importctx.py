from canopsis.engines.core import TaskHandler

class engine(TaskHandler):

    etype = "task_importctx"

    def handle_task(self, job):
        self.logger.error("I am here")
        self.logger.info(job)

        return (0, "OK")
