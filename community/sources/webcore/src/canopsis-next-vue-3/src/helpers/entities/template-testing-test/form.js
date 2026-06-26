import { isArray, mapValues } from 'lodash';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

/**
 * Converts template testing data by extracting IDs from nested objects
 *
 * @param {Object|Array} [templateTestingData={}]
 * @returns {Object|Array|string}
 */
export const templateTestingDataToForm = (templateTestingData = {}) => (
  mapValues(templateTestingData, (dataItem) => {
    if (isArray(dataItem)) {
      return dataItem.map(item => item?._id || item);
    }

    return dataItem?._id;
  })
);

/**
 * Convert template testing test object to form
 *
 * @param {Object} [templateTestingTest={}]
 * @returns {Object}
 */
export const templateTestingTestToForm = (templateTestingTest = {}) => ({
  name: templateTestingTest.name ?? '',
  description: templateTestingTest.description ?? '',
  type: templateTestingTest.type ?? TEMPLATE_TESTING_TEST_TYPES.scenario,
  rule: templateTestingTest.rule ?? { name: '' },
  data: templateTestingDataToForm(templateTestingTest.data),
});

/**
 * Converts a template testing test form object to a template testing test object
 *
 * @param {Object} [form={}]
 * @returns {Object}
 */
export const formToTemplateTestingTest = (form = {}) => ({
  ...form,

  rule: form.rule?._id,
});
