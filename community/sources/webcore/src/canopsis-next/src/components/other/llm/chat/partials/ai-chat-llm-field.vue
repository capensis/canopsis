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
          v-if="!disabled"
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
import { COLORS } from '@/config';

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
  setup() {
    const colors = COLORS.aiChat;

    return {
      colors,
    };
  },
};
</script>
