import * as helpers from '@/helpers/handlebars/helpers';

import { Handlebars } from './handlebars';

const helperRegistrationsByInstance = new WeakMap();

/**
 * Returns the helper registration map for a Handlebars instance.
 * Creates and stores an empty map when the instance has none yet.
 *
 * @param {Handlebars} instance - The Handlebars instance to look up
 * @returns {Map<string, {count: number, owned: boolean}>} Registrations keyed by helper name
 */
const getHelperRegistrations = (instance) => {
  let registrations = helperRegistrationsByInstance.get(instance);

  if (!registrations) {
    registrations = new Map();
    helperRegistrationsByInstance.set(instance, registrations);
  }

  return registrations;
};

/**
 * Register handlebars helper
 *
 * @param {string} name
 * @param {Function} helper
 * @returns {*}
 */
export function registerHelper(name, helper, instance = Handlebars) {
  const registrations = getHelperRegistrations(instance);
  const registration = registrations.get(name);

  if (registration) {
    registration.count += 1;

    return instance.helpers[name];
  }

  const owned = !instance.helpers[name];

  if (owned) {
    instance.registerHelper(name, helper);
  }

  registrations.set(name, { count: 1, owned });

  return instance.helpers[name];
}

/**
 * Unregister handlebars helper
 *
 * @param {string} name
 * @returns {*}
 */
export function unregisterHelper(name, instance = Handlebars) {
  const registrations = helperRegistrationsByInstance.get(instance);
  const registration = registrations?.get(name);

  if (!registration) {
    instance.unregisterHelper(name);

    return;
  }

  registration.count -= 1;

  if (registration.count > 0) {
    return;
  }

  registrations.delete(name);

  if (!registrations.size) {
    helperRegistrationsByInstance.delete(instance);
  }

  if (registration.owned) {
    instance.unregisterHelper(name);
  }
}

/**
 * Registers all custom Handlebars helpers to a Handlebars instance.
 *
 * Registers the following helpers:
 * - duration: Format duration values
 * - state: Display alarm state
 * - request: Make HTTP requests
 * - timestamp: Format timestamps
 * - internal-link: Create internal application links
 * - compare: Compare values
 * - concat: Concatenate strings
 * - sum: Sum numbers
 * - minus: Subtract numbers
 * - mul: Multiply numbers
 * - divide: Divide numbers
 * - capitalize: Capitalize first letter
 * - capitalize-all: Capitalize all words
 * - lowercase: Convert to lowercase
 * - uppercase: Convert to uppercase
 * - replace: Replace string patterns
 * - copy: Copy values
 * - json: Stringify JSON
 * - map: Map values
 *
 * @param {Handlebars} [instance=Handlebars] - The Handlebars instance to register helpers to
 * @returns {Handlebars} The Handlebars instance with all helpers registered
 */
export const registerAllHelpers = (instance = Handlebars) => {
  registerHelper('duration', helpers.durationHelper, instance);
  registerHelper('state', helpers.alarmStateHelper, instance);
  registerHelper('request', helpers.requestHelper, instance);
  registerHelper('timestamp', helpers.timestampHelper, instance);
  registerHelper('now', helpers.nowHelper, instance);
  registerHelper('internal-link', helpers.internalLinkHelper, instance);
  registerHelper('compare', helpers.compareHelper, instance);
  registerHelper('concat', helpers.concatHelper, instance);
  registerHelper('sum', helpers.sumHelper, instance);
  registerHelper('minus', helpers.minusHelper, instance);
  registerHelper('mul', helpers.mulHelper, instance);
  registerHelper('divide', helpers.divideHelper, instance);
  registerHelper('capitalize', helpers.capitalizeHelper, instance);
  registerHelper('capitalize-all', helpers.capitalizeAllHelper, instance);
  registerHelper('lowercase', helpers.lowercaseHelper, instance);
  registerHelper('uppercase', helpers.uppercaseHelper, instance);
  registerHelper('replace', helpers.replaceHelper, instance);
  registerHelper('copy', helpers.copyHelper, instance);
  registerHelper('json', helpers.jsonHelper, instance);
  registerHelper('map', helpers.mapHelper, instance);

  return instance;
};

registerAllHelpers();
