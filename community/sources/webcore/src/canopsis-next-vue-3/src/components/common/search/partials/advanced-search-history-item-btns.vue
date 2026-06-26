<template>
  <c-layout gap="1">
    <c-action-btn
      :icon="pinned ? 'push_pin' : 'push_pin'"
      :tooltip="$t('common.pin')"
      @click="$emit('toggle-pin', id)"
    />
    <c-action-btn
      type="delete"
      @click="$emit('remove', id)"
    />
  </c-layout>
</template>

<script>
export default {
  props: {
    id: {
      type: String,
      required: true,
    },
    pinned: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['remove', 'toggle-pin'],
};
</script>
