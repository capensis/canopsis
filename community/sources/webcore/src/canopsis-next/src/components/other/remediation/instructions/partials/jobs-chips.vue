<template>
  <v-row class="my-2">
    <v-col cols="12">
      <span class="text-subtitle-1 font-weight-bold">
        {{ t('remediation.instruction.job', 2) }}:
      </span>
    </v-col>
    <v-col v-for="job in jobs" :key="job._id" cols="auto">
      <c-instruction-job-chip>
        {{ job.name }}
      </c-instruction-job-chip>
    </v-col>
  </v-row>
</template>

<script>
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    jobs: {
      type: Array,
      default: () => [],
    },
  },
  setup() {
    const { t } = useI18n();

    return {
      t,
    };
  },
};
</script>
