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
 * @typedef {Object} TemplateTestingTest
 * @property {string} name
 * @property {string} description
 * @property {number} rule_type
 * @property {string} rule_name
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

export const getChangesForValidateForm = (form = [], oldForm = []) => ({
  added: differenceBy(form, oldForm, 'key'),
  removed: differenceBy(oldForm, form, 'key'),
});

export const formToTemplateTestingTestValidate = (form = [], type) => {
  const responseKey = type === TEMPLATE_TESTING_TEST_TYPES.eventFilter ? 'external_data' : 'responses';

  return form.reduce((acc, item) => {
    if (item.type === TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.response) {
      if (!acc[responseKey]) {
        acc[responseKey] = {};
      }

      acc[responseKey][Object.keys(acc[responseKey]).length] = item.value;

      return acc;
    }

    acc[item.type] = item.value;

    return acc;
  }, {});
};
