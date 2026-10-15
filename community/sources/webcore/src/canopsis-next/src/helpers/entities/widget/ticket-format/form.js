import { map } from 'lodash';

import { uid } from '@/helpers/uid';

/**
 * @typedef {string} TicketFormat
 */

/**
 * @typedef {ObjectKey} TicketFormatForm
 * @property {TicketFormat} value
 */

/**
 * Converts one ticket format to its editable form representation.
 *
 * @param {TicketFormat} [format='']
 * @returns {TicketFormatForm}
 */
export const ticketFormatToForm = (format = '') => ({
  key: uid(),
  value: format,
});

/**
 * Converts ticket formats to editable form items.
 *
 * @param {TicketFormat[]} [formats=[]]
 * @returns {TicketFormatForm[]}
 */
export const ticketFormatsToForm = (formats = []) => formats.map(ticketFormatToForm);

/**
 * Converts editable ticket format items back to API values.
 *
 * @param {TicketFormatForm[]} [form=[]]
 * @returns {TicketFormat[]}
 */
export const formToTicketFormats = (form = []) => map(form, 'value');
