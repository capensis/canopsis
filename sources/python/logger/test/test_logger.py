import unittest
import tempfile
import os
import logging
import logging.handlers

try:
    from io import StringIO
except ImportError:
    from StringIO import StringIO

from canopsis.logger import Logger, OutputFile, OutputStream

class TestLogger(unittest.TestCase):

    def test_logger_stream(self):

        output = StringIO()
        logger = Logger.get('stream', output, OutputStream)

        logger.info(u'first_line')
        logger.info(u'second_line')

        output.seek(0)
        log_lines = output.readlines()

        self.assertEqual(len(log_lines), 2)

    def test_logger_string_level(self):

        output = StringIO()
        logger = Logger.get('ascii', output, OutputStream, level='info')

        logger.info(u'fline')

        output.seek(0)
        log_lines = output.readlines()

        self.assertEqual(len(log_lines), 1)

    def test_logger_unicode_level(self):

        output = StringIO()
        logger = Logger.get(u'unicode', output, OutputStream, level='info')

        logger.info(u'fline')

        output.seek(0)
        log_lines = output.readlines()

        self.assertEqual(len(log_lines), 1)

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

    def test_logger_file(self):

        tmpf = tempfile.NamedTemporaryFile(delete=False)
        tmpf.close()

        logger = Logger.get('file', tmpf.name, OutputFile)

        logger.info(u'first_line')

        with open(tmpf.name, 'r') as fh:
            content = fh.readlines()
            self.assertEqual(len(content), 1)
            self.assertTrue(content[0][:-1].endswith(u'first_line'))

        os.unlink(tmpf.name)

    def test_logger_memory(self):

        output = StringIO()

        logger = Logger.get('memory', output, OutputStream,
                    level=logging.INFO,
                    memory=True, memory_capacity=2, memory_flushlevel=logging.CRITICAL)

        # under max capacity
        logger.info(u'1')
        output.seek(0)
        lines = len(output.readlines())
        self.assertEqual(lines, 0)

        # max capacity
        logger.info(u'2')
        output.seek(0)
        lines = len(output.readlines())
        self.assertEqual(lines, 2)

        # under max capacity but flush level
        logger.critical(u'3')
        output.seek(0)
        lines = len(output.readlines())
        self.assertEqual(lines, 3)
