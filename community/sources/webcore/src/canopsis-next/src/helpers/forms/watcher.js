import { isUndefined, cloneDeep, omit } from 'lodash';

import { CANOPSIS_STACK, ENTITIES_TYPES } from '@/constants';

import uuid from '@/helpers/uuid';

import { filterToForm, formToFilter } from './filter';

/**
 * @typedef {Object} WatcherForm
 * @property {string} _id
 * @property {string} type
 * @property {string} name
 * @property {boolean} enabled
 * @property {FilterForm} mfilter
 * @property {Object} infos
 * @property {Array} impact
 * @property {Array} depends
 * @property {Array} entities
 * @property {string} output_template
 */

/**
 * @typedef {Object} WatcherPython
 * @property {string} _id
 * @property {string} type
 * @property {string} name
 * @property {boolean} enabled
 * @property {string} mfilter
 * @property {Object} infos
 * @property {Array} impact
 * @property {Array} depends
 */

/**
 * @typedef {Object} WatcherGo
 * @property {string} _id
 * @property {string} type
 * @property {string} name
 * @property {boolean} enabled
 * @property {Object} infos
 * @property {Array} entities
 * @property {string} output_template
 */

/**
 * Convert watcher object to watcher form
 *
 * @param {WatcherGo|WatcherPython|Object} [watcher = {}]
 * @returns {WatcherForm}
 */
export function watcherToForm(watcher = {}) {
  return {
    _id: watcher._id || uuid('watcher'),
    type: ENTITIES_TYPES.watcher,
    name: watcher.name || '',
    enabled: !isUndefined(watcher.enabled) ? watcher.enabled : true,
    infos: watcher.infos ? cloneDeep(watcher.infos) : {},

    /**
     * PYTHON stack fields
     */
    mfilter: filterToForm(JSON.parse(watcher.mfilter || '{}')),
    impact: watcher.impact ? [...watcher.impact] : [],
    depends: watcher.depends ? [...watcher.depends] : [],

    /**
     * GO stack fields
     */
    entities: watcher.entities ? cloneDeep(watcher.entities) : [],
    output_template: watcher.output_template || '',
  };
}

/**
 * Convert watcher form to watcher object by stack
 *
 * @param {WatcherForm} [form = {}]
 * @param {string} [stack = CANOPSIS_STACK.go]
 * @returns {WatcherGo|WatcherPython}
 */
export function formToWatcherByStack(form = {}, stack = CANOPSIS_STACK.go) {
  if (stack === CANOPSIS_STACK.go) {
    return {
      ...omit(form, ['mfilter', 'impact', 'depends']),

      state: {
        method: 'worst',
      },
    };
  }


  return {
    ...omit(form, ['entities', 'output_template', 'mfilter']),

    mfilter: JSON.stringify(formToFilter(form.mfilter)),
  };
}
