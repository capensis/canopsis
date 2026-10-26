<template>
  <v-tabs v-model="activeTab" color="primary" align-tabs="center">
    <v-tab v-for="tab in tabs" :key="tab.id" :value="tab.id">
      {{ tab.label }}
    </v-tab>
  </v-tabs>

  <v-card-text>
    <v-window v-model="activeTab">
      <v-window-item v-for="tab in tabs" :key="`${tab.id}-item`" :value="tab.id">
        <job-management-tab-item
          :ref="element => setTabElement(tab.id, element)"
          :tab-id="tab.id"
        />
      </v-window-item>
    </v-window>
  </v-card-text>
</template>

<script>
import { computed, ref, watch } from 'vue';

import {
  JOBS_TABS,
  JOBS_TABS_PERMISSIONS,
  USER_PERMISSIONS,
} from '@/constants';

import { useCurrentUserPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';

import JobManagementTabItem from './job-management-tab-item.vue';

export default {
  components: {
    JobManagementTabItem,
  },

  setup(_, { expose }) {
    const { t } = useI18n();
    const { checkAccess } = useCurrentUserPermissions();
    const tabs = computed(() => Object.values(JOBS_TABS)
      .filter(tab => checkAccess(USER_PERMISSIONS.api.general[JOBS_TABS_PERMISSIONS[tab]]))
      .map(tab => ({
          id: tab,
          label: t(`jobs.tabs.${tab}`),
        })),
    );
    const activeTab = ref(tabs.value[0]?.id ?? JOBS_TABS.ticketStatus);
    const tabByIdElement = {};
    const setTabElement = (tabId, element) => {
      if (element) {
        tabByIdElement[tabId] = element;

        return;
      }

      delete tabByIdElement[tabId];
    };
    const refresh = () => tabByIdElement[activeTab.value]?.refresh();

    watch(tabs, (newTabs) => {
      if (!newTabs.some(tab => tab.id === activeTab.value)) {
        activeTab.value = newTabs[0]?.id;
      }
    });

    expose({ refresh });

    return {
      activeTab,
      tabs,
      setTabElement,
    };
  },
};
</script>
