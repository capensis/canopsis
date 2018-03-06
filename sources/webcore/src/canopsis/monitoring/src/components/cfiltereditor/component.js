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

Ember.Application.initializer({
    name: 'component-cfiltereditor',
    after: 'IndexesRegistry',
    initialize: function(container, application) {
        var Canopsis = {};
        var indexesregistry = container.lookupFactory('registry:indexes');
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;


        var component = Ember.Component.extend({
            init:function() {
                var cfilter_serialized = get(this, 'cfilter_serialized');

                //FIXME Canopsis object is not accessible anymore
                set(this, 'onlyAllowRegisteredIndexes', get(Canopsis, 'conf.frontendConfig.cfilter_allow_only_optimized_filters'));

                if(get(this, 'content') !== null && get(this, 'content') !== undefined) {
                    set(this, 'cfilter_serialized', get(this, 'content'));
                } else if(cfilter_serialized === undefined || cfilter_serialized === null) {
                    set(this, 'cfilter_serialized', '{}');
                }

                this._super.apply(this, arguments);
            },

            indexesTree: {
                'component': {
                    '_metas': {
                        'name': 'Component'
                    },
                    'resource': {
                        '_metas': {
                            'name': 'Resource',
                            'final':true
                        }
                    }
                },
                'connector': {
                    '_metas': {
                        'name': 'connector'
                    },
                    'component': {
                        '_metas': {
                            'name': 'Component'
                        },
                        'resource': {
                            '_metas': {
                                'name': 'Resource',
                                'final':true
                            }
                        }
                    }
                }
            },

            currentClauseIndex: -1,

            cfilter_serialized : Ember.computed.alias('content'),
            viewTabColumns: [{ name:'component', title:'component' }, { name:'resource', title:'resource' }],

            indexes : indexesregistry,
            selectedIndexName : 'event',

            selectedIndexChanged: function () {
                var selectedIndexName = get(this, 'selectedIndexName');
                var indexes = get(this, 'indexes');

                for (var i = 0, l = indexes.all.length; i < l; i++) {
                    var currentIndex = indexes.all[i];
                    if(currentIndex.name === selectedIndexName) {
                        set(this, 'indexesTree', currentIndex.tree);
                    }
                }
            }.observes('selectedIndexName'),

            clauses: function() {
                var cfilter_serialized = get(this, 'cfilter_serialized');
                var clauses = Ember.A();
                var mfilter;

                try {
                    mfilter = JSON.parse(cfilter_serialized);
                } catch (e) {
                    console.error('unable to parse serialized filter');
                    mfilter = { '$or': {}};
                }

                console.log('deserializeCfilter', cfilter_serialized, clauses.length);
                console.log('mfilter', mfilter);

                if(mfilter.$or === undefined && Object.keys(mfilter).length !== 0) {
                    console.error('This editor cannot display the current filter');
                    set(this, 'filterError', 'This editor cannot display the current filter');
                    return clauses;
                }
                if(mfilter.$or === undefined) {
                    return clauses;
                }

                for (var i = 0, l = mfilter.$or.length; i < l; i++) {
                    var currentMfilterOr = mfilter.$or[i];
                    var currentOr = Ember.Object.create({
                        and: Ember.A()
                    });

                    if(currentMfilterOr.$and === undefined && Ember.keys(currentMfilterOr)[0] !== undefined) {
                        currentMfilterOr.$and = [currentMfilterOr];
                    }

                    if (currentMfilterOr.$and !== undefined) {
                        var currentField;
                        for (var j = 0, lj = currentMfilterOr.$and.length; j < lj; j++) {

                            var currentMfilterAnd = currentMfilterOr.$and[j];

                            var clauseKey = Ember.keys(currentMfilterAnd)[0];
                            var clauseOperator = Ember.keys(currentMfilterAnd[clauseKey])[0];
                            console.log(currentMfilterAnd[clauseKey][clauseOperator]);
                            var clauseValue = currentMfilterAnd[clauseKey][clauseOperator];

                            //deserialize in array value to a string with comma separator
                            if ((clauseOperator === 'not in' || clauseOperator === 'in') && typeof clauseValue === 'object') {
                                clauseValue = clauseValue.join(',');
                            }

                            var keys = this.getIndexesForNewAndClause(currentOr);

                            var field = {
                                isFirst : (j === 0),
                                keyId: get(this, 'cfilterEditId') + '-keys-' + (j + 1),
                                options: {
                                    'available_indexes' : keys
                                },
                                key: clauseKey,
                                value: clauseValue,
                                operator: this.getLabelForMongoOperator(clauseOperator),
                                lastOfList : true,
                                filling: false,
                                finalized: true
                            };

                            currentField = Ember.Object.create(field);

                            console.log('field:', currentField);
                            // currentAnd.pushObject(currentField);
                            currentOr.and.pushObject(currentField);
                        }

                        var useIndexes = get(this, 'onlyAllowRegisteredIndexes') === true;

                        if(!useIndexes || get(currentField, 'options.available_indexes.length') > 0) {
                            console.log ('append empty and clause');
                            this.pushEmptyClause(currentOr);
                        }

                        clauses.pushObject(currentOr);
                    }

                    if(get(this, 'onlyAllowRegisteredIndexes') === false) {
                        this.pushEmptyClause(currentOr);
                    }
                }
                console.log('clause deserialized', clauses);

                return clauses;

            }.property(),

            classNames: ['cfilter'],

            operators: [
                {
                    label: '=',
                    value: '$eq'
                },
                {
                    label: '!=',
                    value: '$ne'
                },
                {
                    label: '<',
                    value: '$lt'
                },
                {
                    label: '>',
                    value: '$gt'
                },
                {
                    label: 'in',
                    value: '$in'
                },
                {
                    label: 'not in',
                    value: '$nin'
                },
                {
                    label: 'regex',
                    value: '$regex'
                },
                {
                    label: '!regex',
                    value: '$not'
                }
            ],

            orButtonHidden: false,

            //initialized in this class' constructor
            onlyAllowRegisteredIndexes: true,

            serializeCfilter: function() {
                var clauses = get(this, 'clauses');

                var mfilter = {
                    '$or': []
                };


                for (var i = 0, l_clauses = clauses.length; i < l_clauses; i++) {
                    var clause = clauses[i];

                    var subfilter = {
                        '$and': []
                    };

                    if (clause.and[0] !== undefined) {
                        set(clause.and[0], 'isFirst', true);
                    }

                    for (var j = 0, l_and = clause.and.length; j < l_and; j++) {
                        var field = clause.and[j];

                        if(j === 0) {
                            set(field, 'isFirst', true);
                        } else {
                            set(field, 'isFirst', false);
                        }

                        // if(get(this, 'onlyAllowRegisteredIndexes') === true) {
                        set(field, 'finalized', true);
                        // }

                        if (j === clause.and.length -1) {
                            set(clause.and[j], 'isLast', true);
                        } else {
                            set(clause.and[j], 'isLast', false);
                        }

                        if (!Ember.isArray(clause.and[j].value) && (clause.and[j].operator === 'in' || clause.and[j].operator === 'not in')) {
                            console.log('Operator in detected');
                            clause.and[j].value = clause.and[j].value.split(',');
                        } else {
                            //manage numbers inputs and cast them to number if numeric.
                            if (typeof clause.and[j].value === 'string' && $.isNumeric(clause.and[j].value)) {
                                clause.and[j].value = parseFloat(clause.and[j].value);
                            }
                        }

                        if (field.key !== undefined) {
                            var item = {};
                            console.log('field', field);
                            var operator = {
                                label: '=',
                                value: '$eq'
                            };

                            if (field.operator !== undefined) {
                                operator = this.getMongoOperatorForLabel(field.operator);
                            }

                            if(isNone(operator.format)) {
                                operator.format = function(x) {
                                    return x;
                                };
                            }

                            item[field.key] = {};
                            item[field.key][operator.value] = operator.format(field.value);

                            subfilter.$and.pushObject(item);
                        }

                    }

                    if (subfilter.$and.length > 0) {
                        if (subfilter.$and.length === 1) {
                            subfilter = subfilter.$and[0];
                        }

                        mfilter.$or.pushObject(subfilter);
                    }
                }

                if(mfilter.$or.length === 0) {
                    mfilter = {};
                }

                mfilter = JSON.stringify(mfilter, null, '');
                return mfilter;
            },

            getMongoOperatorForLabel: function(label) {
                for (var i = 0, l = this.operators.length; i < l; i++) {
                    if (this.operators[i].label === label) {
                        return this.operators[i];
                    }
                }

                return undefined;
            },

            getLabelForMongoOperator: function(mongoOperator) {
                for (var i = 0, l = this.operators.length; i < l; i++) {
                    if (this.operators[i].value === mongoOperator) {
                        return this.operators[i].label;
                    }
                }

                return undefined;
            },

            checkIfNewAndClauseDisplayed : function() {

                var currentClauseIndex = get(this, 'currentClauseIndex');

                if (currentClauseIndex >= 0) {
                    var clauses = get(this, 'clauses');
                    var currentClause = clauses.objectAt(currentClauseIndex);

                    var lastAndOfClause = currentClause.and[currentClause.and.length - 1];

                    var isEmpty = function(value) {
                        if (value === undefined || value === '') {
                            return true;
                        } else {
                            return false;
                        }
                    };

                    if (lastAndOfClause !== undefined && isEmpty(lastAndOfClause.key) && isEmpty(lastAndOfClause.value)) {
                        set(this, 'newAndClauseDisplayed', false);
                        return;
                    }

                    set(this, 'newAndClauseDisplayed', true);

                    return;
                } else {
                    set(this, 'newAndClauseDisplayed', false);
                    return;
                }
            }.observes('currentClauseIndex'),

            clausesChanged: function() {
                var clauses = get(this, 'clauses');

                console.log('clausesChanged', clauses, clauses.length);

                //detect if we have to display the addOrClause button
                if (clauses.length === 0) {
                    set(this, 'orButtonHidden', false);
                } else {
                    var lastOrClause = clauses[clauses.length -1];
                    console.log('last and length', clauses);
                    console.log('last and length', lastOrClause);
                    console.log('last and length', lastOrClause.and.length);
                    var lastAndQueryPart = lastOrClause.and[lastOrClause.and.length -1];
                    console.log('lastAndQueryPart', lastAndQueryPart, lastAndQueryPart.key);
                    if (lastOrClause.and.length <= 1 && (Ember.isNone(lastAndQueryPart.key) || Ember.isNone(lastAndQueryPart.value))) {
                        set(this, 'orButtonHidden', true);
                    } else {
                        set(this, 'orButtonHidden', false);
                    }
                }

                var mfilter = this.serializeCfilter();
                console.log('generated mfilter', mfilter);
                set(this, 'cfilter_serialized', mfilter);
            },

            cfilterId: function() {
                return get(this, 'elementId') + '-cfilter';
            }.property('elementId'),

            cfilter: function() {
                return $('#' + get(this, 'cfilterId'));
            },

            cfilterEditId: function() {
                return get(this, 'cfilterId') + '-edit';
            }.property('cfilterId'),

            cfilterEditTabId: function() {
                return '#' + get(this, 'cfilterEditId');
            }.property('cfilterEditId'),

            cfilterEdit: function() {
                return $(get(this, 'cfilterEditTabId'));
            },

            cfilterRawId: function() {
                return get(this, 'cfilterId') + '-raw';
            }.property('cfilterId'),

            cfilterRawTabId: function() {
                return '#' + get(this, 'cfilterRawId');
            }.property('cfilterRawId'),

            cfilterRaw: function() {
                return $(get(this, 'cfilterRawTabId'));
            },

            cfilterViewId: function() {
                return get(this, 'cfilterId') + '-view';
            }.property('cfilterId'),

            cfilterViewTabId: function() {
                return '#' + get(this, 'cfilterViewId');
            }.property('cfilterViewId'),

            cfilterView: function() {
                return $(get(this, 'cfilterViewTabId'));
            },

            getIndexesForNewAndClause: function(currentClause) {
                console.group('getIndexesForNewAndClause useIndexesOptions : ', get(this, 'onlyAllowRegisteredIndexes'), currentClause);

                if(get(this, 'onlyAllowRegisteredIndexes') === true) {

                    console.group('getIndexesForNewAndClause', currentClause);

                    var indexesTreeCursor = get(this, 'indexesTree');

                    for (var i = 0, l = currentClause.and.length; i < l; i++) {
                        var currentAnd = currentClause.and[i];
                        console.log('currentAnd', currentAnd);
                        if (indexesTreeCursor === undefined) {
                            console.error('bad index management', currentAnd.key);
                        } else {
                            indexesTreeCursor = indexesTreeCursor[currentAnd.key];
                        }
                        console.log('indexesTreeCursor', indexesTreeCursor);
                    }

                    console.info('available indexes', indexesTreeCursor);

                    var available_indexes = [];

                    for (var key in indexesTreeCursor) {
                        console.log('iter', key, indexesTreeCursor[key]);
                        if (key !== '_metas') {
                            available_indexes.pushObject({name: indexesTreeCursor[key]._metas.name, value: key, _metas: indexesTreeCursor[key]._metas});
                        }
                    }

                    console.groupEnd();
                    return available_indexes;
                }

                console.log('no index available because "use indexes" option is disabled');
                console.groupEnd();
            },

            pushEmptyClause: function(currentClause) {
                console.group('pushEmptyAndClause', get(this, 'onlyAllowRegisteredIndexes'));

                var keys = this.getIndexesForNewAndClause(currentClause);

                var useIndexes = get(this, 'onlyAllowRegisteredIndexes') === true;

                console.log('available_indexes', keys);

                var field = {
                    keyId: get(this, 'cfilterEditId') + '-keys-' + (currentClause.and.length + 1),

                    options: {
                        'available_indexes' : keys
                    },
                    key: undefined,
                    value: undefined,
                    operator: '=',
                    lastOfList : true,
                    filling: false,
                    finalized: false
                };

                var lastAndClauseOfList = currentClause.and[currentClause.and.length - 2];
                console.log('and array', currentClause.and);
                console.log('lastAndClauseOfList', lastAndClauseOfList);

                if (lastAndClauseOfList !== undefined) {
                    set(lastAndClauseOfList, 'lastOfList', false);
                }

                if (!useIndexes || get(field, 'options.available_indexes.length') > 0) {
                    currentClause.and.pushObject(Ember.Object.create(field));
                }

                console.groupEnd('pushEmptyClause');
            },

            actions: {
                selectIndexByName: function (name) {
                    set(this, 'selectedIndexName', name);
                },
                unlockIndexes: function() {
                    set(this, 'onlyAllowRegisteredIndexes', false);
                },
                lockIndexes: function() {
                    set(this, 'onlyAllowRegisteredIndexes', true);
                },
                addAndClause: function(wasFinalized) {
                    console.log('Add AND clause');

                    var clauses = this.get('clauses');
                    var currentClauseIndex = this.get('currentClauseIndex');

                    console.log(currentClauseIndex);
                    var currentClause = clauses.objectAt(currentClauseIndex);

                    if(Ember.isNone(currentClause)) {
                        console.error('currentClause seems empty');
                    }

                    if (currentClauseIndex >= 0) {
                        var useIndexes = get(this, 'onlyAllowRegisteredIndexes');

                        console.log('test if it\'s possible to add a and clause', useIndexes, !wasFinalized);
                        if (useIndexes || !wasFinalized) {
                            console.log('try to push empty and clause');

                            this.pushEmptyClause(currentClause);
                        }
                    }

                    console.log('clauses addAndClause', clauses);
                    set(this, 'clauses', clauses);
                    this.clausesChanged();
                },

                addOrClause: function() {
                    var clauses = get(this, 'clauses');
                    console.group('Add OR clause', clauses);

                    var currentClauseIndex = get(this, 'currentClauseIndex');
                    var currentClause;

                    if (currentClauseIndex >= 0) {
                        currentClause = clauses.objectAt(currentClauseIndex);
                        set(currentClause, 'current', false);
                    }

                    currentClause = clauses.pushObject(
                        Ember.Object.create({
                            current: true,
                            and: []
                        })
                    );

                    set(this, 'currentClauseIndex', clauses.length - 1);

                    console.log('calling addAndClause');
                    this.send('addAndClause');
                    this.send('activate', currentClause);

                    console.groupEnd();
                },

                activate: function(clause) {
                    var clauses = get(this, 'clauses');
                    var currentClauseIndex = get(this, 'currentClauseIndex');

                    var newCurrentClauseIndex = clauses.indexOf(clause);

                    if (currentClauseIndex !== newCurrentClauseIndex) {
                        console.log('Activate clause:', clause);

                        if (currentClauseIndex >= 0) {
                            clauses.objectAt(currentClauseIndex).set('current', false);
                        }

                        set(clause, 'current', true);

                        console.log('changing currentClauseIndex');

                        set(this, 'currentClauseIndex', newCurrentClauseIndex);
                    }
                },

                removeAndClause: function(selectedClause, removedAnd) {
                    console.log('removeAndClause');

                    var currentClause;
                    var clauses = get(this, 'clauses');
                    var eraseSuccessors = false;

                    for (var i = 0; i < selectedClause.and.length; i++) {
                        var currentAnd = selectedClause.and.objectAt(i);
                        console.log('currentAnd', currentAnd);
                        if (eraseSuccessors === true) {
                            selectedClause.and.removeAt(i);
                            i--;
                        }

                        if (currentAnd === removedAnd) {
                            selectedClause.and.removeAt(i);
                            console.log('processing delete', selectedClause.and, selectedClause.and.objectAt(i - 1));
                            var lastClauseOfList = selectedClause.and.objectAt(i - 1);
                            if(lastClauseOfList !== undefined) {
                                set(lastClauseOfList, 'lastOfList', true);
                            }
                            i--;

                            if (get(this, 'onlyAllowRegisteredIndexes') === true) {
                                eraseSuccessors = true;
                                if (i === -1) {

                                    console.log('the clause will be empty, drop it and quit');

                                    var removedClause = selectedClause;

                                    for (var j = 0, lj = clauses.length; j < lj; j++) {
                                        currentClause = clauses[j];
                                        if (currentClause === removedClause) {
                                            clauses.removeAt(j);

                                            if (get(this, 'currentClauseIndex') >= j) {
                                                set(this, 'currentClauseIndex', get(this, 'currentClauseIndex') - 1);
                                            }

                                            this.clausesChanged();
                                            return;
                                        }
                                    }
                                }
                            }
                        }
                    }

                    if (selectedClause !== undefined && get(this, 'onlyAllowRegisteredIndexes') === true) {
                        this.pushEmptyClause(selectedClause);
                    }

                    this.clausesChanged();
                }
            }
        });

        application.register('component:component-cfiltereditor', component);
    }
});
