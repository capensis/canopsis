<template>
  <div class="llms-list-expand-panel">
    <v-tabs
      v-model="activeTab"
      bg-color="secondary"
      color="primary"
      align-tabs="center"
    >
      <v-tab :value="TABS.promptsHistory">
        {{ t('llm.expandTabs.promptsHistory') }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab" touchless>
      <v-window-item :value="TABS.promptsHistory">
        <llm-prompts-history-panel :llm-id="llm._id" />
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { ref } from 'vue';

import { useI18n } from '@/hooks/i18n';

import LlmPromptsHistoryPanel from './llm-prompts-history-panel.vue';

const TABS = {
  promptsHistory: 'promptsHistory',
};

export default {
  components: {
    LlmPromptsHistoryPanel,
  },

  props: {
    llm: {
      type: Object,
      required: true,
    },
  },

  setup() {
    const { t } = useI18n();
    const activeTab = ref(TABS.promptsHistory);

    return {
      t,
      activeTab,
      TABS,
    };
  },
};
</script>
