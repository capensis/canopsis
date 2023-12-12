import store from '@/store';

import { Handlebars } from './handlebars';

/**
 * Compile template
 *
 * @param {string} template
 * @param {Object} [context = {}]
 * @param {Object} [instance = Handlebars]
 * @returns {Promise<string>}
 */
export async function compile(template, context = {}, instance = Handlebars) {
  const handleBarFunction = instance.compile(template ?? '');
  const preparedContext = {
    env: store.getters['templateVars/items'] ?? {},
    user: store.getters['auth/currentUser'] ?? {},

    ...context,
  };

  return handleBarFunction(preparedContext);
}
