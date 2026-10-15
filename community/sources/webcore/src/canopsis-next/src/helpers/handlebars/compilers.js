import { useAuthStore } from '@/stores/modules/auth';

import { Handlebars } from './handlebars';
import { normalizeHandlebarsNbsp } from './normalize';

/**
 * Compile template
 *
 * @param {string} template
 * @param {Object} [context = {}]
 * @param {Object} [instance = Handlebars]
 * @returns {Promise<string>}
 */
export async function compile(template, context = {}, instance = Handlebars) {
  const authStore = useAuthStore();

  const handleBarFunction = instance.compile(normalizeHandlebarsNbsp(template ?? ''));
  const preparedContext = {
    //env: store.getters['template/vars/items'] ?? {}, // TODO: continue
    user: authStore.currentUser ?? {},

    ...context,
  };

  return handleBarFunction(preparedContext);
}
