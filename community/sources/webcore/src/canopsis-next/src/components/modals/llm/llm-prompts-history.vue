<template>
  <c-modal-wrapper close>
    <template #title>
      {{ t('llm.expandTabs.promptsHistory') }}
    </template>

    <template #text>
      <llm-prompts-history-panel :llm-id="llmId" />
    </template>

    <template #actions>
      <v-btn
        variant="text"
        @click="close"
      >
        {{ t('common.close') }}
      </v-btn>
    </template>
  </c-modal-wrapper>
</template>

<script setup>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';

import LlmPromptsHistoryPanel from '@/components/other/llm/partials/llm-prompts-history-panel.vue';

defineOptions({
  name: MODALS.llmPromptsHistory,
});

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const { config, close } = useInnerModal(props);

const llmId = computed(() => config.value.llmId);
</script>
