import {
  EVENT_FILTER_TYPES,
  EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES,
  EXTERNAL_DATA_TYPES,
  TEMPLATE_TESTING_TEST_TYPES,
  ACTION_TYPES,
} from '@/constants';

/**
 * @typedef {Object} TemplateTestingTestValidateFormItem
 * @property {string} key
 * @property {string} textKey
 * @property {Object} [textArgs]
 * @property {boolean} [textarea]
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

    externalDataItem.conditions.forEach((_, conditionIndex) => {
      result.push({
        key: `external_data.${index}.conditions.${conditionIndex}.value`,
        textKey: 'templateTesting.conditionValue',
        templateVarsKey: 'external_data',
      });
    });
  });

  return result;
};

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

// TODO: check all forms bellow

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
        key: `actions.${index}.parameters.author`,
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
        textKey: 'templateTesting.webhookUrl',
        textArgs: { number: index + 1 },
        templateVarsKey: webhookTemplateVarsKey,
      }, {
        key: `actions.${index}.parameters.webhook.request.payload`,
        textKey: 'templateTesting.webhookPayload',
        textArgs: { number: index + 1 },
        textarea: true,
        templateVarsKey: webhookTemplateVarsKey,
      });

      Object.keys(request.headers).forEach((headerKey) => {
        result.push({
          key: `actions.${index}.parameters.webhook.request.headers.${headerKey}`,
          textKey: 'templateTesting.webhookHeader',
          textArgs: { number: index + 1, header: headerKey },
          templateVarsKey: webhookTemplateVarsKey,
        });
      });

      if (declareTicket.ticket_id.template) {
        result.push({
          key: `actions.${index}.parameters.webhook.declare_ticket.ticket_id.value`,
          textKey: 'templateTesting.ticketId',
          textArgs: { number: index + 1 },
          templateVarsKey: 'ticket',
        });
      }

      if (declareTicket.ticket_url.template) {
        result.push({
          key: `actions.${index}.parameters.webhook.declare_ticket.ticket_url.value`,
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

    if ([ACTION_TYPES.pbehavior, ACTION_TYPES.pbehaviorRemove, ACTION_TYPES.assocticket].includes(action.type)) {
      result.push({
        key: `actions.${index}.parameters.output`,
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

  if (form.webhooks) {
    form.webhooks.forEach((webhook, index) => {
      if (webhook.request) {
        if (webhook.request.url) {
          result.push({
            key: `webhooks.${index}.request.url`,
            textKey: 'templateTesting.webhookUrl',
            textArgs: { number: index + 1 },
            templateVarsKey: 'webhook',
          });
        }

        if (webhook.request.payload) {
          result.push({
            key: `webhooks.${index}.request.payload`,
            textKey: 'templateTesting.webhookPayload',
            textArgs: { number: index + 1 },
            textarea: true,
            templateVarsKey: 'webhook',
          });
        }

        if (webhook.request.headers) {
          Object.keys(webhook.request.headers).forEach((headerKey) => {
            result.push({
              key: `webhooks.${index}.request.headers.${headerKey}`,
              textKey: 'templateTesting.webhookHeader',
              textArgs: { number: index + 1, header: headerKey },
              templateVarsKey: 'webhook',
            });
          });
        }
      }

      if (webhook.declare_ticket && webhook.declare_ticket.enabled) {
        if (webhook.declare_ticket.ticket_id) {
          result.push({
            key: `webhooks.${index}.declare_ticket.ticket_id`,
            textKey: 'templateTesting.ticketId',
            textArgs: { number: index + 1 },
            templateVarsKey: 'ticket',
          });
        }

        if (webhook.declare_ticket.ticket_url) {
          result.push({
            key: `webhooks.${index}.declare_ticket.ticket_url`,
            textKey: 'templateTesting.ticketUrl',
            textArgs: { number: index + 1 },
            templateVarsKey: 'ticket',
          });
        }
      }
    });
  }

  result.push(...convertExternalDataToTemplateTestingTestValidateForm(form.external_data));

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

  if (form.infos) {
    form.infos.forEach((info, index) => {
      if (info.value) {
        result.push({
          key: `infos.${index}.value`,
          textKey: 'templateTesting.dynamicInfoValue',
          textArgs: { number: index + 1, name: info.name || `Info ${index + 1}` },
          templateVarsKey: 'dynamicInfo',
        });
      }
    });
  }

  result.push(...convertExternalDataToTemplateTestingTestValidateForm(form.external_data));

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

  if (form.steps) {
    form.steps.forEach((step, stepIndex) => {
      if (step.operations) {
        step.operations.forEach((operation, opIndex) => {
          if (operation.type === 'webhook' && operation.request) {
            if (operation.request.url) {
              result.push({
                key: `steps.${stepIndex}.operations.${opIndex}.request.url`,
                textKey: 'templateTesting.instructionWebhookUrl',
                textArgs: { step: stepIndex + 1, operation: opIndex + 1 },
                templateVarsKey: 'instruction',
              });
            }

            if (operation.request.payload) {
              result.push({
                key: `steps.${stepIndex}.operations.${opIndex}.request.payload`,
                textKey: 'templateTesting.instructionWebhookPayload',
                textArgs: { step: stepIndex + 1, operation: opIndex + 1 },
                textarea: true,
                templateVarsKey: 'instruction',
              });
            }

            if (operation.request.headers) {
              Object.keys(operation.request.headers).forEach((headerKey) => {
                result.push({
                  key: `steps.${stepIndex}.operations.${opIndex}.request.headers.${headerKey}`,
                  textKey: 'templateTesting.instructionWebhookHeader',
                  textArgs: { step: stepIndex + 1, operation: opIndex + 1, header: headerKey },
                  templateVarsKey: 'instruction',
                });
              });
            }
          }
        });
      }
    });
  }

  if (form.jobs) {
    form.jobs.forEach((job, index) => {
      if (job.payload) {
        result.push({
          key: `jobs.${index}.payload`,
          textKey: 'templateTesting.instructionJobPayload',
          textArgs: { number: index + 1 },
          textarea: true,
          templateVarsKey: 'instruction',
        });
      }
    });
  }

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

  if (form.payload) {
    result.push({
      key: 'payload',
      textKey: 'templateTesting.jobPayload',
      textarea: true,
      templateVarsKey: 'job',
    });
  }

  if (form.query) {
    Object.keys(form.query).forEach((queryKey) => {
      result.push({
        key: `query.${queryKey}`,
        textKey: 'templateTesting.jobQueryValue',
        textArgs: { field: queryKey },
        templateVarsKey: 'job',
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
  const result = [];

  if (form.output_template) {
    result.push({
      key: 'output_template',
      textKey: 'templateTesting.metaAlarmOutputTemplate',
      textarea: true,
      templateVarsKey: 'metaAlarm',
    });
  }

  if (form.config) {
    if (form.config.component_template) {
      result.push({
        key: 'config.component_template',
        textKey: 'templateTesting.metaAlarmComponentTemplate',
        templateVarsKey: 'metaAlarm',
      });
    }

    if (form.config.resource_template) {
      result.push({
        key: 'config.resource_template',
        textKey: 'templateTesting.metaAlarmResourceTemplate',
        templateVarsKey: 'metaAlarm',
      });
    }
  }

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
    [TEMPLATE_TESTING_TEST_TYPES.eventFilter]: convertEventFilterToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.linkRule]: convertLinkRuleToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.scenario]: convertScenarioToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.widget]: convertWidgetToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.declareTicketRule]: convertDeclareTicketRuleToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.dynamicInfo]: convertDynamicInfoToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.instruction]: convertInstructionToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.job]: convertJobToTemplateTestingTestValidateForm,
    [TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule]: convertMetaAlarmRuleToTemplateTestingTestValidateForm,
  }[type];

  return converter ? converter(form) : [];
};
