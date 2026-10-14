<template>
  <c-layout
    :class="{ 'text-error': isFailedStatus }"
    class="my-2"
    align="center"
  >
    <v-icon :color="iconColor" class="mr-1">
      alarm
    </v-icon>
    <span>{{ label }}</span>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';
import { convertDurationToString, toSeconds } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    failedAt: {
      type: Number,
      required: false,
    },
    completedAt: {
      type: Number,
      required: false,
    },
    startedAt: {
      type: Number,
      required: false,
    },
    timeToComplete: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();

    const isFailedStatus = computed(() => !!props.failedAt);
    const iconColor = computed(() => (isFailedStatus.value ? 'error' : ''));

    const convertDate = date => convertDateToString(date, 'dateTimePicker');

    const label = computed(() => {
      if (props.completedAt) {
        return t('remediation.instructionExecute.completedAt', {
          time: convertDate(props.completedAt),
        });
      }

      if (props.failedAt) {
        return t('remediation.instructionExecute.failedAt', {
          time: convertDate(props.failedAt),
        });
      }

      if (props.startedAt) {
        return t('remediation.instructionExecute.startedAt', {
          time: convertDate(props.startedAt),
        });
      }

      const { value, unit } = props.timeToComplete ?? {};

      return t('remediation.instructionExecute.timeToComplete', {
        duration: convertDurationToString(toSeconds(value, unit)),
      });
    });

    return {
      isFailedStatus,
      iconColor,
      label,
    };
  },
};
</script>
