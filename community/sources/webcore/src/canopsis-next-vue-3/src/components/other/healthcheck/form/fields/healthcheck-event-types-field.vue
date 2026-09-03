<template>
  <c-autocomplete
    :model-value="value"
    :items="types"
    :label="$t('common.eventType')"
    :placeholder="!value.length ? $t('common.all') : ''"
    :persistent-placeholder="!value.length"
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
  >
    <template #chip="{ item, index, props }">
      <v-chip
        v-if="index < showCount"
        v-bind="props"
        size="small"
        closable
      >
        <span class="text-truncate">
          {{ getItemText(item) }}
        </span>
      </v-chip>
      <span v-else-if="index === showCount">
        +{{ value.length - showCount }} {{ $t('common.more') }}
      </span>
      <span v-else />
    </template>
  </c-autocomplete>
</template>

<script>
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

  computed: {
    types() {
      return Object.values(HEALTHCHECK_EVENT_TYPES).map(value => ({
        value,
        text: value,
      }));
    },
  },

  methods: {
    getItemText(item) {
      return item?.text ?? item;
    },
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
