import unittest

try:
    from io import StringIO
except ImportError:
    from StringIO import StringIO

from canopsis.logger import Logger, OutputFile, OutputStream

class TestLogger(unittest.TestCase):

    def test_logger(self):

        output = StringIO()
        logger = Logger.get('single', output, OutputStream)

        logger.info(u'first_line')
        logger.info(u'second_line')

        output.seek(0)
        log_lines = output.readlines()

        self.assertEqual(len(log_lines), 2)

    def test_logger_dedup(self):

        output = StringIO()

        log1 = Logger.get('dedup', output, OutputStream)
        log2 = Logger.get('dedup', output, OutputStream)

        log1.info(u'first_line')
        log2.info(u'second_line')
        log1.info(u'third_line')
        log2.info(u'fourth_line')

        output.seek(0)
        log_lines = output.readlines()

        self.assertEqual(len(log_lines), 4)