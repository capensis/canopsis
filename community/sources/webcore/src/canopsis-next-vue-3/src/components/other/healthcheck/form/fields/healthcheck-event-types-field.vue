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
  />
</template>

<script>
import { HEALTHCHECK_EVENT_TYPES } from '@/constants'

export default {
  props: {
    value: {
      type: Array,
      required: true,
    },
  },

  emits: ['input'],

  computed: {
    types() {
      return Object.values(HEALTHCHECK_EVENT_TYPES).map((value) => ({
        value,
        text: value,
      }))
    },
  },
}
</script>

<style lang="scss">
.healthcheck-event-types-field {
  & ::placeholder {
    color: currentColor !important;
    opacity: 1;
  }
}
</style>
