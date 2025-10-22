import { TEMPLATE_TESTING_DATA_TYPES } from '@/constants';

import { objectToTextPairs, textPairsToObject } from '@/helpers/text-pairs';

/**
 * @typedef {Object} TemplateTestingData
 * @property {string} name
 * @property {string} description
 * @property {number} type
 * @property {Object} body
 * @property {Object} [headers]
 */

/**
 * @typedef {TemplateTestingData} TemplateTestingDataForm
 * @property {string} body
 * @property {TextPairObject[]} headers
 */

/**
 * Convert template testing data object to form
 *
 * @param {TemplateTestingData} templateTestingData
 * @returns {TemplateTestingDataForm}
 */
export const templateTestingDataToForm = (templateTestingData = {}) => ({
  name: templateTestingData.name ?? '',
  description: templateTestingData.description ?? '',
  type: templateTestingData.type ?? TEMPLATE_TESTING_DATA_TYPES.event,
  body: JSON.stringify(templateTestingData.body ?? {}, null, 2),
  headers: templateTestingData.headers ? objectToTextPairs(templateTestingData.headers) : [],
});

/**
 * Convert form object to template testing data
 *
 * @param {TemplateTestingDataForm} form
 * @returns {TemplateTestingData}
 */
export const formToTemplateTestingData = form => ({
  name: form.name,
  description: form.description,
  type: form.type,
  body: JSON.parse(form.body),
  headers: textPairsToObject(form.headers),
});
