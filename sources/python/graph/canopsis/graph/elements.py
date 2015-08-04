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

"""
===========
Description
===========

Functional
==========

A graph is a mean to construct information based on logical understanding of
information.

A graph is composed of vertices such as elementary information and edges which
are logical relationships between information.

With the previous definition, a graph is a complex vertice which results in
providing information.

In such way, an information can be elementary or composed of graph information
such as hypergraphs.

An hypergraph permits to add a context dimension over graph elements. In such
structure, vertices, edges and graphs exist in multiple graphs, and their
surround depends on the graph they are associated at a "time".

Technical
=========

For simplification reasons, a graph is technically solved by such concepts.

Graph Element
-------------

The graph element is the base concept of elements decribed here. It has a
unique identifier among all graph elements and a type for graph specialisation
reasons.

A graph element contains::

    - id: unique identifier among all graph elements.
    - type: type of graph element. A graph could be a topology or something
        else.
    - _cls: python class path.
    - _type: base type which permits to recognize the type of element.

Graph vertice
-------------

A graph vertice inherits from the graph element and can contain data
information.

From a graph vertice, it is possible to resolve neighbour vertices thanks to
edges.

A graph vertice contains::

    - info: vertice information.
    - _type: equals vertice.

Graph edge
----------

Technically, a graph edge is more rich than its representation in the
functional definition in order to ease its manipulation in a richer context
instead of keeping only a logical use of edges. It becomes possible to describe
logical information between two edges.

A graph edge inherits from the graph vertice in order to transport information
and can bind several source with several targets, directly or not.

    - sources: source vertices.
    - targets: target vertices.
    - directed: directed orientation. If False, source and target vertices are
        directly connected, otherwise, only sources are directly connected to
        targets.
    - weight: edge weight. Default 1.
    - _type: equals edge.

A graph inherits from vertice and contains::

    - elts: elements existing in this graph.
    - _type: graph.
"""

from uuid import uuid4 as uuid

from collections import Iterable

from canopsis.common.init import basestring
from canopsis.storage.core import Storage
from canopsis.common.utils import lookup, path

CONF_PATH = 'graph/graph.conf'
CATEGORY = 'GRAPH'


class GraphElement(object):
    """Base class for all graph elements.

    Contains an ID and a type.
    """

    ID = Storage.DATA_ID  #: graph element id.
    TYPE = 'type'  #: graph element type name.
    BASE_TYPE = '_type'  #: base graph element type name.
    _CLS = '_cls'  #: graph element class type.

    __slots__ = (ID, TYPE)

    @staticmethod
    def new_id():
        """Generate a new graph element id.

        :return: new graph element id.
        :rtype: str
        """
        return str(uuid())

    def __init__(self, id=None, type=None, **kwargs):
        """
        :param str id: element id. Generated if None.
        :param str type: element type name. self lower type name if None.
        """

        super(GraphElement, self).__init__()

        self.type = self.__class__.__name__.lower() if type is None else type
        self.id = GraphElement.new_id() if id is None else id

    def __repr__(self):

        result = '{}:('.format(type(self).__name__)

        for index, slot in enumerate(type(self).__slots__):
            value = getattr(self, slot, None)
            result += '{}:{}, '.format(slot, value)

        result += ')'

        return result

    def __eq__(self, other):
        """
        :return: True if other is a GraphElement and public attributes sames.
        """

        # check other type and id
        result = isinstance(other, GraphElement) and self.id == other.id

        return result

    def __ne__(self, other):
        """
        :return: not self.__eq__ == other.
        """

        return not self.__eq__(other)

    def __hash__(self):
        """
        :return: self.id hash.
        """

        return hash(self.id)

    @classmethod
    def new(cls, **kwargs):
        """Instantiate a new element where kwargs contains element attributes.

        :param dict kwargs: new element attributes.
        :return: new element.
        :raises: TypeError if kwargs can not be used in cls.new
        """

        result = cls(**kwargs)
        return result

    @staticmethod
    def new_element(**elt_properties):
        """Instantiate a new graph element related to elt properties.

        :param dict elt_properties: serialized elt properties.
        :return: new elt instance.
        """

        result = None

        cls = elt_properties[GraphElement._CLS]

        if cls is not None:
            cls = lookup(cls)
            result = cls.new(**elt_properties)

        return result

    def to_dict(self):
        """Transform self to a dict in storing public attributes.
        """

        result = {}
        # set public attributes
        for slot in self.__slots__:
            if slot[0] != '_':
                result[slot] = getattr(self, slot)
        # set class type
        self_type = type(self)
        result[GraphElement._CLS] = path(self_type)
        # set base type
        result[GraphElement.BASE_TYPE] = self.BASE_TYPE
        return result

    def resolve_refs(self, elts, manager):
        """Resolve self references with input elts.

        :param dict elts: elts by id.
        """

        pass

    def save(self, manager, graph_ids=None, cache=False):
        """Save self into manager graphs.

        :param GraphManager manager: manager to use in order to save self.
        :param graph_ids: graph ids where save self.
        :type graph_ids: list or str
        :param bool cache: use query cache if True (False by default).
        """

        # save the dict format
        elt = self.to_dict()
        manager.put_elt(elt=elt, graph_ids=graph_ids, cache=cache)

    def delete(self, manager, cache=False):
        """Delete self from manager.

        :param GraphManager manager: manager from where delete self.
        :param bool cache: use query cache if True (False by default).
        """

        # delete self from manager
        manager.del_elts(ids=self.id, cache=cache)


