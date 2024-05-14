<template lang="pug">
  span(v-bind="$attrs")
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
  inheritAttrs: false,
  props: {
    selectedTag: {
      type: String,
      default: '',
    },
  },
};
</script>
