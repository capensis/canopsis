<template>
  <c-event-type-field
    :value="value"
    class="healthcheck-last-hour-filters__event-types"
    multiple
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
};
</script>
