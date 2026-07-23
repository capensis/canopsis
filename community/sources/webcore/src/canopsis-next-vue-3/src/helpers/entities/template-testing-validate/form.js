import { EXTERNAL_DATA_TYPES, TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

/**
 * @typedef {Object} TemplateTestingTestValidateFormItem
 * @property {string} key
 * @property {string} [resultKey]
 * @property {string} textKey
 * @property {Object} [textArgs]
 * @property {boolean} [textarea]
 * @property {boolean} [json]
 * @property {string} [templateVarsKey]
 */

/**
 * @typedef {TemplateTestingTestValidateFormItem[]} TemplateTestingTestValidateForm
 */

/**
 * Converts remediation job form to template testing validation form items
 *
 * @param {RemediationJobForm} [form={}] - Remediation job configuration form
 * @returns {TemplateTestingTestValidateForm} Array of validation form items for remediation job
 */
export const convertJobToTemplateTestingTestValidateForm = (form = {}) => {
  const result = [];

  if (form.payload && form.configType?.with_body) {
    result.push({
      key: 'payload',
      textKey: 'templateTesting.jobPayload',
      textarea: true,
      templateVarsKey: 'payload',
    });
  }

  if (form.configType?.with_query) {
    form.query.forEach((query, index) => {
      result.push({
        key: `query.${index}.value`,
        resultKey: `query.${query.text}`,
        textKey: 'templateTesting.jobQueryValue',
        textArgs: { number: index + 1 },
        templateVarsKey: 'payload',
      });
    });
  }

  return result;
};

/**
 * Converts external auth token rule form to template testing validation form items
 *
 * @param {ExternalAuthTokenRuleForm} [form={}] - External auth token rule configuration form
 * @returns {TemplateTestingTestValidateForm} Array of validation form items for external auth token rule
 */
export const convertExternalAuthTokenRuleToTemplateTestingTestValidateForm = (form = {}) =>
  form.allow_variables
    ? [
        {
          key: 'template',
          textKey: 'common.token',
          templateVarsKey: 'template',
        },
      ]
    : [];

/**
 * Converts external data configuration to template testing validation form items
 *
 * @param {Object[]} [externalData=[]] - Array of external data configurations
 * @returns {TemplateTestingTestValidateForm} Array of validation form items for external data
 */
export const convertExternalDataToTemplateTestingTestValidateForm = (externalData = []) => {
  const result = [];

  externalData.forEach((externalDataItem, index) => {
    if (externalDataItem.type === EXTERNAL_DATA_TYPES.api) {
      result.push(
        {
          key: `external_data.${index}.request.url`,
          textKey: 'templateTesting.requestUrl',
          textArgs: { number: index + 1 },
          templateVarsKey: 'external_data',
        },
        {
          key: `external_data.${index}.request.payload`,
          textKey: 'templateTesting.requestPayload',
          textArgs: { number: index + 1 },
          textarea: true,
          templateVarsKey: 'external_data',
        },
      );

      return;
    }

    externalDataItem.conditions.forEach((condition, conditionIndex) => {
      result.push({
        key: `external_data.${index}.conditions.${conditionIndex}.value`,
        resultKey: `external_data.${index}.${condition.type}.${condition.attribute}`,
        textKey: 'templateTesting.conditionValue',
        templateVarsKey: 'external_data',
      });
    });
  });

  return result;
};

/**
 * Converts rule form to template testing validation form items based on rule type
 *
 * @param {TemplateTestingTestMainForm} [form={}] - Rule configuration form
 * @param {TemplateTestingTestType} type - The type of rule to convert
 * @returns {TemplateTestingTestValidateForm} Array of validation form items for the specified rule type
 */
export const convertRuleToTemplateTestingTestValidateForm = (form = {}, type) => {
  const converter = {
    [TEMPLATE_TESTING_TEST_TYPES.job]: convertJobToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.externalAuthToken]:
      convertExternalAuthTokenRuleToTemplateTestingTestValidateForm,
  }[type];

  return converter ? converter(form) : [];
};
