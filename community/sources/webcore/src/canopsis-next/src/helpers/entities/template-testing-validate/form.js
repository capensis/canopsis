import {
  EVENT_FILTER_TYPES,
  EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES,
  EXTERNAL_DATA_TYPES,
  TEMPLATE_TESTING_TEST_TYPES,
} from '@/constants';

/**
 * @typedef {Object} TemplateTestingTestValidateFormItem
 * @property {string} key
 * @property {string} textKey
 * @property {boolean} [textarea]
 */

/**
 * @typedef {TemplateTestingTestValidateFormItem[]} TemplateTestingTestValidateForm
 */

export const convertExternalData = (externalData = []) => {
  const result = [];

  externalData.forEach((externalDataItem, index) => {
    if (externalDataItem.type === EXTERNAL_DATA_TYPES.api) {
      result.push({
        key: `external_data.${index}.request.url`,
        textKey: 'URL', // TODO: i18n
        templateVarsKey: 'external_data',
      }, {
        key: `external_data.${index}.request.payload`,
        textKey: 'PAYLOAD', // TODO: i18n
        textarea: true,
        templateVarsKey: 'external_data',
      });

      return;
    }

    externalDataItem.conditions.forEach((_, conditionIndex) => {
      result.push({
        key: `external_data.${index}.conditions.${conditionIndex}.value`,
        textKey: 'CONDITION VALUE', // TODO: i18n
        templateVarsKey: 'external_data',
      });
    });
  });

  return result;
};

export const convertEventFilter = (form = {}) => {
  const result = [];
  const isChangeEntity = form.type === EVENT_FILTER_TYPES.changeEntity;
  const isEnrichment = form.type === EVENT_FILTER_TYPES.enrichment;

  if (isChangeEntity || isEnrichment) {
    result.push(...convertExternalData(form.external_data));
  }

  if (isChangeEntity) {
    result.push({
      key: 'config.resource',
      textKey: 'common.resource',
      templateVarsKey: 'config',
    }, {
      key: 'config.component',
      textKey: 'common.component',
      templateVarsKey: 'config',
    }, {
      key: 'config.connector',
      textKey: 'common.connector',
      templateVarsKey: 'config',
    }, {
      key: 'config.connector_name',
      textKey: 'common.connectorName',
      templateVarsKey: 'config',
    });
  }

  if (isEnrichment) {
    form.config.actions.forEach((action, index) => {
      if (
        ![
          EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate,
          EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate,
          EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTagsFromTemplate,
        ].includes(action.type)
      ) {
        return;
      }

      result.push({
        key: `config.actions.${index}.value`,
        textKey: 'ACTION KEY', // TODO: use i18n
        templateVarsKey: 'config',
      });
    });
  }

  return result;
};

export const convertLinkRule = (form = {}, templateVars = {}) => {
  const result = [];

  form.links.forEach((link, index) => {
    result.push({
      key: `links.${index}.label`,
      textKey: 'LABEL', // TODO: i18n
      variables: templateVars.label,
    }, {
      key: `links.${index}.url`,
      textKey: 'URL', // TODO: i18n
      variables: templateVars.url,
    });
  });

  result.push(...convertExternalData(form.external_data));

  return result;
};

export const convertScenario = () => {};

export const convertWidget = () => {};

export const convertDeclareTicketRule = () => {};

export const convertDynamicInfo = () => {};

export const convertInstruction = () => {};

export const convertJob = () => {};

export const convertMetaAlarmRule = () => {};

export const convertRuleToTemplateTestingTestValidateForm = (form = {}, type) => {
  const converter = {
    [TEMPLATE_TESTING_TEST_TYPES.eventFilter]: convertEventFilter,
    [TEMPLATE_TESTING_TEST_TYPES.linkRule]: convertLinkRule,
    [TEMPLATE_TESTING_TEST_TYPES.scenario]: convertScenario,
    [TEMPLATE_TESTING_TEST_TYPES.widget]: convertWidget,
    [TEMPLATE_TESTING_TEST_TYPES.declareTicketRule]: convertDeclareTicketRule,
    [TEMPLATE_TESTING_TEST_TYPES.dynamicInfo]: convertDynamicInfo,
    [TEMPLATE_TESTING_TEST_TYPES.instruction]: convertInstruction,
    [TEMPLATE_TESTING_TEST_TYPES.job]: convertJob,
    [TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule]: convertMetaAlarmRule,
  }[type];

  return converter ? converter(form) : [];
};
