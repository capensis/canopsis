import promisedHandlebars from 'promised-handlebars';
import HandlebarsLib from 'handlebars';

/**
 * TODO: we can remove it in the future because now we are using c-request-helper component
 * instead of sending request directly from request helper
 */
export const Handlebars = promisedHandlebars(HandlebarsLib);
