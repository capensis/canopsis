import {
  EVENT_FILTER_TYPES,
  EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES,
  EXTERNAL_DATA_TYPES,
  TEMPLATE_TESTING_TEST_TYPES,
  ACTION_TYPES,
  META_ALARMS_RULE_TYPES,
  DYNAMIC_INFO_INFORMATION_TYPES,
} from '@/constants';

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
 * Converts external data configuration to template testing validation form items
 *
 * @param {Object[]} [externalData=[]] - Array of external data configurations
 * @returns {TemplateTestingTestValidateForm} Array of validation form items for external data
 */
export const convertExternalDataToTemplateTestingTestValidateForm = (externalData = []) => {
  const result = [];

  externalData.forEach((externalDataItem, index) => {
    if (externalDataItem.type === EXTERNAL_DATA_TYPES.api) {
      result.push({
        key: `external_data.${index}.request.url`,
        textKey: 'templateTesting.requestUrl',
        textArgs: { number: index + 1 },
        templateVarsKey: 'external_data',
      }, {
        key: `external_data.${index}.request.payload`,
        textKey: 'templateTesting.requestPayload',
        textArgs: { number: index + 1 },
        textarea: true,
        templateVarsKey: 'external_data',
      });

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
 * Converts headers array to template testing test validate form structure
 *
 * @param {Object} [options={}] - Options object with headers, prefix, resultPrefix and templateVarsKey
 * @param {Array} [headers=[]] - Array of header objects with text property
 * @param {string} [prefix=''] - Prefix string for form field keys
 * @param {string} [resultPrefix=''] - Prefix string for result keys
 * @param {string} [templateVarsKey=''] - Template variables key for validation
 * @returns {Array} Array of form validation objects with keys, result keys, text keys and template vars
 */
export const convertHeadersToTemplateTestingTestValidateForm = ({
  headers = {},
  prefix = '',
  resultPrefix = prefix,
  templateVarsKey = '',
}) => (
  headers.map((header, index) => ({
    key: `${prefix}.headers.${index}.value`,
    resultKey: `${resultPrefix}.headers.${header.text}`,
    textKey: 'templateTesting.webhookHeader',
    textArgs: { number: index + 1, header: header.text },
    templateVarsKey,
  }))
);

/**
 * Converts event filter form to template testing validation form items
 *
 * @param {EventFilterForm} [form={}] - Event filter configuration form
 * @returns {TemplateTestingTestValidateForm} Array of validation form items for event filter
 */
export const convertEventFilterToTemplateTestingTestValidateForm = (form = {}) => {
  const result = [];
  const isChangeEntity = form.type === EVENT_FILTER_TYPES.changeEntity;
  const isEnrichment = form.type === EVENT_FILTER_TYPES.enrichment;

  if (isChangeEntity || isEnrichment) {
    result.push(...convertExternalDataToTemplateTestingTestValidateForm(form.external_data));
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
        textKey: 'templateTesting.actionValue',
        textArgs: { number: index + 1 },
        templateVarsKey: 'config',
      });
    });
  }

  return result;
};

/**
 * Converts link rule form to template testing validation form items
 *
 * @param {LinkRuleForm} [form={}] - Link rule configuration form
 * @returns {TemplateTestingTestValidateForm} Array of validation form items for link rule
 */
export const convertLinkRuleToTemplateTestingTestValidateForm = (form = {}) => {
  const result = [];

  form.links.forEach((_, index) => {
    result.push({
      key: `links.${index}.label`,
      textKey: 'templateTesting.linkRuleLabel',
      textArgs: { number: index + 1 },
      templateVarsKey: 'label',
    }, {
      key: `links.${index}.url`,
      textKey: 'templateTesting.linkRuleUrl',
      textArgs: { number: index + 1 },
      templateVarsKey: 'url',
    });
  });

  result.push(...convertExternalDataToTemplateTestingTestValidateForm(form.external_data));

  return result;
};

/**
 * Converts scenario form to template testing validation form items
 *
 * @param {ScenarioForm} [form={}] - Scenario configuration form
 * @returns {TemplateTestingTestValidateForm} Array of validation form items for scenario
 */