class Vertice(GraphElement):
    """In charge of managing a Vertice.

    Contains an information.
    """

    INFO = 'info'  #: info attribute name

    BASE_TYPE = 'vertice'  # base type name

    __slots__ = GraphElement.__slots__ + (INFO,)

    def __init__(self, info=None, *args, **kwargs):
        """
        :param info: self info.
        """

        super(Vertice, self).__init__(*args, **kwargs)

        self.info = info

    def delete(self, manager, cache=False):

        super(Vertice, self).delete(manager=manager, cache=cache)

        self_id = self.id

        # remove verties which are linked to self
        edges = manager.get_edges(sources=self_id, targets=self_id)
        links = Edge.SOURCES, Edge.TARGETS
        for edge in edges:
            for link in links:
                if self_id in getattr(edge, link):
                    setattr(edge, link, [
                        elt_id for elt_id in edge[link] if elt_id != self.id
                    ])
            # delete the edge if sources or targets is empty
            if not (edge.sources and edge.targets):
                edge.delete(manager=manager, cache=cache)
            else:  # resolve_refs edge without self in sources and/or targets
                edge.save(manager=manager, cache=cache)


class Edge(Vertice):
    """In charge of managing an Edge.

    Contains sources, targets and a directed information.
    """

    BASE_TYPE = 'edge'  # base type name

    WEIGHT = 'weight'  #: weight edge attribute name
    SOURCES = 'sources'  #: source vertice ids attribute name
    TARGETS = 'targets'  #: target vertice ids attribute name
    DIRECTED = 'directed'  #: directed attribute name
    _DSOURCES = '_dsources'  #: dict sources vertices attribute name
    _DTARGETS = '_dtargets'  #: dict target vertices attribure name
    _GSOURCES = '_gsources'  #: graph source vertices attribute name
    _GTARGETS = '_gtargets'  #: graph target vertices attribure name

    DEFAULT_DIRECTED = True  #: default directed value

    __slots__ = (
        WEIGHT, SOURCES, TARGETS, DIRECTED,
        _DSOURCES, _DTARGETS, _GSOURCES, _GTARGETS,
    ) + Vertice.__slots__

    def __init__(
        self, weight=1, sources=None, targets=None, directed=DEFAULT_DIRECTED,
        _dsources=None, _dtargets=None, _gsources=None, _gtargets=None,
        *args, **kwargs
    ):
        """
        :param float weight: self weight.
        :param list sources: self sources.
        :param list targets: self targets.
        :param bool directed: self directed. (default DEFAULT_DIRECTED)
        :param list _dsources: dict sources.
        :param list _dtargets: dict targets.
        :param dict _gsources: graph vertice targets by id.
        :param dict _gtargets: graph vertice sources by id.
        """

        super(Edge, self).__init__(*args, **kwargs)

        # init weight
        self.weight = weight

        # init sources such as an array of ids
        if sources is None:
            self.sources = []
        else:
            # ensure sources and targets are list if they are string
            if isinstance(sources, (basestring, GraphElement)):
                sources = [sources]
            self.sources = list(
                source.id if isinstance(source, GraphElement) else source
                for source in sources
            )
        # init targets such as an array of ids
        if targets is None:
            self.targets = []
        else:
            if isinstance(targets, (basestring, GraphElement)):
                targets = [targets]
            self.targets = list(
                target.id if isinstance(target, GraphElement) else target
                for target in targets
            )

        self.directed = directed
        self._dsources = [] if _dsources is None else _dsources
        self._dtargets = [] if _dtargets is None else _dtargets
        self._gsources = {} if _gsources is None else _gsources
        self._gtargets = {} if _gtargets is None else _gtargets

    def resolve_refs(self, elts, manager):

        # init self._gsources and self._gtargets
        self._gsources = {}
        self._gtargets = {}

        for source in self.sources:
            if source not in elts:
                elt = manager.get_elts(ids=source)
                elts[source] = elt
            else:
                elt = elts[source]
            self._gsources[source] = elt
        for target in self.targets:
            if target not in elts:
                elt = manager.get_elts(ids=target)
                elts[target] = elt
            else:
                elt = elts[target]
            self._gtargets[target] = elt

    def del_refs(self, ids=None, sources=None, targets=None):
        """Del references of vertices.

        :param ids: vertice ids to remove from self references.
        :type ids: list or str
        :param sources: vertice sources to remove from self references.
        :type sources: list or str
        :param targets: vertice targets to remove from self references.
        :type targets: list or str
        """

        # init params
        if ids is not None:
            # if ids exist, add it to sources and targets
            if isinstance(ids, basestring):
                ids = [ids]
            if sources is None:
                sources = ids
            else:
                if isinstance(sources, basestring):
                    sources = [sources] + ids
                else:
                    sources += ids
            if targets is None:
                targets = ids
            else:
                if isinstance(targets, basestring):
                    targets = [targets] + ids
                else:
                    targets += ids

        # remove sources from self.sources
        if sources is not None:
            if isinstance(sources, basestring):
                sources = [sources]
            for source in sources:
                while source in self.sources:
                    self.sources.remove(source)
                # remove dsources
                self._dsources = [
                    src for src in self._dsources if src not in source
                ]
                # remove gsources
                if source in self._gsources:
                    del self._gsources[source]

        # remove targets from self.targets
        if targets is not None:
            if isinstance(targets, basestring):
                targets = [targets]
            for target in targets:
                while target in self.targets:
                    self.targets.remove(target)
                # remove dtargets
                self._dtargets = [
                    src for src in self._dtargets if src not in target
                ]
                # remove gtargets
                if target in self._gtargets:
                    del self._gtargets[target]


