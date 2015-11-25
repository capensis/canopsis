# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from canopsis.common.init import basestring
from canopsis.task.core import get_task
from canopsis.engines.core import Engine
from canopsis.configuration.configurable import Configurable
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.model import Parameter

CONF_PATH = 'engines/engines.conf'  #: dynamic engine configuration path
CATEGORY = 'ENGINE'  #: dynamic engine configuration category


@conf_paths(CONF_PATH)
class engine(Engine, Configurable):
    """
    Engine which is able to load dynamically its event processing through
    configuration properties.

    :var str event_processing: event processing event_processing path.
    :var dict params: event processing event_processing parameters.
    """

    EVENT_PROCESSING = 'event_processing'  #: event_processing field name
    BEAT_PROCESSING = 'beat_processing'  #: beat_processing field name
    PARAMS = 'params'  #: event processing params field name

    NEXT_AMQP_QUEUES = 'next_amqp_queues'  #: next amqp queues
    NEXT_BALANCED = 'next_balanced'  #: next balanced
    NAME = 'name'  #: self name
    BEAT_INTERVAL = 'beat_interval'  #: beat interval
    EXCHANGE_NAME = 'exchange_name'  #: exchange name
    ROUTING_KEYS = 'routing_keys'  #: routing keys
    CAMQP_CUSTOM = 'camqp_custom'  #: camqp custom

    def __init__(
        self,
        event_processing=None,
        beat_processing=None,
        params=None,
        *args,
        **kwargs
    ):

        super(engine, self).__init__(*args, **kwargs)

        self.event_processing = event_processing
        self.beat_processing = beat_processing
        self.params = params

    @property
    def event_processing(self):
        """
        Event processing event_processing executed in the work
        """

        return self._event_processing

    @event_processing.setter
    def event_processing(self, value):
        """
        Change of event_processing.

        :param value: new event_processing to use. If None or wrong value,
            event_processing is used
        :type value: NoneType, str or function
        """

        # by default, load default event processing
        if value is None:
            value = event_processing
        # if str, load the related function
        elif isinstance(value, basestring):
            try:
                value = get_task(value)
            except ImportError:
                self.logger.error('Impossible to load %s' % value)
                value = event_processing

        # set _event_processing and work
        self._event_processing = value

    @property
    def beat_processing(self):
        """
        Task executed in the beat
        """

        return self._beat_processing

    @beat_processing.setter
    def beat_processing(self, value):
        """
        Change of beat_processing.

        :param value: new beat_processing to use. If None or wrong value,
            beat_processing is used
        :type value: NoneType, str or function
        """

        # by default, load default beat processing
        if value is None:
            value = beat_processing
        # if str, load the related function
        elif isinstance(value, basestring):
            try:
                value = get_task(value)
            except ImportError:
                self.logger.error('Impossible to load %s' % value)
                value = beat_processing

        # set _beat_processing and work
        self._beat_processing = value

    def work(self, event, msg, *args, **kwargs):

        result = self._event_processing(
            engine=self, event=event, msg=msg, logger=self.logger,
            *args, **kwargs
        )

        return result

    def beat(self, *args, **kwargs):
        self._beat_processing(
            engine=self, logger=self.logger,
            *args, **kwargs
        )

    @property
    def params(self):
        """
        Event processing event_processing parameters dictionary
        """

        result = {} if self._params is None else self._params

        return result

    @params.setter
    def params(self, value):
        """
        Change or params.

        :param value: new params to use. If str, it is evaluated.
        :type value: str or dict
        """

        if isinstance(value, basestring):
            value = eval(value)

        self._params = value

    def _conf(self, *args, **kwargs):

        result = super(engine, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=CATEGORY,
            new_content=(
                Parameter(engine.EVENT_PROCESSING),
                Parameter(engine.BEAT_PROCESSING),
                Parameter(engine.PARAMS, parser=eval),
                Parameter(engine.NEXT_AMQP_QUEUES),
                Parameter(engine.NEXT_BALANCED),
                Parameter(engine.NAME),
                Parameter(engine.BEAT_INTERVAL),
                Parameter(engine.EXCHANGE_NAME),
                Parameter(engine.ROUTING_KEYS),
                Parameter(engine.CAMQP_CUSTOM)))

        return result


def load_dynamic_engine(name, *args, **kwargs):
    """Load a new engine in adding a specific conf_path.

    :param str name: dynamic engine name.
    :param tuple args: used in new engine initialization such as varargs.
    :param dict kwargs: used in new engine initialization such as keywords.
    """

    conf_path = 'engines/%s.conf' % name

    # instantiate a new engine with a unified category which corresponds to
    # input name
    result = engine(unified_category=name, *args, **kwargs)

    # set conf_paths in adding input conf_path
    conf_paths = result.conf_paths
    conf_paths.append(conf_path)
    result.conf_paths = conf_paths

    # finally, reconfigure the result
    result.apply_configuration()

    # and returns it
    return result


def event_processing(engine, event, ctx=None, **params):
    """
    Event processing signature to respect in order to process an event.

    A condition may returns a boolean value.

    :param Engine engine: engine which process the event.
    :param dict event: event to process.
    :param dict ctx: event processing context.
    :param dict params: event processing additional parameters.
    """

    return event


def beat_processing(engine, **params):
    """
    Beat processing signature to respect in order to execute a periodic task.

    :param Engine engine: engine which executes the beat.
    :param dict params: beat processing additional parameters.
    """
