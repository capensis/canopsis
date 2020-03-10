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
export function compile(template, context) {
  const handleBarFunction = Handlebars.compile(template);

  return handleBarFunction(context);
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
registerHelper('date', helpers.dateHelper);
registerHelper('state', helpers.alarmStateHelper);
registerHelper('request', helpers.requestHelper);
registerHelper('timestamp', helpers.timestamp);
registerHelper('internal-link', helpers.internalLink);
registerHelper('compare', helpers.compare);