export const convertScenarioToTemplateTestingTestValidateForm = (form = {}) => {
  const result = [];

  let firstWebhook = true;

  form.actions.forEach((action, index) => {
    if (!action.parameters[action.type]?.forward_author) {
      result.push({
        key: `actions.${index}.parameters.${action.type}.author`,
        resultKey: `actions.${index}.parameters.author`,
        textKey: 'templateTesting.actionAuthor',
        textArgs: { number: index + 1 },
        templateVarsKey: 'author',
      });
    }

    if (action.type === ACTION_TYPES.webhook) {
      const webhookTemplateVarsKey = firstWebhook ? 'first_webhook' : 'webhook';
      const { request, declare_ticket: declareTicket } = action.parameters.webhook;

      result.push({
        key: `actions.${index}.parameters.webhook.request.url`,
        resultKey: `actions.${index}.parameters.request.url`,
        textKey: 'templateTesting.webhookUrl',
        textArgs: { number: index + 1 },
        templateVarsKey: webhookTemplateVarsKey,
      }, {
        key: `actions.${index}.parameters.webhook.request.payload`,
        resultKey: `actions.${index}.parameters.request.payload`,
        textKey: 'templateTesting.webhookPayload',
        textArgs: { number: index + 1 },
        textarea: true,
        templateVarsKey: webhookTemplateVarsKey,
      });

      result.push(
        ...convertHeadersToTemplateTestingTestValidateForm(
          {
            headers: request.headers,
            prefix: `actions.${index}.parameters.webhook.request`,
            resultPrefix: `actions.${index}.parameters.request`,
            templateVarsKey: webhookTemplateVarsKey,
          },
        ),
      );

      if (declareTicket.ticket_id.template) {
        result.push({
          key: `actions.${index}.parameters.webhook.declare_ticket.ticket_id.value`,
          resultKey: `actions.${index}.parameters.declare_ticket.ticket_id_tpl`,
          textKey: 'templateTesting.ticketId',
          textArgs: { number: index + 1 },
          templateVarsKey: 'ticket',
        });
      }

      if (declareTicket.ticket_url.template) {
        result.push({
          key: `actions.${index}.parameters.webhook.declare_ticket.ticket_url.value`,
          resultKey: `actions.${index}.parameters.declare_ticket.ticket_url_tpl`,
          textKey: 'templateTesting.ticketUrl',
          textArgs: { number: index + 1 },
          templateVarsKey: 'ticket',
        });
      }

      if (firstWebhook) {
        firstWebhook = false;
      }

      return;
    }

    if (![ACTION_TYPES.pbehavior, ACTION_TYPES.pbehaviorRemove, ACTION_TYPES.assocticket].includes(action.type)) {
      result.push({
        key: `actions.${index}.parameters.${action.type}.output`,
        resultKey: `actions.${index}.parameters.output`,
        textKey: 'templateTesting.noteOutput',
        textArgs: { number: index + 1 },
        textarea: true,
        templateVarsKey: 'output',
      });
    }
  });

  result.push(...convertExternalDataToTemplateTestingTestValidateForm(form.external_data));

  return result;
};

/**
 * Converts widget form to template testing validation form items
 *
 * @param {WidgetForm} [form={}] - Widget configuration form
 * @returns {TemplateTestingTestValidateForm} Array of validation form items for widget
 */
export const convertWidgetToTemplateTestingTestValidateForm = (form = {}) => {
  const result = [];

  if (form.parameters) {
    if (form.parameters.content) {
      result.push({
        key: 'parameters.content',
        textKey: 'templateTesting.widgetContent',
        textarea: true,
        templateVarsKey: 'widget',
      });
    }

    if (form.parameters.templates) {
      Object.keys(form.parameters.templates).forEach((templateKey) => {
        const template = form.parameters.templates[templateKey];
        if (template && typeof template === 'string') {
          result.push({
            key: `parameters.templates.${templateKey}`,
            textKey: 'templateTesting.widgetTemplate',
            textArgs: { template: templateKey },
            textarea: true,
            templateVarsKey: 'widget',
          });
        }
      });
    }

    if (form.parameters.widgetColumns) {
      form.parameters.widgetColumns.forEach((column, index) => {
        if (column.template) {
          result.push({
            key: `parameters.widgetColumns.${index}.template`,
            textKey: 'templateTesting.columnTemplate',
            textArgs: { number: index + 1 },
            textarea: true,
            templateVarsKey: 'widget',
          });
        }
      });
    }

    if (form.parameters.columns) {
      form.parameters.columns.forEach((column, index) => {
        if (column.template) {
          result.push({
            key: `parameters.columns.${index}.template`,
            textKey: 'templateTesting.columnTemplate',
            textArgs: { number: index + 1 },
            textarea: true,
            templateVarsKey: 'widget',
          });
        }
      });
    }
  }

  return result;
};

/**
 * Converts declare ticket rule form to template testing validation form items
 *
 * @param {DeclareTicketRuleForm} [form={}] - Declare ticket rule configuration form
 * @returns {TemplateTestingTestValidateForm} Array of validation form items for declare ticket rule
 */
