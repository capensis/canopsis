import { cloneDeep, omit } from 'lodash';

import { CANOPSIS_STACK, ENTITIES_TYPES } from '@/constants';

import { filterToForm, formToFilter } from './filter';

/**
 * Convert watcher to form
 *
 * @param {Object} [watcher = {}]
 * @returns {Object}
 */
export function watcherToForm(watcher = {}) {
  return {
    name: watcher.name || '',
    mfilter: filterToForm(JSON.parse(watcher.mfilter || '{}')),
    entities: watcher.entities ? cloneDeep(watcher.entities) : [],
    infos: watcher.infos ? cloneDeep(watcher.infos) : {},
    impact: watcher.impact ? [...watcher.impact] : [],
    depends: watcher.depends ? [...watcher.depends] : [],
    output_template: watcher.output_template || '',
  };
}

/**
 * Convert form to watcher depends on stack
 *
 * @param {Object} [form = {}]
 * @param {string} [stack = CANOPSIS_STACK.go]
 * @returns {Object}
 */
export function formToWatcher(form = {}, stack = CANOPSIS_STACK.go) {
  if (stack === CANOPSIS_STACK.go) {
    return {
      ...omit(form, ['mfilter', 'impact', 'depends']),

      type: ENTITIES_TYPES.watcher,
      state: {
        method: 'worst',
      },
    };
  }

  return {
    ...omit(form, ['entities', 'output_template', 'mfilter']),

    type: ENTITIES_TYPES.watcher,
    mfilter: JSON.stringify(formToFilter(form.mfilter)),
  };
}
