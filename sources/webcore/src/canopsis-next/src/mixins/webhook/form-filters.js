import { setInSeveral } from '@/helpers/immutable';
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
