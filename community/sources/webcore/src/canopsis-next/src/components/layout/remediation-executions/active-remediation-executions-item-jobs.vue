<template>
  <v-layout class="gap-2">
    <v-chip
      v-for="chip in chips"
      :key="chip.key"
      :color="chip.color"
      text-color="white"
    >
      {{ chip.text }}
    </v-chip>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { COLORS } from '@/config';

import { useI18n } from '@/hooks/i18n';

const JOBS_STATUSES_COLORS = {
  completed: COLORS.success,
  failed: COLORS.error,
  running: COLORS.warning,
};

const JOBS_STATUSES_TEXTS = {
  completed: 'common.finished',
  failed: 'common.failed',
  running: 'common.inProgress',
};

export default {
  props: {
    jobs: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();

    const chips = computed(() => Object.entries(props.jobs).reduce((acc, [key, value]) => {
      if (value) {
        acc.push({
          key,
          color: JOBS_STATUSES_COLORS[key],
          text: `${value} ${t(JOBS_STATUSES_TEXTS[key])}`,
        });
      }

      return acc;
    }));

    return {
      chips,
    };
  },
};
</script>