export const convertDeclareTicketRuleToTemplateTestingTestValidateForm = (form = {}) => {
  const result = [];

  form.webhooks?.forEach?.((webhook, index) => {
    const templateVarsKey = !index ? 'first_webhook' : 'webhook';

    result.push({
      templateVarsKey,

      key: `webhooks.${index}.request.url`,
      textKey: 'templateTesting.webhookUrl',
      textArgs: { number: index + 1 },
    }, {
      templateVarsKey,

      key: `webhooks.${index}.request.payload`,
      textKey: 'templateTesting.webhookPayload',
      textArgs: { number: index + 1 },
      textarea: true,
    });

    result.push(
      ...convertHeadersToTemplateTestingTestValidateForm(
        {
          headers: webhook.request.headers,
          prefix: `webhooks.${index}.request`,
          templateVarsKey,
        },
      ),
    );

    if (!webhook.declare_ticket.enabled) {
      return;
    }

    if (webhook.declare_ticket.ticket_id.template) {
      result.push({
        key: `webhooks.${index}.declare_ticket.ticket_id.value`,
        textKey: 'templateTesting.ticketId',
        textArgs: { number: index + 1 },
        templateVarsKey: 'ticket',
      });
    }

    if (webhook.declare_ticket.ticket_url.template) {
      result.push({
        key: `webhooks.${index}.declare_ticket.ticket_url.value`,
        textKey: 'templateTesting.ticketUrl',
        textArgs: { number: index + 1 },
        templateVarsKey: 'ticket',
      });
    }
  });

  return result;
};

/**
 * Converts dynamic info form to template testing validation form items
 *
 * @param {DynamicInfoForm} [form={}] - Dynamic info configuration form
 * @returns {TemplateTestingTestValidateForm} Array of validation form items for dynamic info
 */
export const convertDynamicInfoToTemplateTestingTestValidateForm = (form = {}) => {
  const result = [];

  form.infos.forEach((info, index) => {
    if (info.type !== DYNAMIC_INFO_INFORMATION_TYPES.setToInfoFromTemplate) {
      return;
    }

    result.push({
      key: `infos.${index}.value`,
      textKey: 'templateTesting.dynamicInfoValue',
      textArgs: { number: index + 1, name: info.name },
      templateVarsKey: 'value',
    });
  });

  return result;
};

/**
 * Converts remediation instruction form to template testing validation form items
 *
 * @param {RemediationInstructionForm} [form={}] - Remediation instruction configuration form
 * @returns {TemplateTestingTestValidateForm} Array of validation form items for remediation instruction
 */
export const convertInstructionToTemplateTestingTestValidateForm = (form = {}) => {
  const result = [];

  form.steps.forEach((step, stepIndex) => {
    step.operations.forEach((_, operationIndex) => {
      result.push({
        key: `steps.${stepIndex}.operations.${operationIndex}.description`,
        textKey: 'templateTesting.operationDescription',
        textArgs: { step: stepIndex + 1, operation: operationIndex + 1 },
        templateVarsKey: 'operation',
      });
    });
  });

  return result;
};

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
 * Converts meta alarm rule form to template testing validation form items
 *
 * @param {MetaAlarmRuleForm} [form={}] - Meta alarm rule configuration form
 * @returns {TemplateTestingTestValidateForm} Array of validation form items for meta alarm rule
 */
export const convertMetaAlarmRuleToTemplateTestingTestValidateForm = (form = {}) => {
  const result = [{
    key: 'output_template',
    textKey: 'templateTesting.metaAlarmOutputTemplate',
    textarea: true,
    templateVarsKey: 'output',
  }];

  if ([
    META_ALARMS_RULE_TYPES.timebased,
    META_ALARMS_RULE_TYPES.attribute,
    META_ALARMS_RULE_TYPES.complex,
    META_ALARMS_RULE_TYPES.valuegroup,
  ].includes(form.type)) {
    result.push({
      key: 'config.component_template',
      textKey: 'templateTesting.metaAlarmComponentTemplate',
      templateVarsKey: 'entity',
    }, {
      key: 'config.resource_template',
      textKey: 'templateTesting.metaAlarmResourceTemplate',
      templateVarsKey: 'entity',
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
export const convertExternalAuthTokenRuleToTemplateTestingTestValidateForm = (form = {}) => (
  form.allow_variables
    ? [{
      key: 'template',
      textKey: 'common.token',
      templateVarsKey: 'template',
    }]
    : []
);

/**
 * Converts rule form to template testing validation form items based on rule type
 *
 * @param {TemplateTestingTestMainForm} [form={}] - Rule configuration form
 * @param {TemplateTestingTestType} type - The type of rule to convert
 * @returns {TemplateTestingTestValidateForm} Array of validation form items for the specified rule type
 */
export const convertRuleToTemplateTestingTestValidateForm = (form = {}, type) => {
  const converter = {
    [TEMPLATE_TESTING_TEST_TYPES.eventFilter]: convertEventFilterToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.linkRule]: convertLinkRuleToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.scenario]: convertScenarioToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.widget]: convertWidgetToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.declareTicketRule]: convertDeclareTicketRuleToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.dynamicInfo]: convertDynamicInfoToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.instruction]: convertInstructionToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.job]: convertJobToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule]: convertMetaAlarmRuleToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.externalAuthToken]: convertExternalAuthTokenRuleToTemplateTestingTestValidateForm,
  }[type];

  return converter ? converter(form) : [];
};
