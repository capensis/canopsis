<template>
  <v-layout :class="{ 'ai-chat-pattern--active': active }" class="ai-chat-pattern pa-4" column>
    <v-layout class="gap-2">
      <v-flex>
        <v-layout class="ai-chat-pattern__content grey--text" column>
          <strong>
            <div class="ai-chat-pattern__patterns-text text-ucfirst">{{ patternsText }}</div>
            <span v-if="originalVersionText">{{ originalVersionText }}</span>
          </strong>
          <span>{{ versionText }}</span>
          <div>
            <v-btn
              class="ai-chat-pattern__see-pattern-btn mt-2"
              color="primary"
              text
              depressed
              @click="toggleExpanded"
            >
              <v-icon class="mr-2" color="primary">
                {{ expanded ? 'keyboard_arrow_up' : 'keyboard_arrow_down' }}
              </v-icon>
              <span>{{ expanded ? $t('llm.chat.pattern.hidePattern') : $t('llm.chat.pattern.seePattern') }}</span>
            </v-btn>
          </div>
        </v-layout>
      </v-flex>
      <div v-if="!active">
        <c-action-btn
          :tooltip="$t('llm.chat.pattern.restoreVersion')"
          icon="refresh"
          color="primary"
          class="ai-chat-pattern__reset-version-btn"
          left
          @click="restoreVersion"
        />
      </div>
    </v-layout>
    <v-expand-transition>
      <div v-if="expanded">
        <div class="pt-3 text-caption">
          <pre class="ai-chat-pattern__json pa-4">{{ parsedJson }}</pre>
        </div>
      </div>
    </v-expand-transition>
  </v-layout>
</template>

<script>
import { isNil } from 'lodash';
import { computed, ref } from 'vue';

import { stringifyJson } from '@/helpers/json';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    patterns: {
      type: Object,
      required: true,
    },
    version: {
      type: Number,
      default: 1,
    },
    originalVersion: {
      type: Number,
      required: false,
    },
    active: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { t, tc } = useI18n();

    const expanded = ref(false);

    const patternsText = computed(() => {
      const texts = Object.keys(props.patterns).map(key => t(`pattern.patternsFields.${key}`));

      return tc('llm.chat.patternsMessage', texts.length, { patterns: texts.join(', ') });
    });

    const originalVersionText = computed(() => (
      isNil(props.originalVersion) ? '' : ` - ${t('llm.chat.pattern.versionRestored', { version: props.originalVersion + 1 })}`
    ));

    const versionText = computed(() => {
      const currentText = props.active ? ` (${t('common.current')})` : '';

      return `${t('llm.chat.pattern.version', { version: props.version + 1 })}${currentText}`;
    });

    const parsedJson = computed(() => stringifyJson(props.patterns, 2));

    /**
     * Toggles visibility of the JSON pattern preview (`v-expand-transition`).
     */
    const toggleExpanded = () => expanded.value = !expanded.value;

    /**
     * Emits `restore:version` with this card’s `version` so the parent can restore that snapshot.
     */
    const restoreVersion = () => emit('restore:version', props.version);

    return {
      expanded,
      patternsText,
      originalVersionText,
      versionText,
      parsedJson,
      toggleExpanded,
      restoreVersion,
    };
  },
};
</script>

<style lang="scss" scoped>
.ai-chat-pattern {
  --inactive-border-color: #E0E0E0;

  border: 1px solid var(--inactive-border-color);
  border-radius: 8px;

  & &__see-pattern-btn {
    padding-left: 6px;
  }

  & &__reset-version-btn {
    height: auto;
  }

  & &__json {
    border-radius: 8px;
    border: 1px solid var(--v-primary-base);
    background-color: var(--v-grey-lighten-4, #f5f5f5);
    overflow-x: auto;
  }

  &--active {
    border-color: var(--v-primary-base);

    .ai-chat-pattern__content {
      color: var(--v-primary-base) !important;
    }
  }

  &__patterns-text {
    display: inline-block;
  }
}
</style>
