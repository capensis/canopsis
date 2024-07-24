import store from '@/store';

import { Handlebars } from '@/helpers/handlebars/handlebars';

/**
 * Precompile and register template
 *
 * @param {string} id
 * @param {string} template
 * @param {Object} [instance = Handlebars]
 * @returns {Function}
 */
export function registerTemplate(id, template, instance = Handlebars) {
  const precompiledTemplate = instance.precompile(template);

  if (instance.registerPrecompiledTemplate) {
    return instance.registerPrecompiledTemplate(id, precompiledTemplate);
  }

  // eslint-disable-next-line no-new-func
  const compiledTemplateFunc = new Function(`return ${precompiledTemplate}`);

  // eslint-disable-next-line no-param-reassign
  instance.partials[id] = instance.template(compiledTemplateFunc());

  return instance.partials[id];
}

/**
 * Unregisters a template by deleting it from the provided Handlebars instance.
 *
 * @param {string} id - The ID of the template to unregister.
 * @param {Handlebars} [instance = Handlebars] - The Handlebars instance from which to unregister the template.
 */
export function unregisterTemplate(id, instance = Handlebars) {
  // eslint-disable-next-line no-param-reassign
  delete instance.partials?.[id];
}

/**
 * Run handlebars prepared template by id
 *
 * @param {string} id
 * @param {Object} context
 * @param {Object} [instance = Handlebars]
 * @return {Promise<string>}
 */
export async function runTemplate(id, context, instance = Handlebars) {
  const preparedContext = {
    env: store.getters['templateVars/items'] ?? {},
    user: store.getters['auth/currentUser'] ?? {},

    ...context,
  };

  return instance?.partials?.[id]?.(preparedContext) ?? '';
}

/**
 * Check template existence
 *
 * @param {string} id
 * @param {Object} [instance = Handlebars]
 * @return {boolean}
 */
export function hasTemplate(id, instance = Handlebars) {
  return !!instance?.partials?.[id];
}