class Graph(Vertice):
    """In charge of managing a Graph.

    Contains graph elements.
    """

    BASE_TYPE = 'graph'  # base type name

    ELTS = 'elts'  #: elt ids attribute name.
    _GELTS = '_gelts'  #: graph elts attribute name.
    _DELTS = '_delts'  #: dict elts attribute name.
    _SOURCES = '_sources'  #: edge by source vertices attribute name
    _TARGETS = '_targets'  #: edge by target vertices attribute name

    VERTICE_TYPE = Vertice  #: vertice type to use
    EDGE_TYPE = Edge  #: edge type to use

    _UPDATING = '_updating'  #: private attribute name while updating

    __slots__ = (
        ELTS,
        _GELTS, _DELTS, _UPDATING, _SOURCES, _TARGETS
    ) + Vertice.__slots__

    def __init__(
        self,
        elts=None, _gelts=None, _delts=None, _sources=None, _targets=None,
        *args, **kwargs
    ):
        """
        :param list elts: self graph elt ids.
        :param dict _delts: self dict graph elts by id.
        :param dict _gelts: self graph elements by id.
        :param dict _sources: edges by source vertice id.
        :param dict _targets: edges by target vertice id.
        """

        super(Graph, self).__init__(*args, **kwargs)

        # init elts such as an array of ids
        if elts is None:
            self.elts = []
        else:
            if isinstance(elts, basestring):
                elts = [elts]
            self.elts = list(
                elt.id if isinstance(elt, GraphElement) else elt
                for elt in elts
            )

        self._delts = {} if _delts is None else _delts
        self._gelts = {} if _gelts is None else _gelts
        self._updating = False
        self._sources = {} if _sources is None else _sources
        self._targets = {} if _targets is None else _targets

    def __contains__(self, other):
        """True if other is a vertice inside self.

        :param other: possible vertice(s) or vertice(s) id(s) in self graph.
        :type other: str or dict or GraphElement or list of other.
        """

        # by default, result is False
        result = False

        if isinstance(other, basestring):
            result = other in self.elts
        elif isinstance(other, dict):
            result = other[GraphElement.ID] in self.elts
        elif isinstance(other, GraphElement):
            result = other.id in self.elts
        elif isinstance(other, Iterable):
            # ids is a set of graph element ids to compare with self.elts
            ids = set()
            # quick access to the field name GraphElement ID
            id = GraphElement.ID
            for item in other:
                if isinstance(item, basestring):
                    ids.add(item)
                elif isinstance(item, dict):
                    ids.add(item[id])
                elif isinstance(item, GraphElement):
                    ids.add(item.id)
                else:
                    raise Exception(
                        'item {0} is not supported.'.format(item)
                    )
            # transform self elts in a set in order to do quick comparison
            self_elts = set(self.elts)
            result = (self_elts & ids) == ids
        else:
            raise Exception('elt {0} is not supported.'.format(other))

        return result

    def resolve_refs(self, elts, manager):

        if not self._updating:
            self._updating = True
            for gelt in self._gelts:
                gelt.resolve_refs(elts, manager=manager)
                # update self _sources and _targets
                if isinstance(gelt, Edge):
                    _gsources = gelt._gsources
                    for source in _gsources:
                        gsource = _gsources[source]
                        if gsource.id not in self._sources:
                            self._sources[gsource.id] = [gelt]
                        else:
                            self._sources[gsource.id].append(gelt)
                    _gtargets = gelt._gtargets
                    for target in _gtargets:
                        gtarget = _gtargets[target]
                        if gtarget.id not in self._targets:
                            self._targets[gtarget.id] = [gelt]
                        else:
                            self._targets[gtarget.id].append(gelt)
            self._updating = False

    def update_gelts(self, manager, depth=0, _added_elts=None):
        """Update self graph elts with self elt ids and input manager.

        :param manager: self manager to use.
        :param int depth: graph depth (add "depth" levels of graphs).
        :param _added_elts: private parameter for storing new graph elts and
            avoid recursive calls.
        """

        # init _add_elts
        if _added_elts is None:
            _added_elts = {self.id: self}
        # initialize self graph elts
        self._gelts = []
        # get elts
        elts = manager.get_elts(ids=self.elts, serialize=False)
        # for all elt ids
        for elt in elts:
            elt_id = elt[GraphElement.ID]
            # if elt already exists in memory
            if elt_id in _added_elts:
                # get it
                new_elt = _added_elts[elt_id]
            else:  # else, instantiate it
                new_elt = GraphElement.new_element(**elt)
                _added_elts[elt_id] = new_elt
                # fill graph if depth > 0
                if isinstance(new_elt, Graph) and depth > 0:
                    new_elt.update_gelts(
                        manager=manager,
                        depth=depth - 1,
                        _added_elts=_added_elts
                    )
            self._gelts.append(new_elt)
        self.resolve_refs(_added_elts, manager=manager)

    def add_elts(self, *elts):
        """Add elts in self.

        :param elts: elts to add.
        :type elts: tuple of int or dict or GraphElement
        """

        for elt in elts:
            if isinstance(elt, basestring):
                if elt not in self.elts:
                    self.elts.append(elt)
            elif isinstance(elt, dict):
                elt_id = elt[GraphElement.ID]
                if elt_id not in self._delts:
                    self._delts[elt_id] = elt
                    if elt_id in self.elts:
                        self.elts.remove(elt_id)
            elif isinstance(elt, GraphElement):
                elt_id = elt.id
                if elt_id not in self._gelts:
                    self._gelts[elt_id] = elt
                    elt_dict = elt.to_dict()
                    if elt_id not in self._delts:
                        self._delts[elt_id] = elt_dict
                        if elt_id not in self.elts:
                            self.elts.append(elt_id)
            elif isinstance(elt, Iterable):
                self.add_elts(*elt)

    def remove_elts(self, *elts):
        """
        :param elts: elts or elt to remove.
        :type elts: tuple of str, dict or GraphElement
        """

        for elt in elts:
            if isinstance(elt, basestring):
                if elt in self.elts:
                    self.elts.remove(elt)
            elif isinstance(elt, dict):
                elt_id = elt[GraphElement.ID]
                if elt_id in self._delts:
                    del self._delts[elt_id]
                    if elt_id in self.elts:
                        self.elts.remove(elt_id)
            elif isinstance(elt, GraphElement):
                elt_id = elt.id
                if elt_id in self._gelt:
                    del self._gelt[elt_id]
                    if elt_id in self._delts:
                        del self._delts[elt_id]
                        if elt_id in self.elts:
                            self.elts.remove(elt_id)
            elif isinstance(elt, Iterable):
                self.remove_elts(*elt)

    def to_dict(self):

        result = super(Graph, self).to_dict()

        # if self _elts not empty
        if self._delts:  # add delts ids which are not present in elts
            elts = result[Graph.ELTS]
            for delt_id in self._delts:
                if delt_id not in elts:
                    elts.append(delt_id)

        # if gelts not empty
        if self._gelts:
            elts = result[Graph.ELTS]
            for gelt_id in self._gelts:
                if gelt_id not in elts:
                    elts.append(gelt_id)

        return result

    def get_neighbours(self, vertice):
        """Get neighbours vertices by edges of input vertices in respecting
        directed edges.

        :param Vertice vertice: vertice from where find linked vertices by
            edges.
        :return: dict of vertice by edges.
        :rtype: dict
        """

        result = {}
        # get target edges
        target_edges = self._sources[vertice]

        # add all targets
        for target_edge in target_edges:
            if target_edge in result:
                result[target_edge] += target_edge._gtargets
            else:
                result[target_edge] = [target_edge._gtargets]
            # add edge sources if not directed
            if not target_edge.directed:
                sources = list(target_edge._gsources)
                sources.remove(vertice)
                result[target_edge] += sources

        source_edges = self._targets[vertice]

        # add all sources which are not directed
        for source_edge in source_edges:
            if not source_edge.directed:
                if source_edge in result:
                    result[source_edge] += source_edge._gsources
                else:
                    result[source_edge] = [source_edge._gsources]
                targets = list(source_edge._gtargets)
                targets.remove(vertice)
                result[source_edge] += targets

        return result

    def save(self, manager, graph_ids=None, cache=False, *args, **kwargs):

        self.elts = list(set(self.elts))

        super(Graph, self).save(
            manager=manager,
            graph_ids=graph_ids,
            cache=cache,
            *args, **kwargs
        )
        # save graph elements
        if self._gelts:
            for elt_id in self._gelts:
                elt = self._gelts[elt_id]
                elt.save(manager=manager, cache=cache)

        elif self._delts:
            for elt_id in self._delts:
                elt = self._delts[elt_id]
                elt = Graph.new_element(**elt)
                elt.save(manager=manager, cache=cache)

    def delete(self, manager, cache=False, del_orphans=True, *args, **kwargs):

        super(Graph, self).delete(
            manager=manager, cache=cache, *args, **kwargs
        )

        if del_orphans:
            manager.remove_elts(
                ids=self.elts, del_orphans=del_orphans, cache=cache
            )
