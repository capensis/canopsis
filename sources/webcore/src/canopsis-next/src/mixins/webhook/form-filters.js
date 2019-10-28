import { get, omit } from 'lodash';

import { setInSeveral, unsetInSeveralWithConditions } from '@/helpers/immutable';
import { textPairsToObject, objectToTextPairs } from '@/helpers/text-pairs';

export default {
  filters: {
    webhookToForm(webhook) {
      const patternsCustomizer = value => value || [];

      let form = {};

      const webhookFields = {
        'request.headers': objectToTextPairs,
        'hook.event_patterns': patternsCustomizer,
        'hook.alarm_patterns': patternsCustomizer,
        'hook.entity_patterns': patternsCustomizer,
      };

      if (webhook.declare_ticket) {
        const declareTicket = omit(webhook.declare_ticket, ['empty_response']);
        webhookFields.declare_ticket = () => objectToTextPairs(declareTicket);
      }

      if (webhook.declare_ticket && webhook.declare_ticket.empty_response) {
        form.empty_response = webhook.declare_ticket.empty_response;
      }

      form = {
        ...form,
        ...setInSeveral(webhook, webhookFields),
      };

      return form;
    },
    formToWebhook(form) {
      const patternsCondition = value => !value || !value.length;
      const hasAuth = get(form, 'request.auth');

      const pathValuesMap = {
        'request.headers': textPairsToObject,
      };

      if (form.declare_ticket) {
        pathValuesMap.declare_ticket = (value) => {
          const newValue = textPairsToObject(value);

          newValue.empty_response = form.emptyResponse;

          return newValue;
        };
      }

      if (hasAuth) {
        pathValuesMap['request.auth'] = auth => ({
          username: auth.username ? auth.username : null,
          password: auth.password ? auth.password : null,
        });
      }

      const webhook = setInSeveral(omit(form, ['emptyResponse']), pathValuesMap);

      return unsetInSeveralWithConditions(webhook, {
        'hook.event_patterns': patternsCondition,
        'hook.alarm_patterns': patternsCondition,
        'hook.entity_patterns': patternsCondition,
      });
    },
  },
};
