<template>
  <div class="pa-4 secondary lighten-2">
    <v-card>
      <v-card-text>
        <v-tabs
          slider-color="primary"
          centered
        >
          <v-tab>{{ $t('llm.promptsHistory.tabs.allPrompts') }}</v-tab>
          <v-tab-item>
            <llm-prompts-history-table :llm-id="llmId" />
          </v-tab-item>

          <v-tab>{{ $t('llm.promptsHistory.userHistoryTitle') }}</v-tab>
          <v-tab-item>
            <llm-prompts-history-by-user-table :llm-id="llmId" />
          </v-tab-item>
        </v-tabs>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import LlmPromptsHistoryTable from './llm-prompts-history-table.vue';
import LlmPromptsHistoryByUserTable from './llm-prompts-history-by-user-table.vue';

export default {
  components: {
    LlmPromptsHistoryTable,
    LlmPromptsHistoryByUserTable,
  },
  props: {
    llmId: {
      type: String,
      required: true,
    },
  },
};
</script>
