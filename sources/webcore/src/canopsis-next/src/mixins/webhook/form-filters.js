import { get } from 'lodash';

import { setInSeveral, unsetInSeveralWithConditions } from '@/helpers/immutable';
import { textPairsToObject, objectToTextPairs } from '@/helpers/text-pairs';

export default {
  filters: {
    webhookToForm(webhook) {
      const patternsCustomizer = value => value || [];

      return setInSeveral(webhook, {
        declare_ticket: objectToTextPairs,
        'request.headers': objectToTextPairs,
        'hook.event_patterns': patternsCustomizer,
        'hook.alarm_patterns': patternsCustomizer,
        'hook.entity_patterns': patternsCustomizer,
      });
    },
    formToWebhook(form) {
      const patternsCondition = value => !value || !value.length;
      const hasAuth = get(form, 'request.auth');

      const pathValuesMap = {
        declare_ticket: textPairsToObject,
        'request.headers': textPairsToObject,
      };

      if (hasAuth) {
        pathValuesMap['request.auth'] = auth => ({
          username: auth.username ? auth.username : null,
          password: auth.password ? auth.password : null,
        });
      }

      const webhook = setInSeveral(form, pathValuesMap);

      return unsetInSeveralWithConditions(webhook, {
        'hook.event_patterns': patternsCondition,
        'hook.alarm_patterns': patternsCondition,
        'hook.entity_patterns': patternsCondition,
      });
    },
  },
};
