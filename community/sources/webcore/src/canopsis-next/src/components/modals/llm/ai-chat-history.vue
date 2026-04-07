<template>
  <modal-wrapper close>
    <template #title="">
      {{ $t('llm.chat.chatHistoryTitle') }}
    </template>
    <template #text="">
      <v-layout
        v-if="historyItem"
        column
      >
        <pre class="ai-chat-history-modal__payload">{{ historyJson }}</pre>
      </v-layout>
      <span
        v-else
        class="grey--text"
      >{{ $t('common.noData') }}</span>
    </template>
    <template #actions="">
      <v-btn
        depressed
        text
        @click="close"
      >
        {{ $t('common.close') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { useInnerModal } from '@/hooks/modals';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.aiChatHistory,
  components: {
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);

    const historyItem = computed(() => config.value.item ?? null);

    const historyJson = computed(() => (
      historyItem.value ? JSON.stringify(historyItem.value, null, 2) : ''
    ));

    return {
      historyItem,
      historyJson,
      close,
    };
  },
};
</script>

<style lang="scss" scoped>
.ai-chat-history-modal__payload {
  margin: 0;
  max-height: 60vh;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: monospace;
  font-size: 12px;
}
</style>
