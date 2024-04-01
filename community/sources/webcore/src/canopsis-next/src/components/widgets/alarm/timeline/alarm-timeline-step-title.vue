<template>
  <span class="text-subtitle-2">{{ title }}</span>
</template>

<script>
import { computed } from 'vue';

import {
  ALARM_LIST_STEPS,
  ALARM_STEPS_WITH_AUTHOR_IN_TITLE,
  ALARM_STEPS_WITH_CONDITION_FOR_AUTHOR_IN_TITLE,
  ALARM_STEPS_WITH_LAUNCHED_IN_TITLE,
} from '@/constants';

import { convertDurationToString } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    step: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const { t } = useI18n();

    const title = computed(() => {
      const { _t: type, a: author, t: timestamp, val: value } = props.step;
      const hasLaunchedMessage = ALARM_STEPS_WITH_LAUNCHED_IN_TITLE.includes(type);
      const hasConditionForAuthorMessage = ALARM_STEPS_WITH_CONDITION_FOR_AUTHOR_IN_TITLE.includes(type);
      const hasAuthorMessage = ALARM_STEPS_WITH_AUTHOR_IN_TITLE.includes(type)
        && ((hasConditionForAuthorMessage && !props.child) || !hasConditionForAuthorMessage);
      const payload = { author };

      switch (type) {
        case ALARM_LIST_STEPS.statusdec:
        case ALARM_LIST_STEPS.statusinc:
          payload.status = t(`common.statusTypes.${value}`);
          break;
        case ALARM_LIST_STEPS.snooze:
          payload.duration = convertDurationToString(value - timestamp);
          break;
      }

      let result = t(`alarm.timeline.steps.${type}`, payload);

      if (hasLaunchedMessage) {
        result += `, ${t('alarm.timeline.launched')}`;
      }

      if (hasAuthorMessage) {
        result += ` ${t('alarm.timeline.by')} ${author}`;
      }

      return result;
    });

    return {
      title,
    };
  },
};
</script>
