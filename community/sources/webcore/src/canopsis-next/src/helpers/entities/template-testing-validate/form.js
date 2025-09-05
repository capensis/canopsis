import {
  EVENT_FILTER_TYPES,
  EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES,
  EXTERNAL_DATA_TYPES,
  TEMPLATE_TESTING_TEST_TYPES,
} from '@/constants';

export const convertEventFilter = (form = {}) => {
  const result = [];
  const isChangeEntity = form.type === EVENT_FILTER_TYPES.changeEntity;
  const isEnrichment = form.type === EVENT_FILTER_TYPES.enrichment;

  if (isChangeEntity || isEnrichment) {
    form.external_data.forEach((externalData, index) => {
      if (externalData.type === EXTERNAL_DATA_TYPES.api) {
        result.push({
          key: `external_data.${index}.request.url`,
          textKey: 'URL',
        }, {
          key: `external_data.${index}.request.payload`,
          textKey: 'PAYLOAD',
          textarea: true,
        });

        return;
      }

      externalData.conditions.forEach((condition, conditionIndex) => {
        result.push({
          key: `external_data.${index}.conditions.${conditionIndex}.value`,
          textKey: 'CONDITION VALUE',
        });
      });
    });
  }

  if (isChangeEntity) {
    result.push({
      key: 'config.resource',
      textKey: 'common.resource',
    }, {
      key: 'config.component',
      textKey: 'common.component',
    }, {
      key: 'config.connector',
      textKey: 'common.connector',
    }, {
      key: 'config.connector_name',
      textKey: 'common.connectorName',
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
        textKey: 'ACTION KEY',
      });
    });
  }

  return result;
};

export const convertLinkRule = () => {};

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
