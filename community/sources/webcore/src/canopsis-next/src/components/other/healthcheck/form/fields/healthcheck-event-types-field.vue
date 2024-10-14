<template>
  <c-event-type-field
    :value="value"
    :types="types"
    :placeholder="!value.length ? $t('common.all') : ''"
    :persistent-placeholder="!value.length"
    class="healthcheck-event-types-field"
    multiple
    clearable
    @input="$emit('input', $event)"
  >
    <template #selection="{ item, index, parent }">
      <v-chip
        v-if="index < showCount"
        small
        close
        @click:close="parent.onChipInput(item)"
      >
        <span class="text-truncate">{{ item }}</span>
      </v-chip>
      <span v-else-if="index === showCount">+{{ value.length - showCount }} {{ $t('common.more') }}</span>
      <span v-else />
    </template>
  </c-event-type-field>
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
  setup() {
    const types = Object.values(HEALTHCHECK_EVENT_TYPES).map(value => ({
      value,
      text: value,
    }));

    return {
      types,
    };
  },
};
</script>

<style lang="scss">
.healthcheck-event-types-field {
  & ::placeholder {
    .theme--dark & {
      color: var(--v-text-dark-primary) !important;
    }

    .theme--light & {
      color: var(--v-text-light-primary) !important;
    }
  }
}
</style>
