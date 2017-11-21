# -*- coding: utf-8 -*-

"""
VENDOR https://github.com/pklaus/bottlelog

ORIGIN LICENSE:

Copyright (c) 2013 Philipp Klaus.  All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1.  Redistributions of source code must retain the above copyright notice,
this list of conditions and the following disclaimer.
2.  Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.
3.  Neither the name of the Portable Site Information Project nor the names
of its contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.


    Source: http://docs.python.org/3/library/datetime.html → "Example tzinfo classes"
    Idea:   http://stackoverflow.com/a/2071364/183995
"""

from __future__ import unicode_literals

import time
from datetime import datetime as dt, tzinfo, timedelta
from bottle import request, response
from bottle import HTTPResponse

import time as _time

ZERO = timedelta(0)
STDOFFSET = timedelta(seconds = -_time.timezone)
if _time.daylight:
    DSTOFFSET = timedelta(seconds = -_time.altzone)
else:
    DSTOFFSET = STDOFFSET

DSTDIFF = DSTOFFSET - STDOFFSET

class LocalTimezone(tzinfo):

    def utcoffset(self, dt):
        if self._isdst(dt):
            return DSTOFFSET
        else:
            return STDOFFSET

    def dst(self, dt):
        if self._isdst(dt):
            return DSTDIFF
        else:
            return ZERO

    def tzname(self, dt):
        return _time.tzname[self._isdst(dt)]

    def _isdst(self, dt):
        tt = (dt.year, dt.month, dt.day,
              dt.hour, dt.minute, dt.second,
              dt.weekday(), 0, 0)
        stamp = _time.mktime(tt)
        tt = _time.localtime(stamp)
        return tt.tm_isdst > 0

Local = LocalTimezone()

# -*- coding: utf-8 -*-

"""
This is a plugin for Bottle web applications.
When you install it to an application, it will
log all requests to your site. It's even imitating
Apache's combined log format to allow you to use
any of the many tools for Apache log file analysis.

Homepage: https://github.com/pklaus/bottlelog

Copyright (c) 2013, Philipp Klaus.
License: BSD (see LICENSE for details)
"""

__all__ = ['LoggingPlugin']

def format_NCSA_log(request, response, bodylen):
    """
      Apache log format 'NCSA extended/combined log':
      "%h %l %u %t \"%r\" %>s %b \"%{Referer}i\" \"%{User-agent}i\""
      see http://httpd.apache.org/docs/current/mod/mod_log_config.html#logformat
    """

    # Let's collect log values
    val = dict()
    val['host'] = request.remote_addr or request.remote_route[0]
    val['logname'] = '-'
    val['user'] = '-'
    val['time'] = dt.now(tz=Local).strftime("%d/%b/%Y:%H:%M:%S %z")
    val['request'] = "{} {} {}".format(
          request.method,
          request.path,
          request.environ.get('SERVER_PROTOCOL', '')
        )
    val['status'] = response.status_code
    val['size'] = bodylen
    val['referer'] = request.get_header('Referer','')
    val['agent'] = request.get_header('User-agent','')

    # see http://docs.python.org/3/library/string.html#format-string-syntax
    FORMAT = '{host} {logname} {user} [{time}] "{request}" '
    FORMAT += '{status} {size} "{referer}" "{agent}"'
    return FORMAT.format(**val)


def format_with_response_time(*args, **kwargs):
    """
      This is the format for TinyLogAnalyzer:
      https://pypi.python.org/pypi/TinyLogAnalyzer
    """
    rt_ms = kwargs.get('rt_ms', 0)
    return format_NCSA_log(*args) + " {}/{}".format(int(rt_ms/1000000), rt_ms)


class LoggingPlugin(object):
    ''' This is the Bottle logging plugin itself. '''

    name = 'logging'
    api = 2

    def __init__(self, logger):
        self.logger = logger

    def __call__(self, callback):
        def wrapper(*args, **kwargs):
            start = time.clock()
            body = callback(*args, **kwargs)
            runtime = int((time.clock() - start) * 10**6)
            bodylen = len(body) if not isinstance(body, HTTPResponse) else 0
            #msg = format_NCSA_log(request, response, bodylen)
            msg = format_with_response_time(request, response, bodylen, rt_ms=runtime)
            self.logger.info(msg)
            return body
        return wrapper

