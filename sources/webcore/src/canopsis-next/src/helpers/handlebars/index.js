import promisedHandlebars from 'promised-handlebars';
import HandlebarsLib from 'handlebars';

import * as helpers from './helpers';

const Handlebars = promisedHandlebars(HandlebarsLib);

/**
 * Compile template
 *
 * @param {string} template
 * @param {Object} context
 * @returns {Promise}
 */
export async function compile(template, context) {
  const handleBarFunction = Handlebars.compile(template);
  const result = await handleBarFunction(context);

  const element = document.createElement('div');

  element.innerHTML = result;

  return element.innerHTML;
}

/**
 * Register handlebars helper
 *
 * @param {string} name
 * @param {Function} helper
 * @returns {*}
 */
export function registerHelper(name, helper) {
  return Handlebars.registerHelper(name, helper);
}

/**
 * Unregister handlebars helper
 *
 * @param {string} name
 * @returns {*}
 */
export function unregisterHelper(name) {
  return Handlebars.unregisterHelper(name);
}

/**
 * Register global helpers
 */
registerHelper('duration', helpers.durationHelper);
registerHelper('state', helpers.alarmStateHelper);
registerHelper('request', helpers.requestHelper);
registerHelper('timestamp', helpers.timestampHelper);
registerHelper('internal-link', helpers.internalLinkHelper);
registerHelper('compare', helpers.compareHelper);
registerHelper('concat', helpers.concatHelper);
registerHelper('sum', helpers.sumHelper);
registerHelper('minus', helpers.minusHelper);
registerHelper('mul', helpers.mulHelper);
registerHelper('divide', helpers.divideHelper);
