<template>
  <v-layout class="gap-1">
    <v-chip
      v-for="chip in chips"
      :key="chip.key"
      :color="chip.color"
      class="cursor-pointer"
      text-color="white"
      small
    >
      <span class="text-ucfirst">{{ chip.text }}</span>
    </v-chip>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { COLORS } from '@/config';

import { useI18n } from '@/hooks/i18n';

const JOBS_TYPES = {
  running: 'running',
  completed: 'completed',
  failed: 'failed',
};

const JOBS_STATUSES_COLORS = {
  [JOBS_TYPES.completed]: COLORS.success,
  [JOBS_TYPES.failed]: COLORS.error,
  [JOBS_TYPES.running]: COLORS.warning,
};

const JOBS_STATUSES_TEXTS = {
  [JOBS_TYPES.completed]: 'common.finished',
  [JOBS_TYPES.failed]: 'common.failed',
  [JOBS_TYPES.running]: 'common.inProgress',
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

    const chips = computed(() => {
      const notEmptyJobs = Object.entries(props.jobs).filter(([, value]) => !!value);

      return notEmptyJobs.reduce((acc, [key, value]) => {
        if (value) {
          const textPrefix = notEmptyJobs.length === 1 && notEmptyJobs[0] !== JOBS_TYPES.running
            ? `${t('common.all')}`
            : `${value}`;

          acc.push({
            key,
            color: JOBS_STATUSES_COLORS[key],
            text: `${textPrefix} ${t(JOBS_STATUSES_TEXTS[key])}`,
          });
        }

        return acc;
      }, []);
    });

    return {
      chips,
    };
  },
};
</script>
