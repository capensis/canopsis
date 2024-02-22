<template>
  <c-alarm-actions-chips
    :items="tags"
    :active-item="selectedTag"
    :small="small"
    :inline-count="inlineCount"
    item-text="text"
    item-value="text"
    row
    v-on="$listeners"
  />
</template>

<script>
import { entitiesAlarmTagMixin } from '@/mixins/entities/alarm-tag';

export default {
  mixins: [entitiesAlarmTagMixin],
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    selectedTag: {
      type: String,
      required: false,
    },
    inlineCount: {
      type: [Number, String],
      default: 2,
    },
    small: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    tags() {
      return (this.alarm.tags ?? []).map(tag => ({
        text: tag,
        color: this.getTagColor(tag),
      }));
    },
  },
};
</script>
