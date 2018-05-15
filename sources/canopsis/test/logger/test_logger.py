from __future__ import unicode_literals

import logging
import logging.handlers
import os
import tempfile
import unittest

try:
    from io import StringIO
except ImportError:
    from StringIO import StringIO
from canopsis.common import root_path
from canopsis.logger import Logger, OutputFile, OutputStream
import xmlrunner


class TestLogger(unittest.TestCase):

    def test_logger_stream(self):

        output = StringIO()
        logger = Logger.get('stream', output, OutputStream)

        logger.info('first_line')
        logger.info('second_line')

        output.seek(0)
        log_lines = output.readlines()

        self.assertEqual(len(log_lines), 2)

    def test_logger_string_level(self):

        output = StringIO()
        logger = Logger.get('ascii', output, OutputStream, level='info')

        logger.info('fline')

        output.seek(0)
        log_lines = output.readlines()

        self.assertEqual(len(log_lines), 1)

    def test_logger_unicode_level(self):

        output = StringIO()
        logger = Logger.get('unicode', output, OutputStream, level='info')

        logger.info('fline')

        output.seek(0)
        log_lines = output.readlines()

        self.assertEqual(len(log_lines), 1)

    def test_logger_dedup(self):

        output = StringIO()

        log1 = Logger.get('dedup', output, OutputStream)
        log2 = Logger.get('dedup', output, OutputStream)

        log1.info('first_line')
        log2.info('second_line')
        log1.info('third_line')
        log2.info('fourth_line')

        output.seek(0)
        log_lines = output.readlines()

        self.assertEqual(len(log_lines), 4)

    def test_logger_file(self):

        tmpf = tempfile.NamedTemporaryFile(delete=False)
        tmpf.close()

        logger = Logger.get('file', tmpf.name, OutputFile)

        logger.info('first_line')

        with open(tmpf.name, 'r') as fh:
            content = fh.readlines()
            self.assertEqual(len(content), 1)
            self.assertTrue(content[0][:-1].endswith('first_line'))

        os.unlink(tmpf.name)

    def test_logger_memory(self):

        output = StringIO()

        logger = Logger.get('memory', output, OutputStream,
                            level=logging.INFO,
                            memory=True,
                            memory_capacity=2,
                            memory_flushlevel=logging.CRITICAL)

        # under max capacity
        logger.info('1')
        output.seek(0)
        lines = len(output.readlines())
        self.assertEqual(lines, 0)

        # max capacity
        logger.info('2')
        output.seek(0)
        lines = len(output.readlines())
        self.assertEqual(lines, 2)

        # under max capacity but flush level
        logger.critical('3')
        output.seek(0)
        lines = len(output.readlines())
        self.assertEqual(lines, 3)

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
