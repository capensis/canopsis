import unittest
import clogging
import os


LOGGER_NAME = 'LOGGER_NAME'


class LoggingTest(unittest.TestCase):

    def setUp(self):
        pass

    def _testLogger(self, logger, name):

        self.assertEquals(logger.name, name)

        old_size = os.path.getsize(logger.handler.baseFilename)

        logger.info('here is an information message')

        new_size = os.path.getsize(logger.handler.baseFilename)

        self.assertNotEquals(old_size, new_size)

        old_size = new_size

        self.assertEquals(logger.handler.formatter._fmt, clogging.INFO_FORMAT)

        logger.debug('here is a debug message which shouldn\'t be displayed')

        new_size = os.path.getsize(logger.handler.baseFilename)

        self.assertEquals(old_size, new_size)

        old_size = new_size

        self.assertEquals(logger.handler.formatter._fmt, clogging.INFO_FORMAT)

        logger.setLevel("DEBUG")

        logger.debug('here is a debug message which should be displayed')

        new_size = os.path.getsize(logger.handler.baseFilename)

        self.assertNotEquals(old_size, new_size)

        old_size = new_size

        self.assertEquals(logger.handler.formatter._fmt, clogging.DEBUG_FORMAT)

        logger.info('here is an information message which should be displayed')

        new_size = os.path.getsize(logger.handler.baseFilename)

        self.assertNotEquals(old_size, new_size)

        old_size = new_size

        self.assertEquals(logger.handler.formatter._fmt, clogging.INFO_FORMAT)

    def testLogger(self):

        logger = clogging.getLogger()

        self._testLogger(logger, "clogging-Myunittest")

        logger = clogging.getLogger("plop")

        self._testLogger(logger, "plop")

        pass

if __name__ == '__main__':
    unittest.main()
