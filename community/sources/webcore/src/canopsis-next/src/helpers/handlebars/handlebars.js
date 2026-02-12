import promisedHandlebars from 'promised-handlebars';
import HandlebarsLib from 'handlebars';

/**
 * Creates a new Handlebars instance with promise support.
 *
 * TODO: we can remove it in the future because now we are using c-request-helper component
 * instead of sending request directly from request helper
 *
 * @returns {Handlebars} A Handlebars instance wrapped with promised-handlebars for async template compilation
 */
export const createHandlebarsInstance = () => promisedHandlebars(HandlebarsLib);

export const Handlebars = createHandlebarsInstance();
