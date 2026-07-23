<template>
  <v-tabs v-model="activeTab" color="primary" align-tabs="center">
    <v-tab v-for="tab in tabs" :key="tab.id" :value="tab.id">
      {{ tab.label }}
    </v-tab>
  </v-tabs>

  <v-card-text>
    <v-window v-model="activeTab">
      <v-window-item v-for="tab in tabs" :key="`${tab.id}-item`" :value="tab.id">
        <ticket-status-jobs-tab-item :tab-id="tab.id" />
      </v-window-item>
    </v-window>
  </v-card-text>
</template>

<script>
import { computed, ref } from 'vue';

import { JOBS_TABS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import TicketStatusJobsTabItem from './ticket-status-jobs-tab-item.vue';

export default {
  components: {
    TicketStatusJobsTabItem,
  },

  setup() {
    const { t } = useI18n();
    const tabs = computed(() =>
      Object.values(JOBS_TABS).map((tab) => ({
        id: tab,
        label: t(`jobs.tabs.${tab}`),
      })),
    );
    const activeTab = ref(tabs.value[0]?.id ?? JOBS_TABS.ticketStatus);

    return {
      activeTab,
      tabs,
    };
  },
};
</script>
