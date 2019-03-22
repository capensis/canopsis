import { setInSeveral } from '@/helpers/immutable';
import { textPairsToObject, objectToTextPairs } from '@/helpers/text-pairs';

export default {
  filters: {
    webhookToForm(webhook) {
      return setInSeveral(webhook, {
        'request.headers': objectToTextPairs,
        declare_ticket: objectToTextPairs,
      });
    },
    formToWebhook(form) {
      const patternsCustomizer = value => (value && value.length ? value : null);

      return setInSeveral(form, {
        declare_ticket: textPairsToObject,
        'request.headers': textPairsToObject,
        'hook.event_patterns': patternsCustomizer,
        'hook.alarm_patterns': patternsCustomizer,
        'hook.entity_patterns': patternsCustomizer,
      });
    },
  },
};
