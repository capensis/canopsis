<template>
  <span
    :class="{ 'error--text': !!errorMessage, 'text--secondary': !errorMessage }"
    class="text-caption"
  >
    <span v-html="errorMessage || message" />
  </span>
</template>

<script>
import { computed } from 'vue';

import { TIME_UNITS, PATTERN_DURATION_FORMAT } from '@/constants';

import { sanitizeHtml } from '@/helpers/html';
import { convertDurationToString } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    alarmCounter: {
      type: Object,
      required: false,
    },
    entityCounter: {
      type: Object,
      required: false,
    },
    errorMessage: {
      type: String,
      required: false,
    },
  },
  setup(props) {
    const { t, tc } = useI18n();

    const message = computed(() => {
      const messages = [];

      if (props.alarmCounter) {
        const alarmCount = props.alarmCounter?.count ?? 0;

        messages.push(tc('pattern.alarmFound', alarmCount, { count: alarmCount }));
      }

      if (props.entityCounter) {
        const entityCount = props.entityCounter?.count ?? 0;

        messages.push(tc('pattern.entityFound', entityCount, { count: entityCount }));
      }

      if (!messages.length) {
        return '';
      }

      const searchTime = convertDurationToString(
        Math.max(props.alarmCounter?.ms ?? 0, props.entityCounter?.ms ?? 0),
        PATTERN_DURATION_FORMAT,
        TIME_UNITS.millisecond,
      );

      return sanitizeHtml(t('pattern.found', { message: messages.join(', '), searchTime }));
    });

    return {
      message,
    };
  },
};
</script>
