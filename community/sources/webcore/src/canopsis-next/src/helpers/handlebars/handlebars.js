import promisedHandlebars from 'promised-handlebars';
import HandlebarsLib from 'handlebars';

/**
 * Creates a new Handlebars instance with promise support.
 *
 * @returns {Handlebars} A Handlebars instance wrapped with promised-handlebars for async template compilation
 */
export const createHandlebarsInstance = () => promisedHandlebars(HandlebarsLib);

export const Handlebars = createHandlebarsInstance();
