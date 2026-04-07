<template>
  <v-layout column>
    <v-tabs
      slider-color="primary"
      centered
    >
      <v-tab>{{ $t('llm.promptsHistory.tabs.allPrompts') }}</v-tab>
      <v-tab-item>
        <llm-prompts-history-all-prompts-table
          :llm-id="llmId"
          @see-chat="seeChat"
        />
      </v-tab-item>

      <v-tab>{{ $t('llm.promptsHistory.userHistoryTitle') }}</v-tab>
      <v-tab-item>
        <llm-prompts-history-by-user
          :llm-id="llmId"
          @see-chat="seeChat"
        />
      </v-tab-item>
    </v-tabs>
  </v-layout>
</template>

<script>
import LlmPromptsHistoryAllPromptsTable from './llm-prompts-history-all-prompts-table.vue';
import LlmPromptsHistoryByUser from './llm-prompts-history-by-user.vue';

export default {
  components: {
    LlmPromptsHistoryAllPromptsTable,
    LlmPromptsHistoryByUser,
  },
  props: {
    llmId: {
      type: String,
      required: true,
    },
  },
  setup(_, { emit }) {
    const seeChat = item => emit('see-chat', item);

    return {
      seeChat,
    };
  },
};
</script>
