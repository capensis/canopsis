# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from canopsis.common.util import resolve_element
from canopsis.engines import Engine
from canopsis.configuration import (
        Configurable, add_category, conf_paths, Parameter
        )

CONF_PATH = 'engines/engines.conf'  #: dynamic engine configuration path
CATEGORY = 'ENGINE'  #: dynamic engine configuration category


@add_category(CATEGORY)
@conf_paths(CONF_PATH)
class DynamicEngine(Engine, Configurable):
    """
    Engine which is able to load dynamically its event processing through
    configuration properties.

    :var str task: event processing task path.
    :var dict params: event processing task parameters.
    """

    TASK = 'task'  #: task field name
    PARAMS = 'params'  #: params field name

    NEXT_AMQP_QUEUES = 'next_amqp_queues'
    NEXT_BALANCED = 'next_balanced'
    NAME = 'name'
    BEAT_INTERVAL = 'beat_interval'
    EXCHANGE_NAME = 'exchange_name'
    ROUTING_KEYS = 'routing_keys'
    CAMQP_CUSTOM = 'camqp_custom'

    def __init__(
        self,
        task=None,
        params=None,
        *args,
        **kwargs
    ):

        super(DynamicEngine, self).__init__(*args, **kwargs)

        self.task = task
        self.params = params

    @property
    def task(self):
        """
        Event processing task executed in the work
        """

        return self._task

    @task.setter
    def task(self, value):
        """
        Change of task.

        :param value: new task to use. If None or wrong value, event_processing
         is used
        :type value: NoneType, str or function
        """

        # by default, load default event processing
        if value is None:
            value = event_processing
        # if str, load the related function
        elif isinstance(value, str):
            try:
                value = resolve_element(value)
            except ImportError:
                self.logger.error('Impossible to load %s' % value)
                value = event_processing

        # set _task and work
        self._task = self.work = value

    @property
    def params(self):
        """
        Event processing task parameters dictionary
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

        if isinstance(value, str):
            value = eval(value)

        self._params = value

    def _conf(self, *args, **kwargs):

        result = super(DynamicEngine, self)._conf(*args, **kwargs)

        result.add_unified_category(CATEGORY,
            Parameter(DynamicEngine.TASK, self.task),
            Parameter(DynamicEngine.PARAMS, self.params, eval),
            Parameter(DynamicEngine.NEXT_AMQP_QUEUES, self.next_amqp_queues),
            Parameter(DynamicEngine.NEXT_BALANCED, self.next_balanced),
            Parameter(DynamicEngine.NAME, self.name),
            Parameter(DynamicEngine.BEAT_INTERVAL, self.beat_interval),
            Parameter(DynamicEngine.EXCHANGE_NAME, self.exchange_name),
            Parameter(DynamicEngine.ROUTING_KEYS, self.routing_keys),
            Parameter(DynamicEngine.CAMQP_CUSTOM, self.camqp_custom))

        return result

    def _configure(self, unified_conf, *args, **kwargs):

        super(DynamicEngine, self)._configure(
            unified_conf=unified_conf, *args, **kwargs)

        params = [
            DynamicEngine.TASK,
            DynamicEngine.PARAMS,
            DynamicEngine.NEXT_BALANCED,
            DynamicEngine.NEXT_AMQP_QUEUES,
            DynamicEngine.NAME,
            DynamicEngine.BEAT_INTERVAL,
            DynamicEngine.EXCHANGE_NAME,
            DynamicEngine.ROUTING_KEYS,
            DynamicEngine.CAMQP_CUSTOM]

        for param in params:
            self._update_property(
                unified_conf=unified_conf,
                param_name=param)


def load_dynamic_engine(conf_path, category, *args, **kwargs):
    """
    Load a new engine in adding a specific conf_path.

    :param str conf_path: final conf_path to set at the end of new dynamic
        engine conf paths

    :param tuple args: used in new engine initialization such as *args
    :param dict kwargs: used in new engine initialization such as **kwargs
    """

    # instantiate a new engine with a unified category which corresponds to
    # input category
    engine = DynamicEngine(
        unified_category=category,
        *args, **kwargs)

    # set conf_paths in adding input conf_path
    conf_paths = engine.conf_paths
    conf_paths.append(conf_path)
    engine.conf_paths = conf_paths

    # finally, reconfigure the engine
    engine.apply_configuration()

    # and returns it
    return engine


def event_processing(event, ctx=None, **params):
    """
    Event processing signature to respect in order to process an event.

    A condition may returns a boolean value.

    :param dict event: event to process
    :param dict ctx: event processing context
    :param dict params: event processing additional parameters
    """

    raise NotImplementedError()
