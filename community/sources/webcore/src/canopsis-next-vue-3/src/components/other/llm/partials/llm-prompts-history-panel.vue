<template>
  <div class="pa-4 secondary lighten-2">
    <v-card>
      <v-card-text>
        <v-tabs v-model="activeTab" color="primary" align-tabs="center">
          <v-tab :value="TABS.allPrompts">
            {{ t('llm.promptsHistory.tabs.allPrompts') }}
          </v-tab>
          <v-tab :value="TABS.userHistory">
            {{ t('llm.promptsHistory.userHistoryTitle') }}
          </v-tab>
        </v-tabs>

        <v-window v-model="activeTab" touchless>
          <v-window-item :value="TABS.allPrompts">
            <llm-prompts-history-table :llm-id="llmId" />
          </v-window-item>

          <v-window-item :value="TABS.userHistory">
            <llm-prompts-history-by-user-table :llm-id="llmId" />
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import { ref } from 'vue';

import { useI18n } from '@/hooks/i18n';

import LlmPromptsHistoryByUserTable from './llm-prompts-history-by-user-table.vue';
import LlmPromptsHistoryTable from './llm-prompts-history-table.vue';

const TABS = {
  allPrompts: 'allPrompts',
  userHistory: 'userHistory',
};

export default {
  components: {
    LlmPromptsHistoryByUserTable,
    LlmPromptsHistoryTable,
  },

  props: {
    llmId: {
      type: String,
      required: true,
    },
  },

  setup() {
    const { t } = useI18n();
    const activeTab = ref(TABS.allPrompts);

    return {
      t,
      activeTab,
      TABS,
    };
  },
};
</script>
