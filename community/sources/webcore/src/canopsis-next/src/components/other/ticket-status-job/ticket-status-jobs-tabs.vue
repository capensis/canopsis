<template>
  <v-card>
    <v-tabs
      slider-color="primary"
      centered
    >
      <template v-for="tab in tabs">
        <v-tab :key="tab.id" :href="`#${tab.id}`">
          {{ $t(tab.label) }}
        </v-tab>
        <v-tab-item :key="`${tab.id}-item`" :value="tab.id">
          <ticket-status-jobs-tab-item :tab-id="tab.id" />
        </v-tab-item>
      </template>
    </v-tabs>
  </v-card>
</template>

<script>
import { computed } from 'vue';

import { JOBS_TABS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import TicketStatusJobsTabItem from '@/components/other/ticket-status-job/ticket-status-jobs-tab-item.vue';

export default {
  components: { TicketStatusJobsTabItem },
  setup() {
    const { t } = useI18n();

    const tabs = computed(() => Object.values(JOBS_TABS).map(tab => ({
      id: tab,
      label: t(`jobs.tabs.${tab}`),
    })));

    return {
      tabs,
    };
  },
};
</script>
