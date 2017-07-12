# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from abc import ABCMeta, abstractmethod
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category
from canopsis.configuration.model import Parameter, ParamList

CONF_FILE = "link_builder/manager.conf"
CONF_CAT = "LINK_BUILDER"
CONF_CONTENT = [
    ParamList(parser=Parameter.bool)
]


@conf_paths(CONF_FILE)
@add_category(CONF_CAT, content=CONF_CONTENT)
class HypertextLinkManager:

    def __init__(self, config):
        self.config = config

    def links_for_entity(self, entity, options):
        """Generate links for the entity with the builder specify in the
        configuration.

        :param dict entity: the entity to handle
        :param options: the options
        :return list: a list of links as a string."""



class HypertextLinkBuilder:
    __metaclass_ = ABCMeta

    @abstractmethod
    def build(self, entity, options):
        """Build links from an entity and the given option
        :param dict entity: the entity to handle
        :param options: the options
        :return list: a list of links as a string."""
        pass

class BasicLinkBuilder(HypertextLinkBuilder):

    def __init__(self, options):
        pass

    def build(self, entity, options):
        pass
