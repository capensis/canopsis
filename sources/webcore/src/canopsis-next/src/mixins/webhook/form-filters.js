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
      return setInSeveral(form, {
        'request.headers': textPairsToObject,
        declare_ticket: textPairsToObject,
      });
    },
  },
};
