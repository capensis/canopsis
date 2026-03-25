<template>
  <c-select-chip
    v-field="value"
    :items="items"
    :aria-label="$t('llm.model')"
    :disabled="disabled || pending"
    :color="colors.chipBackground"
    :text-color="colors.chipText"
    item-text="name"
    item-value="_id"
    active-list-item-class="primary--text"
    return-object
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
        <span>{{ value?.name }}</span>
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

import { useModelField } from '@/hooks/form/model-field';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: false,
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

    const { updateModel } = useModelField(props, emit);

    watch(() => props.items, () => {
      const defaultLlm = props.items.find(llm => llm.default);

      if (defaultLlm) {
        updateModel(defaultLlm);
      }
    }, { immediate: true });

    return {
      colors,
    };
  },
};
</script>
