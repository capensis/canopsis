<template>
  <c-select-chip
    v-field="value"
    :items="items"
    :aria-label="$t('llm.model')"
    :disabled="disabled || pending"
    :color="colors.chipBackground"
    :text-color="colors.chipText"
    item-text="name"
    item-value="name"
    active-list-item-class="primary--text"
    rounded
  >
    <template #selection>
      <v-progress-circular
        v-if="pending"
        color="primary"
        size="16"
        width="2"
        indeterminate
      />
      <span v-else>
        <span>{{ value }}</span>
        <v-icon
          v-if="items.length > 1"
          :color="colors.chipText"
          class="ml-1"
          small
        >
          keyboard_arrow_down
        </v-icon>
      </span>
    </template>
  </c-select-chip>
</template>

<script>
import { watch } from 'vue';

import { COLORS } from '@/config';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    items: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const colors = COLORS.aiChat;

    const applyDefaultLlm = () => {
      const list = props.items;

      if (!list.length) {
        if (props.value) {
          emit('input', '');
        }

        return;
      }

      const stillValid = list.some(llm => llm.name === props.value);

      if (stillValid) {
        return;
      }

      const defaultLlm = list.find(llm => llm.default);

      emit('input', defaultLlm?.name ?? list[0].name);
    };

    watch(() => props.items, applyDefaultLlm, { immediate: true });

    return {
      colors,
    };
  },
};
</script>
