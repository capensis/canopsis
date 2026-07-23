import { isArray, mapValues } from 'lodash';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

/**
 * @typedef {0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 } TemplateTestingTestType
 * Template testing test type values:
 * - 0: eventFilter
 * - 1: linkRule
 * - 2: scenario
 * - 3: widget
 * - 4: declareTicketRule
 * - 5: dynamicInfo
 * - 6: instruction
 * - 7: job
 * - 8: metaAlarmRule
 */

/**
 * @typedef {Object|Array} TemplateTestingDataInput
 * @property {string} [_id] - Unique identifier of the data item
 * @property {string} [name] - Name of the data item
 * @property {*} [*] - Any other properties that may contain nested objects or arrays
 */

/**
 * @typedef {Object|Array|string} TemplateTestingDataForm
 * Object with nested objects converted to their IDs, arrays processed recursively, or string IDs
 */

/**
 * @typedef {Object} TemplateTestingTestForm
 * @property {string} name - Name of the template testing test
 * @property {string} description - Description of the template testing test
 * @property {TemplateTestingTestType} type - Type of the template testing test
 * @property {Object} rule - Rule object containing rule configuration
 * @property {string} [rule._id] - ID of the rule
 * @property {string} [rule.name] - Name of the rule
 * @property {TemplateTestingDataForm} data - Additional data for the template testing test
 * @property {TemplateTestingTestValidate} [validate] - Validation data for the test
 */

/**
 * @typedef {Object} TemplateTestingTestValidateFormItem
 * @property {string} type - The type of the validate form item
 * @property {Object} params - Additional parameters for the item
 * @property {number} [index] - Index of the item in the array
 * @property {boolean} required - Whether the item is required
 * @property {boolean} someRequired - Whether the item has some requirement
 * @property {string} key - Unique key for the item
 * @property {string} value - Value of the item
 */

/**
 * @typedef {Object} ValidateFormChanges
 * @property {TemplateTestingTestValidateFormItem[]} added - Items added to the form
 * @property {TemplateTestingTestValidateFormItem[]} removed - Items removed from the form
 */

/**
 * @typedef {Object} TemplateTestingTestValidate
 * @property {string} [event] - Event data
 * @property {string} [alarm] - Alarm data
 * @property {string} [entity] - Entity data
 * @property {string} [user] - User data
 * @property {Object<number, string>} [responses] - Response data mapping
 */

/**
 * @typedef {
 *   EventFilterForm |
 *   LinkRuleForm |
 *   ScenarioForm |
 *   WidgetForm |
 *   DeclareTicketRuleForm |
 *   DynamicInfoForm |
 *   RemediationInstructionForm |
 *   RemediationJobForm |
 *   MetaAlarmRuleForm
 * } TemplateTestingTestMainForm
 */

/**
 * @typedef {Object} TemplateTestingTestValidateFormItemConfig
 * @property {string} type - The type of the validate form item
 * @property {Object} [params] - Additional parameters for the item
 * @property {boolean} [required] - Whether the item is required
 * @property {boolean} [someRequired] - Whether the item has some requirement
 * @property {number} [index] - Index of the item in the array
 * @property {string} [key] - Unique key for the item, defaults to type if not provided
 */

/**
 * Converts template testing data by extracting IDs from nested objects
 *
 * @param {Object|Array} [templateTestingData={}]
 * @returns {Object|Array|string}
 */
export const templateTestingDataToForm = (templateTestingData = {}) =>
  mapValues(templateTestingData, (dataItem) => {
    if (isArray(dataItem)) {
      return dataItem.map((item) => item?._id || item);
    }

    return dataItem?._id;
  });

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
