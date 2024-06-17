<template>
  <span v-bind="$attrs">
    <slot />
    <c-alarm-action-chip
      v-if="selectedTag"
      :color="getTagColor(selectedTag)"
      class="px-2"
      closable
      @close="$emit('clear')"
    >
      {{ selectedTag }}
    </c-alarm-action-chip>
  </span>
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
