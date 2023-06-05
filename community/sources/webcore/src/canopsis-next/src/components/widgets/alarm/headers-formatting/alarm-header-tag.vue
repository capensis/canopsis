<template lang="pug">
  span
    slot
    c-alarm-action-chip(
      v-if="selectedTag",
      :color="getTagColor(selectedTag)",
      closable,
      @close="$emit('clear')"
    ) {{ selectedTag }}
</template>

<script>
import { entitiesAlarmTagMixin } from '@/mixins/entities/alarm-tag';

export default {
  mixins: [entitiesAlarmTagMixin],
  props: {
    selectedTag: {
      type: String,
      default: '',
    },
  },
};
</script>
