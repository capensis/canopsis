/*eslint-disable*/
/*
 * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
 *
 * This file is part of Canopsis.
 *
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Canopsis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
 */

/**
* This widget contains graph records (vertices, edges and graphs).
* They are saved in controller.recordsById dictionary (like record elements).
*
* The process begins in executing the updateModel method.
* This last aims to retrieve a graph with its elements.
* Once the graph is retrieved, all elements including the graph, are
* transformed into schemas records and saved in the dictionary graph._delts where keys are record uid.
* A vertice becomes a node, an edge becomes a set of node and links, and the graph is a node as well.
*/

// define([
//     'canopsis/uibase/widgets/topology/view',
//     'canopsis/uibase/widgets/topology/adapter',
//     'link!canopsis/uibase/widgets/topology/style.css'
// ], function(WidgetFactory, formsUtils, dataUtils, TopologyViewMixin) {

Ember.Application.initializer({
    name: 'TopologyWidget',
    after: ['WidgetFactory', 'FormsUtils', 'DataUtils', 'TopologyViewMixin'],
    initialize: function(container, application) {
        var WidgetFactory = container.lookupFactory('factory:widget');
        var formsUtils = container.lookupFactory('utility:forms');
        var dataUtils = container.lookupFactory('utility:data');
        var TopologyViewMixin = container.lookupFactory('mixin:topology-view');

        var get = Ember.get,
            set = Ember.set;

        var widget = WidgetFactory('topology', {

            viewMixins: [
                TopologyViewMixin
            ],

            graphModel: {
                graph: null, // graph record
                recordsById: {}, // record elements by id,
                selected: {} // selected records
            },

            graph_cls: 'canopsis.topology.elements.Topology', // default graph class

            defaultVerticeCls: 'canopsis.topology.elements.TopoNode', // default vertice class
            defaultEdgeCls: 'canopsis.topology.elements.TopoEdge', // default edge class

            graphEltType: 'topo', // graph elt type
            verticeEltType: 'toponode', // vertice elt type
            edgeEltType: 'topoedge', // edge elt type

            init: function() {
                this._super.apply(this, arguments);
            },

            /**
            * Get entities from a foreign server.
            * @param entityIds list of entity ids from where get entities.
            * @param success handler to request success.
            * @param failure handler to request failure.
            */
            getEntitiesFromServer: function(entityIds, success, failure) {
                if (entityIds.length > 0) {
                    var doAjax = function(resolve, reject) {
                        $.ajax(
                            {
                                url: '/context/ids',
                                type: 'POST',
                                data: {
                                    'ids': JSON.stringify(entityIds),
                                }
                            }
                        ).then(resolve, reject);
                    };
                    var promise = new Ember.RSVP.Promise(
                        doAjax
                    );
                    promise.then(success, failure);
                } else {
                    success({total: 0});
                }
            },

            /**
            * Add a graph and all its related elements into this records.
            *
            * @param graph graph to add with all inner elements.
            */
            _addGraph: function(graph) {
                var me = this;
                var recordsById = this.graphModel.recordsById;
                // register the graph
                this.graphModel.graph = graph;
                recordsById[graph.get('id')] = graph;
                // add all graph elts
                var _delts = graph.get('_delts');
                var addRecord = function(elt_id) {
                    var elt = _delts[elt_id];
                    elt.id = elt_id;
                    var record = me.widgetDataStore.createRecord(
                        elt.type, elt
                    );
                    recordsById[elt_id] = record;
                    // save record in order to bind it to the adapter
                    return record.save();
                };
                // and in doing the server request
                records_to_save = Object.keys(_delts).map(
                    addRecord
                );
                var refresh = function() {
                    // refresh the view
                    me.trigger('refresh');
                };
                Ember.RSVP.all(records_to_save).then(
                    refresh,
                    console.error
                );
            },

            /**
            * display nodes in the widget
            */
            findItems: function() {
                // load business
                var me = this;
                // get graphId
                var graphId = get(this, 'model.graph_id');
                if (graphId !== undefined) {
                    // delete old records in memory
                    if (this.graphModel.graph !== null && this.graphModel.graph.get('id') !== graphId) {
                        this.graphModel.graph = null;
                        this.graphModel.recordsById = {};
                        this.graphModel.selected = {};
                    }
                    var query = {
                        ids: graphId,
                        add_elts: true
                    };
                    // update current graph
                    var updateGraph = function(result) {
                        var graph = null;
                        if (result.content.length === 1) {  // if content exists
                            // get graph and elt ids
                            graph = result.content[0];
                            // if old graph exists
                            if (me.graphModel.graph !== null) {
                                var _delts = graph.get('_delts');
                                // delete old records
                                var recordsById = me.graphModel.recordsById;
                                var selected = me.graphModel.selected;
                                Object.keys(recordsById).forEach(
                                    function(recordId) {
                                        if (recordId !== graphId && _delts[recordId] === undefined) {
                                            delete recordsById[recordId];
                                            delete selected[recordId];
                                        }
                                    }
                                );
                            }
                            me._addGraph(graph);
                        } else {  // if no graph exists
                            console.error('Several graph obtained from a simple request: ' + result.content);
                        }
                    };
                    // if no graph exists, create a new one
                    var newGraph = function(reason) {
                        var graph = me.widgetDataStore.createRecord(
                            me.graphEltType, {id: graphId}
                        );
                        var addGraph = function() {
                            me._addGraph(graph);
                        };
                        graph.save().then(
                            addGraph
                        );
                    };
                    this.widgetDataStore.find('graph', query).then(
                        updateGraph,
                        newGraph
                    );
                }
            },

            /**
            * Delete record(s).
            * @param records record(s) to delete.
            * @param success fired if deleting successed.
            * @param failure fired if deleting failed.
            * @param context success/failure execution context.
            */
            deleteRecords: function(records, success, failure, context) {
                // ensure records is an array of records
                if (! Array.isArray(records)) {
                    records = [records];
                }
                var me = this;
                var graphId = this.graphModel.graph.get('id');
                var recordsById = this.graphModel.recordsById;
                var selected = this.graphModel.selected;
                // create an array of promises
                var destroyRecord = function(record) {
                    var result = null;
                    var recordId = record.get('id');
                    if (recordId === graphId) {
                        console.error('Impossible to delete the graph');
                    } else {
                        delete recordsById[recordId];
                        delete selected[recordId];
                        result = record.destroyRecord();
                    }
                    return result;
                };
                var recordsToDelete = records.map(
                    destroyRecord
                );
                var thenDelete = function() {
                    if (success !== undefined) {
                        success.call(context, arguments);
                    }
                    me.findItems();
                };
                var failDelete = function(reason) {
                    console.error(reason);
                    if (failure !== undefined) {
                        failure.call(context, arguments);
                    }
                    me.trigger('refresh');
                };
                // execute all promises
                Ember.RSVP.all(recordsToDelete).then(
                    thenDelete,
                    failDelete
                );
            },

            /**
            * Select input record(s) in order to get detail informations.
            *
            * @param records record(s) to select.
            */
            select: function(records) {
                if (! Array.isArray(records)) {
                    records = [records];
                }
                var updateSelected = function(record) {
                    this.graphModel.selected[record.get('id')] = record;
                };
                records.forEach(
                    updateSelected,
                    this
                );
                this.trigger('refresh');
            },

            /**
            * Unselect input record(s) in order to get detail informations.
            *
            * @param records record(s) to unselect.
            */
            unselect: function(records) {
                if (! Array.isArray(records)) {
                    records = [records];
                }
                var deleteRecord = function(record) {
                    delete this.graphModel.selected[record.get('id')];
                };
                records.forEach(
                    deleteRecords,
                    this
                );
                this.trigger('refresh');
            },

            /**
            * Save records.
            *
            * @param records records to save.
            */
            saveRecords: function(records, success, failure, context) {
                // ensure records is an array of records
                if (! Array.isArray(records)) {
                    records = [records];
                }
                // save all records in an array of promises
                var promises = records.map(
                    function(record) {
                        var info = record.get('info');
                        // save view elt information in the dictionary of view graphId
                        if (info.view === undefined) {
                            info.view = {};
                        }
                        info.view[get(this, 'model.graph_id')] = record.view;
                        // save the record
                        return record.save();
                    }
                );
                // execute promises
                Ember.RSVP.all(promises).then(success, context).catch(failure, context);
            },

            /**
            * Convert input elt to a record, or returns it if already a record.
            *
            * @param elt element to convert. If undefined, get default vertice element.
            * @param edit edit record with a form if true.
            * @param success edition success callback. Takes record in parameter.
            * @param failure edition failure callback. Takes record in parameter.
            */
            newRecord: function(type, properties, edit, success, failure, context) {
                var me = this;
                // initialize type with default vertice
                if (type === undefined) {
                    type = this.verticeEltType;
                }
                var result = this.widgetDataStore.createRecord(type, properties);
                // any failure would result in calling input failure
                var _failure = function(reason) {
                    console.error(reason);
                    if (failure !== undefined) {
                        failure.call(context, reason);
                    }
                };
                // callback called if the record has been created
                var _success = function(record) {
                    var recordId = record.get('id');
                    var recordsById = this.graphModel.recordsById;
                    var graph = this.graphModel.graph;
                    var oldRecord = recordsById[recordId];
                    // if record already exists in graph
                    if (oldRecord !== undefined) {
                        this.deleteRecords(oldRecord);
                    } else { // add record to the graph
                        var elts = graph.get('elts');
                        elts.push(recordId);
                        graph.set('elts', elts);
                        // update record in self records by id
                        var saved = function() {
                            // update reference of the record
                            recordsById[recordId] = record;
                            if (success !== undefined) {
                                success.call(context, record);
                            }
                        };
                        var failed = function(reason) {
                            var failure = function(reason) {
                                console.error(reason);
                                _failure(reason);
                            };
                            record.destroyRecord().then(
                                _failure
                            );
                        };
                        graph.save().then(
                            saved,
                            failed
                        );
                    }
                };
                if (edit) {
                    this.editRecord(result, _success, _failure, this);
                } else {
                    _success.call(this, result);
                }
                return result;
            },

            /**
            * Edit a record.
            *
            * @param record record to edit. Can be a record id.
            * @param success edition success callback. Takes record in parameter.
            * @param failure edition failure callback. Takes record in parameter.
            */
            editRecord: function(record, success, failure, context) {
                var me = this;
                // ensure record is a record in case of record id
                if (typeof record === 'string') {
                    record = this.graphModel.recordsById[record];
                }
                // fill operator data from record
                var recordType = record.get('_type');
                switch(recordType) {
                    case 'edge':
                        break;
                    case 'vertice':
                    case 'graph':
                        var states = ['ok', 'minor', 'major', 'critical'];
                        /**
                        * Get task id.
                        *
                        * @param  task task from where get task id. Can be a string or a dictionary.
                        */
                        var getShortId = function(task) {
                            var taskId = task.id || task;
                            var lastIndex = taskId.lastIndexOf('.');
                            var result = taskId.substring(lastIndex + 1);
                            return result;
                        };
                        /**
                        * Update form properties related to input params.
                        *
                        * @param params params from where get action params.
                        * @param record form record.
                        * @param actionName action name to retrieve from params.
                        * @param stateName state name to retrieve from params.
                        */
                        var updateAction = function(params, record, actionName, stateName) {
                            // set actionState
                            var action = params[actionName];
                            if (action !== undefined) {
                                var taskId = getShortId(action);
                                var actionState = taskId;
                                if (taskId === 'change_state') {
                                    var actionParams = action.params;
                                    if (actionParams !== undefined) {
                                        actionState = actionParams.state;
                                        actionState = states[actionState];
                                    }
                                }
                                record.set(stateName, actionState);
                            }
                        };
                        /**
                        * Update record info.
                        *
                        * @param record record to update.
                        * @param actionName action name to update in info.
                        * @param action action to retrieve from record task params.
                        */
                        var updateInfoAction = function(task, record, actionName, action) {
                            // set thenState
                            var conState = record.get(actionName);
                            task.params[action] = {
                                id: 'canopsis.topology.rule.action.',
                                params: {}
                            };
                            switch(conState) {
                                case 'worst_state':
                                case 'best_state':
                                    task.params[action].id += conState;
                                    break;
                                default:
                                    conState = states.indexOf(conState);
                                    task.params[action].id += 'change_state';
                                    task.params[action].params.state = conState;
                            }
                        };
                        var info = record.get('info');
                        if (info !== undefined) {
                            // set entity
                            var entity = info.entity;
                            if (entity !== undefined) {
                                record.set('entity', entity);
                            }
                            // set label
                            var label = info.label;
                            if (label !== undefined) {
                                record.set('label', label);
                            }
                            var task = info.task;
                            if (task !== undefined) {
                                var operator = task.id || task;
                                if (operator !== undefined) {
                                    var operatorName = getShortId(task);
                                    if (operatorName === 'condition') {  // at least / nok
                                        var params = task.params;
                                        if (params !== undefined) {
                                            var condition = params.condition;
                                            if (condition !== undefined) {
                                                var conditionName = getShortId(condition.id || condition);
                                                record.set('operator', conditionName);
                                                var condParams = condition.params;
                                                if (condParams !== undefined) {
                                                    // set inState
                                                    if (conditionName === 'nok') {
                                                        record.set('operator', 'nok');
                                                    } else {
                                                        var inState = condParams.state;
                                                        var f = condParams.f;
                                                        if (f === 'canopsis.topology.rule.condition.is_nok') {
                                                            inState = 'nok';
                                                        } else {
                                                            if (inState !== undefined) {
                                                                inState = states[inState];
                                                                record.set('in_state', inState);
                                                            }
                                                        }
                                                    }
                                                    // set minWeight
                                                    var minWeight = condParams.min_weight;
                                                    if (minWeight !== undefined) {
                                                        record.set('min_weight', minWeight);
                                                    }
                                                    // set thenState
                                                    updateAction(params, record, 'statement', 'then_state');
                                                    // set elseState
                                                    updateAction(params, record, '_else', 'else_state');
                                                }
                                            }
                                        }
                                    } else {  // simple task
                                        record.set('operator', operatorName);
                                    }
                                }
                            }
                        }
                        break;
                    default: break;
                }
                var recordWizard = formsUtils.showNew(
                    'modelform',
                    record,
                    {inspectedItemType: record.get('type')}
                );
                /**
                * form execution success.
                */
                var done = function() {
                    switch(recordType) {
                        case 'edge':
                            break;
                        case 'vertice':
                        case 'graph':
                            var info = record.get('info');
                            var task = null;
                            if (info !== undefined) {
                                task = info.task;
                                if (typeof task === 'string') {
                                    info.task = {
                                        id: task
                                    };
                                    task = info.task;
                                } else if (task === undefined) {
                                    info.task = {};
                                    task = info.task;
                                }
                            } else {
                                info = {
                                    task: {}
                                };
                                task = info.task;
                            }
                            // set entity
                            var entity = record.get('entity');
                            if (entity !== undefined) {
                                info.entity = entity;
                            }
                            // set label
                            var label = record.get('label');
                            if (label !== undefined) {
                                info.label = label;
                            }
                            var operator = record.get('operator');
                            switch(operator) {
                                case 'change_state':
                                case 'worst_state':
                                case 'best_state':
                                    task.id = 'canopsis.topology.rule.action.' + operator;
                                    task.params = {};
                                    break;
                                case '_all':
                                case 'at_least':
                                case 'nok':
                                    task.id = 'canopsis.task.condition.condition';
                                    task.params = {};
                                    task.params.condition = {
                                        id: 'canopsis.topology.rule.condition.' + operator,
                                        params: {}
                                    };
                                    // set minWeight
                                    if (operator !== '_all') {
                                        var minWeight = record.get('min_weight');
                                        task.params.condition.params.min_weight = minWeight;
                                    }
                                    // set inState
                                    var inState = record.get('in_state');
                                    if (inState === 'nok') {
                                        if (operator !== 'nok') {
                                            task.params.condition.params.f = 'canopsis.topology.rule.condition.is_nok';
                                            task.params.condition.params.state = null;
                                        }
                                    } else {
                                        task.params.condition.params.state = states.indexOf(inState);
                                    }
                                    // set statement
                                    updateInfoAction(task, record, 'then_state', 'statement');
                                    // set else
                                    updateInfoAction(task, record, 'else_state', '_else');
                                    break;
                                default: break;
                            }
                            // update info
                            record.set('info', info);
                            break;
                        default: break;
                    }
                    var _success = function(record) {
                        var _record = record[0];
                        if (success !== undefined) {
                            success.call(context, _record);
                        }
                        me.trigger('refresh');
                    };
                    var _failure = function(record) {
                        console.error(record);
                        if (failure !== undefined) {
                            failure.call(context, record);
                        }
                    };
                    // save the record
                    me.saveRecords(record, _success, _failure, me);
                };
                /**
                * form execution failure.
                */
                var fail = function() {
                    if (failure !== undefined) {
                        failure.call(context, record);
                    }
                };
                recordWizard.submit.done(
                    done
                ).fail(
                    fail
                );
            },

        });
        application.register('widget:topology', widget);

    }
});
