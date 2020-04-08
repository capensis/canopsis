import { setSeveralFields, unsetSeveralFieldsWithConditions } from '@/helpers/immutable';
import { textPairsToObject, objectToTextPairs } from '@/helpers/text-pairs';
import { POST_PROCESSOR_TYPES } from '@/constants';

/**
 * Create default webhook request field
 * @return {Object}
 */
export const getDefaultRequestField = () => ({
  method: '',
  url: '',
  headers: [],
  payload: '{}',
  auth: {
    username: '',
    password: '',
  },
  withAuth: false,
});

/**
 * Create default hook
 * @return {Object}
 */
export const getDefaultHookField = () => ({
  triggers: [],
  event_patterns: [],
  alarm_patterns: [],
  entity_patterns: [],
});

/**
 * Create default post processor
 * @return {Object}
 */
export const getDefaultPostProcessorField = () => ({
  emptyResponse: false,
  ticketId: '',
  fields: [],
});

/**
 * Create webhook form
 * @return {Object}
 */
export const getDefaultWebhookForm = () => ({
  retry: {},
  hook: getDefaultHookField(),
  requests: [getDefaultRequestField()],
  postProcessors: [getDefaultPostProcessorField()],
  disable_if_active_pbehavior: false,
  enabled: true,
});

/**
 * Prepare post processors to form
 * @param {Array} postProcessors
 * @return {Array}
 */
function preparePostProcessorsToForm(postProcessors) {
  return postProcessors.map(({ parameters }) => {
    const {
      empty_response: emptyResponse,
      ticket_id: ticketId,
      fields,
      ...postProcessor
    } = parameters;

    return {
      ...postProcessor,
      fields: objectToTextPairs(fields),
      emptyResponse,
      ticketId,
    };
  });
}

/**
 * Prepare requests to form
 * @param {Array} requests
 * @return {Array}
 */
function prepareRequestsToForm(requests) {
  return requests.map(({ headers, auth, ...otherRequestFields }) => ({
    headers: objectToTextPairs(headers),
    withAuth: !!auth,
    auth,
    ...otherRequestFields,
  }));
}

/**
 * Convert webhook to object for form
 * @param {Array} postProcessors
 * @param {Boolean} enabled
 * @param {Object} webhook
 * @return {Object}
 */
export function webhookToForm({
  post_processors: postProcessors, enabled, ...webhook
}) {
  const patternsFieldsCustomizer = value => value || [];

  return setSeveralFields(webhook, {
    'hook.event_patterns': patternsFieldsCustomizer,
    'hook.alarm_patterns': patternsFieldsCustomizer,
    'hook.entity_patterns': patternsFieldsCustomizer,
    requests: prepareRequestsToForm,
    enabled: enabled === undefined ? true : enabled,
    postProcessors: preparePostProcessorsToForm(postProcessors),
  });
}

/**
 *
 * @param requests
 * @return {Object}
 */
function formRequestFieldToWebhook(requests) {
  return requests.map(({
    withAuth, auth, headers, ...otherFields
  }) => {
    const request = {
      headers: textPairsToObject(headers),
      ...otherFields,
    };

    if (withAuth) {
      request.auth = auth;
    }

    return request;
  });
}

/**
 * Prepare post processors to webhook object
 * @param {Array} postProcessors
 * @return {Array}
 */
function formPostProcessorsToWebhook(postProcessors) {
  return postProcessors.map(({ fields, ticketId, emptyResponse }) => ({
    type: POST_PROCESSOR_TYPES.declareTicket,
    parameters: {
      fields: textPairsToObject(fields),
      ticket_id: ticketId,
      empty_response: emptyResponse,
    },
  }));
}

/**
 * Transform webhook "form" object to valid webhook to the API
 * @param {Object} form
 * @returns {Object}
 */
export function formToWebhook({ requests, postProcessors, ...webhookForm }) {
  const patternsCondition = value => !value || !value.length;
  const hasValue = v => !v;

  const webhook = unsetSeveralFieldsWithConditions(webhookForm, {
    'hook.event_patterns': patternsCondition,
    'hook.alarm_patterns': patternsCondition,
    'hook.entity_patterns': patternsCondition,
    'retry.count': hasValue,
    'retry.unit': hasValue,
    'retry.delay': hasValue,
  });

  return {
    requests: formRequestFieldToWebhook(requests),
    post_processors: formPostProcessorsToWebhook(postProcessors),
    ...webhook,
  };
}
