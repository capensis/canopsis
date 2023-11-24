import * as helpers from '@/helpers/handlebars/helpers';

import { Handlebars } from './handlebars';

/**
 * Register handlebars helper
 *
 * @param {string} name
 * @param {Function} helper
 * @returns {*}
 */
export function registerHelper(name, helper) {
  if (Handlebars.helpers[name]) {
    return;
  }

  Handlebars.registerHelper(name, helper);
}

/**
 * Unregister handlebars helper
 *
 * @param {string} name
 * @returns {*}
 */
export function unregisterHelper(name) {
  Handlebars.unregisterHelper(name);
}

registerHelper('duration', helpers.durationHelper);
registerHelper('state', helpers.alarmStateHelper);
registerHelper('tags', helpers.alarmTagsHelper);
registerHelper('request', helpers.requestHelper);
registerHelper('timestamp', helpers.timestampHelper);
registerHelper('internal-link', helpers.internalLinkHelper);
registerHelper('compare', helpers.compareHelper);
registerHelper('concat', helpers.concatHelper);
registerHelper('sum', helpers.sumHelper);
registerHelper('minus', helpers.minusHelper);
registerHelper('mul', helpers.mulHelper);
registerHelper('divide', helpers.divideHelper);
registerHelper('capitalize', helpers.capitalizeHelper);
registerHelper('capitalize-all', helpers.capitalizeAllHelper);
registerHelper('lowercase', helpers.lowercaseHelper);
registerHelper('uppercase', helpers.uppercaseHelper);
registerHelper('replace', helpers.replaceHelper);
registerHelper('copy', helpers.copyHelper);
registerHelper('json', helpers.jsonHelper);
