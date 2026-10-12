<template>
  <c-autocomplete
    :model-value="value"
    :items="types"
    :label="$t('common.eventType')"
    :placeholder="!value.length ? $t('common.all') : ''"
    :persistent-placeholder="!value.length"
    :show-count="showCount"
    class="healthcheck-event-types-field"
    item-title="text"
    item-value="value"
    density="compact"
    multiple
    clearable
    chips
    closable-chips
    hide-details
    @update:model-value="$emit('input', $event)"
  />
</template>

<script>
import { computed } from 'vue';

import { HEALTHCHECK_EVENT_TYPES } from '@/constants';

export default {
  props: {
    value: {
      type: Array,
      required: true,
    },
    showCount: {
      type: Number,
      default: 2,
    },
  },

  emits: ['input'],

  setup() {
    const types = computed(() => Object.values(HEALTHCHECK_EVENT_TYPES).map(value => ({
      value,
      text: value,
    })));

    return {
      types,
    };
  },
};
</script>

<style lang="scss">
.healthcheck-event-types-field {
  & ::placeholder {
    color: currentColor !important;
    opacity: 1;
  }
}
</style>
