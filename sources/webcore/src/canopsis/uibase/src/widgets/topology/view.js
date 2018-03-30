/*eslint-disable*/
/*
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
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
*/

/**
* View model
* ----------
*
* This widget uses 2 layers of models:
* - Dom element ---> view element (node and link).
* View elements are stored in the shapesById dictionary where ids correspond to controller record ids.
* Dom element to view element: field __data__
* view element to canopsis record: using view id in controller.graphModel.recordsById.
* Both Dom and view elements have the same id.
*
* A view element contains such view properties:
* - hidden: if True, the element is not displayed.
* - fixed: if True, the element is not impacted by the auto-layout.
* - x, y: respective abscisse and ordinate coordinates.
*
* All view elements are saved in the property shapesById.
* nodes and links properties are dedicated to the layout engine.
* The process begins updating the view model, then in drawing shapes.
* During update, records are retrieved form the model, and if existing id already exist in view model, the view keeps existing view properties (coordinates, etc.) and update business properties (dynamic weight, etc.).
* Drawing step can be fired by the 'redraw' event from the controller.
*
* Layout
* ------
*
* Several layouts are provided.
*
* All respect common properties:
* - size: screen size.
* - translate: translation information.
* - scale: zoom information.
*
* And specific properties are given in the 'layout' property:
* - type: a string value corresponding to the layout type among force, cluster, tree, pack and partition.
* - activated: a boolean value which specifies if the auto layout is activated.
* - {force, cluster, pack, partition, tree}: dedicated layout properties.
* - engine: auto layout engine.
*
* Therefore, if you want the current layout properties 'P', you have to do view.layout[view.layout.type].P.
*
* Toolbox
* -------
*
* A toolbox is provided in order to enrich interaction with the graph. Its instance is managed in the toolbox property. This last contains data or null if toolbox is not activated.
*/
Ember.Application.initializer({
    name: 'TopologyViewMixin',
    initialize: function(container, application) {
    var d3 = window.d3;

    var get = Ember.get,
        set = Ember.set;

    var TopologyViewMixin = Ember.Mixin.create({

        shapesById: {}, // dictionary of shapes by id

        nodes: [], // list of nodes
        links: [], // list of links

        nodeClass: 'node', // node class name
        linkClass: 'link', // link class name
        edgeClass: 'edge', // edge class name

        toolbox: null, // global toolbox which may be unique

        _layout: {
            type: 'force', // current layout type
            force: null,
            tree: null,
            cluster: null,
            partition: null,
            pack: null,
            activated: true, // if true, layout is activated,
            engine: null, // layout engine
        }, // current layout parameters

        panel: null, // display panel

        source: null, // source position in multi selection (link, multi-select, etc.)

        toolboxGap: 50, // toolbox gap between shapes
        duration: 500, // transition duration in ms

        _unit: 32, // unit shape size

        translate: [1, 1], // translation
        scale: 1.5, // layout scale

        _defaultLayout: {
            partition: {
                sort: function comparator(a, b) {
                    var result = d3.ascending(a.name, b.name);
                    return result;
                },
                children: function children(d) {
                    return d.sources;
                },
            }, // default pack layout parameters

            pack: {
                sort: function comparator(a, b) {
                    var result = d3.ascending(a.name, b.name);
                    return result;
                },
                children: function children(d) {
                    return d.sources;
                },
            }, // default pack layout parameters

            cluster: {
                separation: function separation(a, b) {
                    var result = (a.parent == b.parent ? 1 : 2) / a.depth;
                    return result;
                },
                nodeSize: null,
                sort: function comparator(a, b) {
                    var result = d3.ascending(a.name, b.name);
                    return result;
                },
                children: function children(d) {
                    return d.sources;
                },

            }, // default tree layout parameters

            tree: {
                separation: function separation(a, b) {
                    var result = (a.parent == b.parent ? 1 : 2) / a.depth;
                    return result;
                },
                nodeSize: null,
                sort: function comparator(a, b) {
                    var result = d3.ascending(a.name, b.name);
                    return result;
                },
                children: function children(d) {
                    return d.sources;
                },

            }, // default tree layout parameters

            force: {
                charge: -30, // charge
                linkDistance: 20, // link distance
                linkStrength: 1, // link strength
                friction: 0.9, // friction
                theta: 0.8, // theta
                gravity: 0.1, // gravity
            }, // default force parameters
        }, // default layout parameters

        optionsCharge: function() {
            return {
                min: -256,
                max: 256
            };
        }.property('optionsCharge'),
        linkStrengthCharge: function() {
            return {
                min: 0,
                max: 1,
                step: 0.05
            };
        }.property('linkStrengthCharge'),
        linkDistanceCharge: function() {
            return {
                min: 0,
                max: 150,
                step: 10
            };
        }.property('linkDistanceCharge'),
        frictionCharge: function() {
            return {
                min: 0,
                max: 1,
                step: 0.05
            };
        }.property('frictionCharge'),
        thetaCharge: function() {
            return {
                min: 0,
                max: 1,
                step: 0.05
            };
        }.property('thetaCharge'),
        gravityCharge: function() {
            return {
                min: 0,
                max: 1,
                step: 0.05
            };
        }.property('gravityCharge'),

        filter: function() {
            return this._filter;
        }.property('filter'),
        filterChanged: function() {
            this._filter = get(this, 'filter');
        }.observes('filter'),
        autoLayout: function() {
            return this._layout.activated;
        }.property('autoLayout'),
        autoLayoutChanged: function() {
            this._layout.activated = get(this, 'autoLayout');
            if (this._layout.activated) {
                this._layout.engine.start();
            } else {
                this._layout.engine.stop();
            }
        }.observes('autoLayout'),
        charge: function() {
            var result = this._layout[this._layout.type].charge;
            return result;
        }.property('charge'),
        chargeChanged: function() {
            var charge = get(this, 'charge');
            this._layout[this._layout.type].charge = charge;
            this._layout.engine.charge(charge);
            this.redraw();
        }.observes('charge'),
        linkDistance: function() {
            var result = this._layout[this._layout.type].linkDistance;
            return result;
        }.property('linkDistance'),
        linkDistanceChanged: function() {
            var linkDistance = this.get('linkDistance');
            this._layout[this._layout.type].linkDistance = linkDistance;
            this._layout.engine.linkDistance(linkDistance);
            this.redraw();
        }.observes('linkDistance'),
        linkStrength: function() {
            var result = this._layout[this._layout.type].linkStrength;
            return result;
        }.property('linkStrength'),
        linkStrengthChanged: function() {
            var linkStrength = this.get('linkStrength');
            this._layout[this._layout.type].linkStrength = linkStrength;
            this._layout.engine.linkStrength(linkStrength);
            this.redraw();
        }.observes('linkStrength'),
        friction: function() {
            var result = this._layout[this._layout.type].friction;
            return result;
        }.property('friction'),
        frictionChanged: function() {
            var friction = this.get('friction');
            this._layout[this._layout.type].friction = friction;
            this._layout.engine.friction(friction);
            this.redraw();
        }.observes('friction'),
        theta: function() {
            var result = this._layout[this._layout.type].theta;
            return result;
        }.property('theta'),
        thetaChanged: function() {
            var theta = this.get('theta');
            this._layout[this._layout.type].theta = theta;
            this._layout.engine.theta(theta);
            this.redraw();
        }.observes('theta'),
        gravity: function() {
            var result = this._layout[this._layout.type].gravity;
            return result;
        }.property('gravity'),
        gravityChanged: function() {
            var gravity = this.get('gravity');
            this._layout[this._layout.type].gravity = gravity;
            this._layout.engine.gravity(gravity);
            this.redraw();
        }.observes('gravity'),
        showProperties: function() {
            return this._showProperties;
        }.property('showProperties'),

        rerender: function() {
            this._super.apply(this, arguments);
            this.didInsertElement();
        },

        //Controller -> View Hooks
        registerHooks: function() {
            result = this._super.apply(this, arguments);
            get(this, "controller").on('redraw', this, this.redraw);
            return result;
        },

        unregisterHooks: function() {
            result = this._super.apply(this, arguments);
            get(this, "controller").off('redraw', this, this.redraw);
            return result;
        },

        didInsertElement: function() {
            this._super.apply(this, arguments);
            if (get(this, 'controller').graphModel.graph === null) {
                return;
            }
            // update the view
            this.updateModel();
        },

        /**
        * Update the view model related to model records.
        */
        updateModel: function() {
            // get global reference to this in updateModel function
            var me = this;
            // get record elements by id
            var recordsById = get(this, 'controller').graphModel.recordsById;
            // add nodes and links in graph
            // contain nodes by entity_ids
            var nodesByEntityIds = {};
            /**
            * Save node related to entity ids in memory.
            */
            var saveRefToEntity = function(node) {
                var nodeId = node.id;
                var info = recordsById[nodeId].get('info');
                // add reference between entity id and node
                if (info) {
                    var entity = info.entity;
                    if (entity !== undefined) {
                        if (nodesByEntityIds[entity] === undefined) {
                            nodesByEntityIds[entity] = [node];
                        } else {
                            nodesByEntityIds[entity].push(node);
                        }
                    }
                }
            };

            // delete old shapes
            var shapesById = this.shapesById;
            var arrayToDelete = {
                nodes: [],
                links: []
            };
            Object.keys(shapesById).forEach(
                function(shapeId) {
                    if (recordsById[shapeId] === undefined) {
                        var shape = shapesById[shapeId];
                        if (shape.type === 'link') {
                            arrayToDelete.links.push(shape.index);
                        } else {
                            arrayToDelete.nodes.push(shape.index);
                        }
                        delete shapesById[shapeId];
                    }
                }
            );
            ['nodes', 'links'].forEach(
                function(type) {
                    indexToDelete = arrayToDelete[type];
                    var count = 0;
                    indexToDelete.forEach(
                        function(index) {
                            this[type].splice(index - count++, 1);
                        },
                        this
                    );
                },
                this
            );

            Object.keys(recordsById).forEach(
                function(recordId) {
                    var record = recordsById[recordId];
                    // get record shape from record
                    var recordShape = this.getNode(record);
                    // update its reference in shapesById
                    this.shapesById[recordId] = recordShape;
                    if (record.get('_type') === 'edge') {
                        // add links if record is an edge.
                        this.weaveLinks(recordShape);
                    } else {
                        // initialize weight
                        recordShape._weight = 1;
                    }
                    // register the node in memory
                    saveRefToEntity(recordShape);
                },
                this
            );

            // resolve weight
            this.nodes.forEach(
                function(node, i) {
                    var nodeId = node.id;
                    var record = recordsById[nodeId];
                    if (record.get('_type') === 'edge') {
                        var sourceId = record.get('sources')[0];
                        var sourceNode = this.shapesById[sourceId];
                        if (sourceNode !== undefined) {
                            sourceNode._weight += record.get('weight');
                        }
                    }
                    node.index = i;
                    node.entity = undefined;
                },
                this
            );

            // resolve sources and targets in links
            this.links.forEach(
                function(link) {
                    if (typeof link.source === 'string') {
                        link.source = this.shapesById[link.source];
                    }
                    if (typeof link.target === 'string') {
                        var parentShape = this.shapesById[link.target];
                        link.parent = link.target = parentShape;
                    }
                },
                this
            );

            // resolve entities
            get(this, 'controller').getEntitiesFromServer(
                Object.keys(nodesByEntityIds),
                function(result) {
                    if (result.total !== 0) {
                        result.data.forEach(
                            function(entity) {
                                nodesByEntityIds[entity._id].forEach(
                                    function(node) {
                                        node.entity = entity;
                                    }
                                );
                            }
                        );
                    }
                    me.trigger('redraw');
                }
            );
        },

        redraw: function() {
            // get global reference to this for inner functions
            var me = this;
            // calculate the right size
            var width = get(this, 'controller.model.width');
            if (!width) width = this.$().width();
            var height = get(this, 'controller.model.height');
            if (!height) height = width * 9 / 16;
            // get the right layout
            var layoutType = get(this, 'controller.model.layout');
            var layout = this._layout;
            layout.type = layoutType;
            // initialize the engine if necessary
            var engine = layout.engine;
            if (layout[layoutType] === null) {
                // initialize layout with default properties
                engine = layout.engine = d3.layout[layoutType]();
                layout[layoutType] = this._defaultLayout[layoutType];
            }
            // apply layout properties
            /*Object.keys(layout[layoutType]).forEach(
                function(layout_property_id) {
                    engine[layout_property_id](layout[layoutType][layout_property_id]);
                }
            );*/
            // apply size
            engine.size([width, height]);
            // get panel
            this.panel = d3.select(this.$('svg .panel')[0]);
            if (this.panel.size() === 0) {
                /**
                * zoom function.
                */
                var zoom = function() {
                    /*me.eventZoom = d3.event;
                    console.log(me.eventZoom);
                    if (d3.event.sourceEvent.type !== 'mousemove') {*/
                        if (!me.dragging) {
                            console.log(me.translate, d3.event.translate);
                            var translate = [];
                            [0, 1].forEach(function(index) {
                                translate[index] = (d3.event.translate[index]);// me.translate[index] + d3.event.translate[index]) * d3.event.scale;
                            });
                            me.translate = translate;
                            me.scale = Math.max(d3.event.scale, 1.5);
                            var transform = "translate(" + translate + ")scale(" + me.scale + ")";
                            me.panel.attr("transform", transform);
                        }
                    /*} else {
                        var translate = [d3.event.translate[0] * d3.event.scale, d3.event.translate[1] * d3.event.scale];
                        me.panel.attr("transform", "translate(" + translate + ")scale(" + d3.event.scale + ")");
                    }*/
                };
                var drag = function() {
                    var translate = [
                        me.translate[0] + d3.event.dx,
                        me.translate[1] + d3.event.dy
                    ];
                    me.translate = translate;
                    me.panel.attr('transform', 'translate(' + translate + ') scale(' + me.scale + ')');
                    /*me.panel.select('.overlay')
                        .attr(
                            {
                                'x': -translate[0],
                                'y': -translate[1]
                            }
                        )
                    ;*/
                };
                var dragstart = function(d, i) {
                    force.stop(); // stops the force auto positioning before you start dragging
                };
                var dragmove = function(d, i) {
                    d.px += d3.event.dx;
                    d.py += d3.event.dy;
                    d.x += d3.event.dx;
                    d.y += d3.event.dy;
                    tick(); // this is the key to make it work together with updating both px,py,x,y on d !
                };
                var dragend = function(d, i) {
                    d.fixed = true; // of course set the node to fixed so the force doesn't include the node in its auto positioning stuff
                    tick();
                    if (layout.activated) {
                        force.resume();
                    }
                };
                var drag = d3.behavior.drag()
                    .on('dragstart', dragstart)
                    .on('drag', dragmove)
                    .on('dragend', dragend)
                    ;
                var zoom = d3.behavior.zoom()
                    //.translate([0, 0])
                    //.scale(1)
                    .scaleExtent([1.5, 10])
                    .on('zoom', zoom)
                ;
                // or create it if it does not exist
                this.panel = d3.select(this.$('svg')[0])
                    .attr(
                        {
                            width: width,
                            height: height
                        }
                    )
                    .append('g')
                        .classed(
                            {
                                panel: true
                            }
                        )
                        .attr(
                            {
                            //    transform: 'scale(1.5)'
                            }
                        )
                        .on('click', function(){ return me.clickAction(this); })
                        .on('mousemove', function() { return me.moveAction(this); })
                        //.on('mouseover', function() { return me.overAction(this);})
                        .on('mouseout', function() { return me.overAction(this); })
                        .on('dblclick', function() { return me.addHandler(); })
                        .call(zoom)
                        //.call(drag)
                ;
            }
            // recover translate and scale
            var transform = "translate(" + this.translate + ")scale(" + this.scale + ")";
            this.panel.attr("transform", transform);

            // apply an overlay for better graph zooming and selection
            var overlay = this.panel.select('.overlay');
            if (overlay.size() === 0) {
                this.panel
                    .append("rect")
                        .classed("overlay", true)
                        .attr("width", width)
                        .attr("height", height);
            }
            // save overlay in memory
            me.overlay = overlay;
            // load nodes and links into engine layout
            engine.nodes(this.nodes).links(this.links);
            // get link model
            var linkModel = this.panel.selectAll('.' + this.linkClass)
                .data(this.links, function(link) {return link.id;});

            // process links to delete, add and update.
            var linksToAdd = linkModel.enter();
            this.addLinks(linksToAdd);
            var linksToUpdate = linkModel;
            this.updateLinks(linksToUpdate);
            var linksToDelete = linkModel.exit();
            this.delLinks(linksToDelete);

            // get node model
            var nodeModel = this.panel.selectAll('.' + this.nodeClass)
                .data(this.nodes, function(node) { return node.id; });

            // process nodes to delete, add and update.
            var nodesToAdd = nodeModel.enter();
            this.addNodes(nodesToAdd);
            var nodesToUpdate = nodeModel;
            this.updateNodes(nodesToUpdate);
            var nodesToDelete = nodeModel.exit();
            this.delNodes(nodesToDelete);

            // refresh selected shapes
            this.refreshSelectedShapes();

            // refresh locked shapes
            this.refreshLockedShapes();

            engine.on(
                'tick',
                function() {
                    if (!layout.activated) {
                        engine.stop();
                    }
                    linkModel
                        .attr(
                            {
                                "d": function(d) {
                                    return "M"+d.source.x+" "+d.source.y+"L"+d.target.x+" "+d.target.y;
                                }
                            }
                        )
                    ;
                    nodeModel
                        .attr(
                            {
                                "transform": function(d) {
                                    return "translate(" + d.x + "," + d.y + ")";
                                }
                            }
                        )
                    ;
                }
            );
            // finish to run the layout
            if(layout.activated) {
                engine.start();
            }
        },

        /**
        * Weave edge links from a d3 elt.
        *
        * @param edge d3 graph edge.
        */
        weaveLinks: function(d3Edge) {
            var me = this;
            // save count of linkId per targets/sources
            var recordsById = get(this, 'controller').graphModel.recordsById;
            var id = d3Edge.id;
            var edge = recordsById[id];
            // cause a d3 link has one source and one target, we may create as many link as there are composition of sources and targets.
            var sources = edge.get('sources');
            var targets = edge.get('targets');
            // update weight
            d3Edge._weight = edge.get('weight');
            // save links in d3Edge links
            d3Edge.sources = {};
            d3Edge.targets = {};
            var neighbourCounts = {
                sources: {},
                targets: {}
            };
            // array of links view pos to delete at the end for better complexity time resolution
            var linkPosToDelete = [];
            /**
            * Add links in d3Edge.
            * @param key value among targets and sources
            */
            var updateLinks = function(key) {
                var isSource = key === 'sources';
                var neighbours = edge.get(key);

                neighbours.forEach(
                    function(d) {
                        // get number of time we meet d
                        var neighbourCount = neighbourCounts[key][d];
                        if (neighbourCount === undefined) {
                            neighbourCount = 0;
                            neighbourCounts[key][d] = neighbourCount;
                        }
                        // increment neighbour count
                        neighbourCount++;
                        neighbourCounts[key][d] = neighbourCount;
                        // get d links
                        var linksById = d3Edge[key][d];
                        if (linksById === undefined) {
                            linksById = {};
                            d3Edge[key][d] = linksById;
                        }
                        // if a link is missing
                        var linksByIdLength = Object.keys(linksById).length;
                        if (linksByIdLength < neighbourCount) {
                            // create it
                            link = {
                                edge: d3Edge,
                                type: 'link',
                                isSource: isSource,
                                source: isSource? d : id,
                                target: (!isSource)? d : id,
                                id: ''+Math.random(), // uuid
                                index: this.links.length, // pos in view
                            };
                            // add in view
                            this.links.push(link);
                            this.shapesById[link.id] = link;
                            // try to inject references to sources
                            var source = this.shapesById[link.source];
                            if (source !== undefined) {
                                link.source = source;
                            }
                            // and targets
                            var target = this.shapesById[link.target];
                            if (target !== undefined) {
                                link.target = target;
                            }
                            // push link in d3Edge model
                            linksById[link.id] = link;
                        }
                    },
                    this
                );
                // remove useless links
                var cleanLink = function(linkId) {
                    var link = this.shapesById[linkId];
                    // delete from model
                    delete d3Edge[key][verticeId][linkId];
                    // delete from view
                    delete this.shapesById[linkId];
                    linkPosToDelete.push(link.index);
                };
                for (var verticeId in d3Edge[key]) {
                    var neighbourCount = neighbourCounts[key][verticeId];
                    var index = 0; // index to start to remove links
                    if (neighbourCount !== undefined) {
                        index = neighbourCount;
                    }
                    // delete links from model
                    var linkIdToDelete = Object.keys(d3Edge[key][verticeId]).splice(index);
                    linkIdToDelete.forEach(
                        cleanLink,
                        this
                    );
                }
            };
            // add links for all sources
            updateLinks.call(this, 'sources');
            // add links for all targets
            updateLinks.call(this, 'targets');
        },

        /**
        * Refresh the view related to selected data.
        */
        refreshSelectedShapes: function() {
            var selected = get(this, 'controller').graphModel.selected;
            // get a list of 'selected' items
            var selectedShapes = this.panel.selectAll('.shapegroup')
                .data(selected, function(d){ return d; });
            // select newly selected items
            selectedShapes.classed('selected', true);
            // unselect old selected
            selectedShapes.exit().classed('selected', false);
        },

        /**
        * Refresh the view related to locked data.
        */
        refreshLockedShapes: function() {
            this.panel.selectAll('.shapegroup').filter(
                function(d) {
                    return d.fixed;
                }
            ).classed('locked', true);
            this.panel.selectAll('.shapegroup').filter(
                function(d) {
                    return !d.fixed;
                }
            ).classed('locked', false);
        },

        coordinates: function() {
            var result = [1, 1];
            if (this.panel !== null) {
                result = d3.mouse(this.panel[0][0]);
            }
            return result;
        },

        /**
        * Show tool box related to input shape.
        */
        showToolBox: function(shape) {
            if (d3.event.defaultPrevented) {
                return;
            }
            var data = shape.__data__;
            // if toolbox already exists, destroy it.
            //if (this.toolbox !== null) {
            //    this.destroyToolBox(shape);
            //}
            // add a new toolbox with specific node toolbox items
            var toolboxItems = this.getToolBoxItems(data);
            // get coordinates
            var coordinates = data? [data.x, data.y] : this.coordinates();
            var me = this;
            // create generic toolbox
            this.toolbox = this.panel.selectAll('.toolbox')
                .data(toolboxItems, function(d) {return d.name;});
            // add new toolbox items
            var newToolbox = this.toolbox.enter()
                .append('g')
                    .style("opacity", 0) // generic style
                    .attr(
                        {
                            id: function(d) { return 'toolbox_' + d.name;}, // id
                            class: function(d) { return 'toolbox ' + d.name;}, // class
                            transform: function(d, i) {
                                var length = 2 * i * Math.PI / toolboxItems.length;
                                var x = coordinates[0] + Math.cos(length) * me.toolboxGap;
                                var y = coordinates[1] + Math.sin(length) * me.toolboxGap;
                                var translate = 'translate(' + x + ',' + y + ')';
                                var scale = 'scale(' + (1 / me.scale) + ')';
                                return translate + scale + (d.transform ? d.transform : '');
                            }
                        }
                    )
                ;
            newToolbox // add shape
                .append('path')
                    .attr(
                        {
                            d: function(d, i) {
                                return d.symbol.size(20 * me._unit)(d, i);
                            }
                        }
                    )
                    // add text
                    .append('title')
                        .text(function(d) {return d.name;})
                ;
            newToolbox // add hyperlink with name
                .append('text')
                    .text(function(d) { return d.name;});
            // move all toolbox to mouse coordinates
            this.toolbox
            // attach handlers
                .on('click', function(d){
                    me[d.name + 'Handler'](data, shape);
                    me.destroyToolBox();
                })
                .transition()
                    .style("opacity", 1) // show them
                    .attr( // and move them
                        {
                            transform: function(d, i) {
                                var length = 2 * i * Math.PI / toolboxItems.length;
                                var x = coordinates[0] + Math.cos(length) * me.toolboxGap;
                                var y = coordinates[1] + Math.sin(length) * me.toolboxGap;
                                return "translate(" + x + "," + y + ") " + (d.transform ? d.transform : '');
                            }
                        }
                    )
                ;
            // delete old toolbox
            this.toolbox.exit()
                .transition()
                    .duration(this.duration)
                    .style("opacity", 0) // generic exit
                    .remove()
                ;
        },

        /**
        * Destroy toolbox.
        */
        destroyToolBox: function() {
            // delete old toolbox
            if (this.toolbox) {
                this.toolbox
                    .transition()
                        .duration(this.duration)
                        .style("opacity", 0)
                        .remove()
                    ;
            }
        },

        closeHandler: function(data) {
            d3.event.stopPropagation();
        },
        editHandler: function(data) {
            d3.event.stopPropagation();
            get(this, 'controller').editRecord(data.id);
        },
        deleteHandler: function(data) {
            d3.event.stopPropagation();
            // ensure data is an array
            if (!Array.isArray(data)) {
                data = [data];
            }
            // save records to delete
            var recordsToDel = [];
            var controller = get(this, 'controller');
            var recordsById = controller.graphModel.recordsById;
            data.forEach(
                function(d) {
                    var record = recordsById[d.id];
                    // delete links
                    if (d.isSource !== undefined) {
                        var dest = d.isSource ? 'source' : 'target';
                        var dests = record.get(dest + 's');
                        var index = dests.indexOf(d[dest]);
                        dests.splice(index, 1);
                        // update the record
                        record.set(dest+'s', dests);
                    } else {  // push record in records to delete
                        recordsToDel.push(record);
                    }
                }
            );
            if (recordsToDel.length > 0) {
                controller.deleteRecords(recordsToDel);
            }
        },
        linkHandler: function(data, shape) {
            d3.event.stopPropagation();
            this.source = data; // save link source
            //var coordinates = [data.x, data.y];
            /*this.tmpLink = this.panel
                .append('line')
                    .attr(
                        {
                            id: 'tmpLink',
                            x1: coordinates[0],
                            y1: coordinates[1],
                            x2: coordinates[0],
                            y2: coordinates[1]
                        }
                    )
                    .classed('tmpLink', true)
                ;*/
        },
        addHandler: function(data) {
            d3.event.stopPropagation();
            var controller = get(this, 'controller');
            var callback = function(record) {
                this.getNode(record);
                controller.trigger('refresh');
            };
            if (this.source === null) { // in case of node
                var record = controller.newRecord(
                    controller.verticeEltType, undefined, true, callback, undefined, this
                );
            } else { // in case of edge
                var success = function(record) {
                    this.removeTmpLink();
                    controller.trigger('refresh');
                };
                var failure = function(record) {
                    this.removeTmpLink();
                };
                this.addLink(this.source, data, true, success, failure, this);
            }
        },
        unselectHandler: function(data) {
            d3.event.stopPropagation();
            get(this, 'controller').unselect(data.id);
        },
        selectHandler: function(data) {
            d3.event.stopPropagation();
            get(this, 'controller').select(data.id);
        },
        unlockHandler: function(data) {
            d3.event.stopPropagation();
            this.lock(data);
        },
        lockHandler: function(data) {
            d3.event.stopPropagation();
            this.lock(data, true);
        },
        cancelHandler: function(data) {
            d3.event.stopPropagation();
            this.removeTmpLink();
        },
        eventpoolHandler: function(data) {
            document.location.href = "/eventpool/?" + data.elt.get('info').entity;
        },

        /**
        * Remove tmp link
        */
        removeTmpLink: function() {
            this.source = null;
            if (this.tmpLink) {
                this.tmpLink.remove();
            }
            this.tmpLink = null;
        },

        /**
        * Lock shapes related to input data.
        * @param data data to lock.
        * @param enable if true, lock input data shapes.
        */
        lock: function(data, enable) {
            if (!Array.isArray(data)) {
                data = [data];
            }
            data.forEach(
                function(d) {
                    d.fixed = enable;
                }
            );
            this.refreshLockedShapes();
        },

        /**
        * called when adding nodes.
        */
        addNodes: function(nodes) {
            var me = this;
            var recordsById = get(this, 'controller').graphModel.recordsById;
            // create the graphical element
            var shapes = nodes
                .append('g')
                    .classed(
                        {
                            shapegroup: true,
                            node: true,
                            edge: function(d) {
                                var record = recordsById[d.id];
                                return record.get('_type') === 'edge';
                            },
                            graph: function(d) {
                                var record = recordsById[d.id];
                                return record.get('_type') === 'graph';
                            }
                        }
                    )
                    .attr(
                        {
                            id: function(d) {return d.id;}
                        }
                    )
                    .on('click', function() {return me.clickAction(this);})
                    .on('dblclick', function(){return me.dblClickAction(this);})  // add menu selection
                    .on('mouseover', function() {me.overAction(this);})
                    .on('mouseout', function() {me.outAction(this);})
                    .on('mousemove', function() {me.moveAction(this);})
                ;
            shapes
                .append('path')
                    .classed(
                        {
                            shape: true
                        }
                    )
                ;
            var nodeShapes = shapes.filter(function(d) {
                var record = recordsById[d.id];
                return record.get('_type') !== 'edge';
            });
            nodeShapes.append('title');  // add title
            nodeShapes.append('text').classed('entity', true);  // add entity name text
            nodeShapes.append('text').classed('operator', true);  // add operator name text
            var edgeShapes = shapes.filter(function(d) {
                var record = recordsById[d.id];
                return record.get('_type') === 'edge';
            });
            edgeShapes.append('text').classed('weight', true); // add weight text
            // node drag
            var node_drag = this._layout.engine.drag
                .on(
                    'dragstart',
                    function (d) {
                        d3.event.sourceEvent.stopPropagation();
                        me.dragging = true;
                        me.panel.select(this).classed("dragging", true);
                    }
                )
                .on(
                    'dragend',
                    function (d) {
                        me.dragging = false;
                        me.panel.select(this).classed("dragging", false);
                    }
                );
            // ensure dragging of circle
            shapes.call(node_drag);
        },

        /**
        * called during node deletion.
        */
        delNodes: function(nodes) {
            nodes
                .transition()
                    .duration(this.duration)
                    .style('opacity', 0)
                        .remove()
            ;
        },

        /**
        * called during node updating.
        */
        updateNodes: function(nodes) {
            var me = this;
            var recordsById = get(this, 'controller').graphModel.recordsById;
            var graphId = get(this, 'controller.model.graph_id');
            nodes
                .select('path')
                    .classed(
                        {
                            'ndok': function(d) {
                                var record = recordsById[d.id];
                                return (!record.get('info').state) || record.get('info').state == 0;
                            },
                            'ndminor': function(d) {
                                var record = recordsById[d.id];
                                return record.get('info').state == 1;
                            },
                            'ndmajor': function(d) {
                                var record = recordsById[d.id];
                                return record.get('info').state == 2;
                            },
                            'ndcritical': function(d) {
                                var record = recordsById[d.id];
                                return record.get('info').state == 3;
                            }
                        }
                    )
                    .attr(
                        'd',
                        function(d, i) {
                            var f = d3.svg.symbol();
                            if (d.x === undefined) {
                                d.x = 0; d.y = 0;
                            }
                            var record = recordsById[d.id];
                            switch(record.get('_type')) {
                                case 'graph':
                                    f = f.type('circle');
                                    if (d.id === graphId) {
                                        f = f.size(30 * me._unit);
                                    } else {
                                        f = f.size((d._weight>=1? d._weight : 1) * me._unit);
                                    }
                                    break;
                                case 'edge': f = f.type('diamond').size(d._weight * me._unit); break;
                                case 'vertice': f = f.type('square').size((d._weight>=1? d._weight : 1) * me._unit); break;
                                default: f = f.type('triangle-down');
                            }
                           var result = f(d, i);
                           return result;
                        }
                    )
            ;
            // update a title which is the entity id or graph element id
            nodes.select('title').text(
                function(d) {
                    var result = '';
                    var record = recordsById[d.id];
                    var info = record.get('info');
                    if (info !== undefined) {
                        var entity = info.entity;
                        if (entity !== undefined) {
                            result += entity;
                        }
                    }
                    return result;
                }
            );
            // update an entity id which is the entity id or graph element id
            nodes.select('.entity').text(
                function (d) {
                    var result = '';
                    var record = recordsById[d.id];
                    var info = record.get('info');
                    if (info && info.label) {
                        result += info.label;
                    } else {
                        var entity = d.entity;
                        if (entity !== undefined) {
                            result += entity.name;
                        }
                    }
                    return result;
                }
            );
            nodes.select('.operator').text(
                function(d) {
                    var result = '';
                    var record = recordsById[d.id];
                    var info = record.get('info');
                    if (info !== undefined) {
                        var operator = info.task;
                        if (operator) {
                            var operatorId = operator.id || operator;
                            if (operatorId === 'canopsis.task.condition.condition') {
                                operator = operator.params.condition;
                                operatorId = operator.id || operator;
                            }
                            operatorId = operatorId.substring(operatorId.lastIndexOf('.') + 1);
                            result += operatorId;
                        }
                        if (result === 'change_state') {  // do not display default operator
                            result = '';
                        }
                    }
                    return result;
                }
            );
            nodes.select('.weight').text(
                function(d) {
                    var result = '';
                    var record = recordsById[d.id];
                    var weight = record.get('weight');
                    if (weight !== undefined && weight !== 1) {
                        result += weight;
                    }
                    return result;
                }
            );
        },

        /**
        * called during link addition.
        */
        addLinks: function(links) {
            var me = this;
            // create the graphical element
            var shapes = links
                .append('path') // line representation
                    .classed(
                        {
                            shape: true, // set shape link
                            link: true, // set class link
                            targetLink: function(d) { return !d.isSource; }
                        }
                    )
                    .attr(
                        {
                            'id': function(d) { return d.id; },
                        }
                    )
                    .on('click', function() { return me.clickAction(this);})
                    .on('mouseout', function() { me.outAction(this);})
                    .on('mousemove', function() { me.moveAction(this);})
                    .on('mouseover', function() { me.overAction(this);})
            ;
        },

        /**
        * called during link deletion.
        */
        delLinks: function(links) {
            links
                .transition()
                    .duration(this.duration)
                    .style('opacity', 0)
                        .remove()
            ;
        },

        /**
        * called during link updating.
        */
        updateLinks: function(links) {
            var me = this;
            var recordsById = get(this, 'controller').graphModel.recordsById;
            links
                .classed(
                    {
                        directed: function(d) {
                            var record = recordsById[d.edge.id];
                            return record.get('directed');
                        },
                        'lnok': function(d) {
                            var record = recordsById[d.edge.id];
                            return (!record.get('info').state) || record.get('info').state == 0;
                        },
                        'lnminor': function(d) {
                            var record = recordsById[d.edge.id];
                            return record.get('info').state == 1;
                        },
                        'lnmajor': function(d) {
                            var record = recordsById[d.edge.id];
                            return record.get('info').state == 2;
                        },
                        'lncritical': function(d) {
                            var record = recordsById[d.edge.id];
                            return record.get('info').state == 3;
                        }
                    }
                )
                .style(
                    {
                        'stroke-width': function(d) {
                            var result = 1;
                            var record = recordsById[d.edge.id];
                            var weight = record.get('weight');
                            if (weight !== undefined) {
                                result += (weight - 1);
                            }
                            return result;
                        },
                        'marker-end': function(d) {
                            var record = recordsById[d.edge.id];
                            return (record.get('directed') && !d.isSource)? "url(#markerArrow)" : "";
                        }
                    }
                )
            ;
        },

        moveAction: function(shape) {
            if (d3.event.defaultPrevented) return;
            this.overAction(shape);
        },

        overAction: function(shape) {
            if (d3.event.defaultPrevented) return;
            var data = shape.__data__;
            if (data !== undefined) {
                d3Shape = d3.select(shape);
                d3Shape.classed('over', true);
            }
            if (this.source !== null) {
                console.log(data);
                var coordinates = data? [data.x, data.y] : this.coordinates();
                /*this.tmpLink.attr(
                    {
                        x2: coordinates[0],
                        y2: coordinates[1]
                    }
                )
                if (this.checkTargetLink(shape.__data__)) {
                    this.tmpLink.attr(
                        {
                            id: 'tmpLink'
                        }
                    )
                } else {
                    this.tmpLink.attr(
                        {
                            id: 'wrongTmpLink'
                        }
                    )
                }*/
            }
        },

        /**
        * Check if a shape is a good target for a new link.
        *
        * @param shape target shape.
        */
        checkTargetLink: function(data) {
            var recordsById = get(this, 'controller').graphModel.recordsById;
            var result = data === undefined;
            if (!result) {
                var record = recordsById[data.id || data];
                result = record.get('_type') !== 'edge' && data.id !== this.source.id;
            }
            return result;
        },

        outAction: function(shape) {
            if (d3.event.defaultPrevented) return;
            d3.select(shape).classed('over', false);
        },

        dblClickAction: function(shape) {
            if (d3.event.defaultPrevented) return;
            d3.event.stopPropagation();
            get(this, 'controller').editRecord(shape.__data__.id);
            this.destroyToolBox();
        },

        clickAction: function(shape) {
            if (d3.event.defaultPrevented) return;
            d3.event.stopPropagation();
            this.showToolBox(shape);
        },

        /**
        * Get a node from a record and register it in this.shapesById and nodes.
        * If node does not exist, create it.
        *
        * @param record source record.
        * @return record node.
        */
        getNode: function(record) {
            // get record id
            var recordId = record.get('id');
            // get result in this memory
            var result = this.shapesById[recordId];
            // get graphId
            var graphId = get(this, 'controller.model.graph_id');
            // get node from db
            var recordNode;
            var view = record.get('info').view;
            if (view !== undefined) {
                recordNode = view[graphId];
            }
            // create a new node if it does not exist
            if (result === undefined) {
                // if old node does not exist
                if (recordNode === undefined) {
                    // apply mouse coordinates to the result
                    result = {
                        fixed: false, // new node is fixed
                        hidden: false, // and displayed,
                        id: recordId, // with record id
                        index: this.nodes.length,
                        type: record.get('type'),
                        x: 0,
                        y: 0,
                        px: 0,
                        py: 0
                    };
                    // add node in nodes
                    this.nodes.push(result);
                } else { // result becomes old node
                    result = recordNode;
                }
                // update this shapesById
                this.shapesById[recordId] = result;
            }
            return result;
        },

        /**
        * Add a link.
        *
        * @param source source link.
        * @param target optional target link.
        * @param edit enable link properties edition.
        * @param success edition success callback. Takes record in parameter.
        * @param failure edition failure callback. Takes record in parameter.
        * @param callback execution context.
        * @return new edge node.
        */
        addLink: function(source, target, edit, success, failure, context) {
            var me = this;
            var result = null;
            var controller = get(this, 'controller');
            var sourceRecord = controller.graphModel.recordsById[source.id];
            var sourceType = sourceRecord.get('_type');
            var graphId = get(controller, 'model.graph_id');
            // ensure source and target are ok
            if (source.id === graphId || !this.checkTargetLink(target)) {
                throw 'Wrong parameters';
            }
            // default failure function
            var _failure = function(reason) {
                console.error(reason);
                if (failure !== undefined) {
                    failure.call(context, reason);
                }
            };
            // get a target
            if (target === undefined) {
                // create a callback
                var coordinates = this.coordinates();
                var _success = function(record) {
                    var target = this.getNode(record);
                    target.px = target.x = coordinates[0];
                    target.py = target.y = coordinates[1];
                    this.addLink(source, target, edit, success, failure, context);
                };
                // edit a new vertice if target does not exist
                target = controller.newRecord(controller.verticeEltType, undefined, edit, _success, _failure, this);
                if (edit) return;
            }
            // get result
            // if source is an edge, then edge is source
            if (sourceType === 'edge') {
                result = source;
                // and add right now the target
                record = controller.graphModel.recordsById[result.id];
                var targets = record.get('targets');
                targets.push(target.id);
                record.set('targets', targets);
                record.save().then(
                    function() {
                        me.weaveLinks(result);
                    },
                    function(reason) {
                        console.error(reason);
                    }
                );
            } else { // else create a new edge
                var callback2 = function(record) {
                    var edge = this.getNode(record);
                    this.weaveLinks(edge);
                    if (success !== undefined) {
                        success.call(context, record);
                    }
                };
                var edge = controller.newRecord(
                    controller.edgeEltType,
                    {
                        sources: [source.id],
                        targets: [target.id]
                    },
                    edit,
                    callback2,
                    _failure,
                    this
                );
            }
            return result;
        },

        /**
        * Create a new toolbox item with a name, a symbol and a transformation.
        * @param name toolbox item name
        * @param symbol see d3.svg.symbol().type().
        * @param transform transformation attribute.
        */
        newToolBoxItem: function(name, symbol, transform) {
            var result = {
                name: name,
                symbol: d3.svg.symbol().type(symbol),
                transform: transform,
                type: 'toolbox'
            };
            return result;
        },

        /**
        * Get toolbox names by node.
        *
        */
        getToolBoxItems: function(data) {
            var result = [];
            var graphId = get(this, 'controller.model.graph_id');
            if (this.source !== null) {
                result.push(
                    this.newToolBoxItem('close', 'triangle-down'),
                    this.newToolBoxItem('cancel', 'cross', 'rotate(45)'),
                    this.newToolBoxItem('add', 'cross')
                );
            } else {
                // default result
                result.push(
                    this.newToolBoxItem('close', 'triangle-down')
                );

                if (data !== undefined && data.type !== 'link') { // is it a node
                    var record = get(this, 'controller').graphModel.recordsById[data.id];
                    var info = record.get('info');
                    if (info) {
                        if (info.entity) {
                            result.push(this.newToolBoxItem('eventpool', 'triangle-up'));
                        }
                    }
                    if (get(this, 'controller').graphModel.selected[data.id] !== undefined) {
                        result.push(this.newToolBoxItem('unselect', 'diamond'));
                    } else {
                        result.push(this.newToolBoxItem('select', 'diamond'));
                    }
                    if (data.fixed) {
                        result.push(this.newToolBoxItem('unlock'));
                    } else {
                        result.push(this.newToolBoxItem('lock'));
                    }
                    if (data.id !== graphId) { // all nodes but topology
                        result.push(
                            this.newToolBoxItem('link', 'triangle-up'), // new link
                            this.newToolBoxItem('delete', 'cross', 'rotate(45)') // elt deletion
                        );
                    }
                    //if (data.elt.get('type') !== 'topology') { // all elts but topologies
                        result.push(this.newToolBoxItem('edit', 'square')); // edit elt
                    //}
                } else {
                    result.push(
                        this.newToolBoxItem('add', 'cross') // add elt
                    );
                }
            }

            return result;
        },

        /**
        * Show only d3 elements where entity names correspond to input filter.
        * @param filter string which specifies node regex entity name to show. If undefined, show all d3 nodes.
        */
        /*filter: function(filter) {
            if (filter === undefined) {
                filter = '*';
            }
            var toShow = this.panel.selectAll('.node').where(function(d) {
                var result = d.entity && d.entity.name;
            });
        },*/

        actions: {
            /**
            * Change boolean value of controller.showProperties with the opposite.
            */
            showHideProperties: function() {
                var showProperties = get(this, 'showProperties');
                set(this, 'showProperties', !showProperties);
            },
            /**
            * Apply filter.
            */
            filter: function() {
                var filter = get(this, 'filter');
                this.filter(filter);
                set(this, 'filter', filter);
            }
        },
    });
    application.register('mixin:topology-view', TopologyViewMixin);
}});
