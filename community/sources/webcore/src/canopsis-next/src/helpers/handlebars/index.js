import { registerAllHelpers } from './registers';
import { createHandlebarsInstance as createPromisedHandlebarsInstance } from './handlebars';

export { createHandlebarsInstance, Handlebars } from './handlebars';
export { compile } from './compilers';
export { getTemplateVariables } from './variables';
export { registerHelper, unregisterHelper, registerAllHelpers } from './registers';
export { registerTemplate, unregisterTemplate, runTemplate, hasTemplate } from './templates';

/**
 * Creates a new Handlebars instance and registers all custom helpers.
 *
 * @returns {Handlebars} A Handlebars instance with all registered custom helpers
 */
export const createHandlebarsInstanceWithHelpers = () => registerAllHelpers(createPromisedHandlebarsInstance());

export const ServiceWeatherHandlebars = createHandlebarsInstanceWithHelpers();
