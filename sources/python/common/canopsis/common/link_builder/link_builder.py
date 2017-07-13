# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from abc import ABCMeta, abstractmethod
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category
from canopsis.configuration.model import Parameter, ParamList

CONF_FILE = "common/link_builder/manager.conf"
BUILDERS_CAT = "LINK_BUILDERS"
BUILDERS_CONTENT = [ParamList(parser=Parameter.bool)]

DEFAULT_BUILDER_CAT = "DEFAULT_BUILDER"
DEFAULT_BUIDLER_CONTENT = [
    Parameter("column"), Parameter("baseurl"), Parameter(
        "managed_entities", Parameter.array())
]


@conf_paths(CONF_FILE)
@add_category(BUILDERS_CAT, content=BUILDERS_CONTENT)
class HypertextLinkManager:
    def __init__(self, config):
        self.config = config

    def links_for_entity(self, entity, options):
        """Generate links for the entity with the builder specify in the
        configuration.

        :param dict entity: the entity to handle
        :param options: the options
        :return list: a list of links as a string."""
        pass


class HypertextLinkBuilder:
    __metaclass__ = ABCMeta

    @abstractmethod
    def build(self, entity, options):
        """Build links from an entity and the given option
        :param dict entity: the entity to handle
        :param options: the options
        :return list: a list of links as a string."""
        pass


@conf_paths(CONF_FILE)
@add_category(DEFAULT_BUILDER_CAT, content=DEFAULT_BUIDLER_CONTENT)
class BasicLinkBuilder(HypertextLinkBuilder):
    def __init__(self, options):
        self.options = options

    def build(self, entity, options):
        pass
