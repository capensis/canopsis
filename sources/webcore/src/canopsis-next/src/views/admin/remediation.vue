<template lang="pug">
  v-container
    the-page-header {{ $t('common.instructions') }}
    v-layout(row, wrap)
      v-flex(xs12)
        v-card.ma-2
          v-tabs(v-model="activeTab", fixed-tabs, slider-color="primary")
            template
              v-tab(href="#instructions") {{ $t('remediation.tabs.instructions') }}
              v-tab-item(:value="$constants.REMEDIATION_TABS.instructions")
                v-card-text
                  span {{ $t('remediation.tabs.instructions') }}
            template
              v-tab(href="#configurations") {{ $t('remediation.tabs.configurations') }}
              v-tab-item(:value="$constants.REMEDIATION_TABS.configurations")
                v-card-text
                  span {{ $t('remediation.tabs.configurations') }}
            template
              v-tab(href="#jobs") {{ $t('remediation.tabs.jobs') }}
              v-tab-item(:value="$constants.REMEDIATION_TABS.jobs")
                v-card-text
                  span {{ $t('remediation.tabs.jobs') }}
    fab-buttons(@create="create", @refresh="refresh", :has-access="hasAccess")
      span {{ tooltipText }}
</template>

<script>
import { REMEDIATION_TABS } from '@/constants';

import FabButtons from '@/components/other/fab-buttons/fab-buttons.vue';

export default {
  components: {
    FabButtons,
  },
  data() {
    return {
      activeTab: REMEDIATION_TABS.instructions,
    };
  },
  computed: {
    tooltipText() {
      return '';
    },

    hasAccess() {
      return true;
    },
  },
  methods: {
    refresh() {
      switch (this.activeTab) {
        case REMEDIATION_TABS.instructions:
          this.fetchInstructionsList();
          break;
        case REMEDIATION_TABS.configurations:
          this.fetchConfigurationsList();
          break;
        case REMEDIATION_TABS.jobs:
          this.fetchJobsList();
          break;
      }
    },

    create() {
      switch (this.activeTab) {
        case REMEDIATION_TABS.instructions:
          this.showCreateInstructionModal();
          break;
        case REMEDIATION_TABS.configurations:
          this.showCreateConfigurationModal();
          break;
        case REMEDIATION_TABS.jobs:
          this.showCreateJobModal();
          break;
      }
    },

    fetchInstructionsList() {
    },

    fetchConfigurationsList() {
    },

    fetchJobsList() {
    },

    showCreateInstructionModal() {
    },

    showCreateConfigurationModal() {
    },

    showCreateJobModal() {
    },
  },
};
</script>
