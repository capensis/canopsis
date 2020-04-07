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
 *
 * @param {Array} postProcessors
 * @param {Object} webhook
 * @return {Object}
 */
export function webhookToForm({
  post_processors: postProcessors, enabled, request, ...webhook
}) {
  const patternsFieldsCustomizer = value => value || [];

  return setSeveralFields(webhook, {
    'hook.event_patterns': patternsFieldsCustomizer,
    'hook.alarm_patterns': patternsFieldsCustomizer,
    'hook.entity_patterns': patternsFieldsCustomizer,
    enabled: enabled === undefined ? true : enabled,
    requests: request.map(({ headers, auth, ...otherRequestFields }) => ({
      headers: objectToTextPairs(headers),
      withAuth: !!auth,
      auth,
      ...otherRequestFields,
    })),
    postProcessors: postProcessors.map(({ parameters }) => {
      const { empty_response: emptyResponse, fields, ...postProcessor } = parameters;

      return {
        ...postProcessor,
        fields: objectToTextPairs(fields),
        emptyResponse,
      };
    }),
  });
}

/**
 *
 * @param requests
 * @return {Object}
 */
function formRequestFieldToWebhook(requests) {
  return requests.map(({
    emptyResponse, withAuth, auth, headers, ...otherFields
  }) => {
    const request = {
      empty_response: !!emptyResponse,
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
  return postProcessors.map(({
    fields,
    ticketId,
    ...postProcessorParameters
  }) => ({
    type: POST_PROCESSOR_TYPES.declareTicket,
    parameters: {
      fields: textPairsToObject(fields),
      ticket_id: ticketId,
      ...postProcessorParameters,
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
    request: formRequestFieldToWebhook(requests),
    post_processors: formPostProcessorsToWebhook(postProcessors),
    ...webhook,
  };
}
