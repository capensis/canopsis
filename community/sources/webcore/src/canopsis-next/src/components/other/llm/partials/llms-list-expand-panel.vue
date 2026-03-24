<template>
  <v-tabs
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>
      {{ $t('llm.expandTabs.promptsHistory') }}
    </v-tab>
    <v-tab-item>
      <v-layout class="py-3 secondary lighten-2">
        <v-flex xs12>
          <v-card>
            <v-card-text>
              <llm-prompts-history-panel :llm-id="llm._id" />
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-tab-item>
  </v-tabs>
</template>

<script>
import LlmPromptsHistoryPanel from './llm-prompts-history-panel.vue';

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
};
</script>
