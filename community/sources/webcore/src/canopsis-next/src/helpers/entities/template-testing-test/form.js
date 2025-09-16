import { differenceBy } from 'lodash';

import {
  TEMPLATE_TESTING_TEST_TYPES,
  TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES,
  ACTION_TYPES,
  LINK_RULE_TYPES,
  EXTERNAL_DATA_TYPES,
  PATTERNS_FIELDS,
} from '@/constants';

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
 * @typedef {Object} TemplateTestingTest
 * @property {string} name
 * @property {string} description
 * @property {number} rule_type
 * @property {string} rule_name
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
 * Convert template testing test object to form
 *
 * @param {TemplateTestingTest} templateTestingTest
 * @returns {TemplateTestingTest}
 */
export const templateTestingTestToForm = (templateTestingTest = {}) => ({
  name: templateTestingTest.name ?? '',
  description: templateTestingTest.description ?? '',
  rule_type: templateTestingTest.rule_type ?? TEMPLATE_TESTING_TEST_TYPES.scenario,
  rule_name: templateTestingTest.rule_name ?? '',
});

/**
 * Creates a template testing test validate form item
 *
 * @param {TemplateTestingTestValidateFormItemConfig} config - Configuration object
 * @returns {TemplateTestingTestValidateFormItem} The created form item
 */
export const getTemplateTestingTestValidateFormItem = ({
  type,
  params = {},
  required = false,
  someRequired = false,
  index,
  key,
}) => ({
  type,
  params,
  index,
  required,
  someRequired,
  key: key ?? type,
  value: '',
});

/**
 * Converts a form object to template testing test validate form based on the type
 *
 * @param {TemplateTestingTestMainForm} [form] - The form object containing rule configuration
 * @param {TemplateTestingTestType} type - The type of template testing test
 * @returns {TemplateTestingTestValidateFormItem[]} Array of validate form items
 */
export const formToTemplateTestingTestValidateForm = (form = {}, type) => {
  if (type === TEMPLATE_TESTING_TEST_TYPES.eventFilter) {
    return [
      getTemplateTestingTestValidateFormItem({
        type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.event,
        required: true,
        params: { [PATTERNS_FIELDS.event]: [[{ field: 'event_type', cond: { type: 'eq', value: 'check' } }]] },
      }),
      ...form.external_data.reduce((acc, externalData, index) => {
        if (externalData.type === EXTERNAL_DATA_TYPES.api) {
          acc.push(getTemplateTestingTestValidateFormItem({
            index,
            required: true,
            type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.response,
            key: externalData.key,
          }));
        }

        return acc;
      }, []),
    ];
  }

  if (type === TEMPLATE_TESTING_TEST_TYPES.scenario) {
    return [
      getTemplateTestingTestValidateFormItem({
        type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.event,
        required: true,
      }),

      ...form.actions.reduce((acc, action, index) => {
        if (action.type === ACTION_TYPES.webhook) {
          acc.push(getTemplateTestingTestValidateFormItem({
            index,
            required: true,
            key: action.key,
            type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.response,
          }));
        }

        return acc;
      }, []),
    ];
  }

  if (type === TEMPLATE_TESTING_TEST_TYPES.linkRule) {
    return [
      getTemplateTestingTestValidateFormItem({
        type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.user,
      }),

      form.type === LINK_RULE_TYPES.alarm
        ? getTemplateTestingTestValidateFormItem({
          type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.alarm,
          required: true,
          params: { opened: true },
        })
        : getTemplateTestingTestValidateFormItem({
          type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.entity,
          required: true,
        }),
    ];
  }

  if (type === TEMPLATE_TESTING_TEST_TYPES.widget) {
    return [
      getTemplateTestingTestValidateFormItem({
        type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.alarm,
        required: true,
      }),
    ];
  }

  if (type === TEMPLATE_TESTING_TEST_TYPES.declareTicketRule) {
    return [
      getTemplateTestingTestValidateFormItem({
        type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.alarm,
        required: true,
        params: { opened: true },
      }),

      ...form.webhooks.map((webhook, index) => (
        getTemplateTestingTestValidateFormItem({
          index,
          required: true,
          key: webhook.key,
          type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.response,
        })
      )),
    ];
  }

  if (type === TEMPLATE_TESTING_TEST_TYPES.dynamicInfo) {
    return [
      getTemplateTestingTestValidateFormItem({
        type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.alarm,
        someRequired: true,
        params: { opened: true },
      }),
      getTemplateTestingTestValidateFormItem({
        type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.event,
        someRequired: true,
      }),
    ];
  }

  if (type === TEMPLATE_TESTING_TEST_TYPES.instruction) {
    return [
      getTemplateTestingTestValidateFormItem({
        type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.alarm,
        required: true,
        params: { opened: true },
      }),
    ];
  }

  if (type === TEMPLATE_TESTING_TEST_TYPES.job) {
    return [
      getTemplateTestingTestValidateFormItem({
        type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.alarm,
        required: true,
        params: { opened: true },
      }),
    ];
  }

  if (type === TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule) {
    return [
      getTemplateTestingTestValidateFormItem({
        type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.alarm,
        required: true,
        params: { opened: true },
      }),
    ];
  }

  return [];
};

/**
 * Gets the changes between two validate forms by comparing their keys
 *
 * @param {TemplateTestingTestValidateFormItem[]} [form] - The new form array
 * @param {TemplateTestingTestValidateFormItem[]} [oldForm] - The old form array
 * @returns {ValidateFormChanges} Changes between the two forms
 */
export const getChangesForValidateForm = (form = [], oldForm = []) => ({
  added: differenceBy(form, oldForm, 'key'),
  removed: differenceBy(oldForm, form, 'key'),
});

/**
 * Converts a form array to template testing test validate object
 *
 * @param {TemplateTestingTestValidateFormItem[]} [form] - Array of form items
 * @returns {TemplateTestingTestValidate} The validate object with entity values and responses
 */
export const formToTemplateTestingTestValidate = (form = []) => form.reduce((acc, item) => {
  if (item.type === TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.response) {
    if (!acc.responses) {
      acc.responses = {};
    }

    acc.responses[Object.keys(acc.responses).length] = item.value;

    return acc;
  }

  acc[item.type] = item.value;

  return acc;
}, {});

/**
 * Converts template testing test validate object back to form array
 *
 * @param {TemplateTestingTestValidateFormItem[]} [originalForm] - The original form structure
 * @param {TemplateTestingTestValidate} [validate] - The validate object with entity values
 * @returns {TemplateTestingTestValidateFormItem[]} Updated form array with populated values
 */
export const templateTestingTestValidateToForm = (originalForm = [], validate = {}) => {
  const responses = Object.values(validate.responses || {});

  return originalForm.map((item) => {
    if (item.type === TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.response) {
      const response = responses.pop();

      return {
        ...item,
        value: response?.value ?? '',
      };
    }

    const entityValue = validate[item.type] ?? item.value;

    return {
      ...item,
      value: entityValue,
    };
  });
};
