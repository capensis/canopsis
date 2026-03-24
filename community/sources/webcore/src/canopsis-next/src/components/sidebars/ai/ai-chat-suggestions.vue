<template>
  <div>
    <div class="text-body-2 font-weight-medium text-center mb-4">
      {{ $t('llm.chat.tryLabel') }}
    </div>
    <v-layout class="ai-chat-suggestions gap-2 mt-2" justify-space-between>
      <c-chip
        v-for="suggestion in suggestions"
        :key="suggestion.type"
        :aria-label="suggestion.label"
        :color="colors.chipBackground"
        :text-color="colors.chipText"
        class="ma-0 px-2"
        rounded
        @click="selectSuggestion(suggestion.type)"
      >
        <v-layout align-center>
          <v-icon
            v-if="suggestion.icon"
            class="mr-1"
            small
          >
            {{ suggestion.icon }}
          </v-icon>
          <span>{{ suggestion.label }}</span>
        </v-layout>
      </c-chip>
    </v-layout>
  </div>
</template>

<script>
import { computed } from 'vue';

import { COLORS } from '@/config';

import { useI18n } from '@/hooks/i18n';

export default {
  setup(props, { emit }) {
    const { t } = useI18n();
    const colors = COLORS.aiChat;

    /**
     * TODO: move it into constants in the future
     */
    const suggestions = computed(() => ([
      {
        type: 'createPattern',
        label: t('llm.chat.suggestions.createPattern'),
        icon: 'add',
      },
      {
        type: 'editPattern',
        label: t('llm.chat.suggestions.editPattern'),
        icon: 'edit',
      },
      {
        type: 'validatePattern',
        label: t('llm.chat.suggestions.validatePattern'),
        icon: 'check_circle',
      },
    ]));

    /**
     * Notifies the parent that the user picked a quick-action suggestion.
     *
     * @param {string} type - Suggestion key (`createPattern`, `editPattern`, or `validatePattern`).
     */
    const selectSuggestion = type => emit('select', type);

    return {
      colors,
      suggestions,
      selectSuggestion,
    };
  },
};
</script>
