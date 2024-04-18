<template>
  <c-alarm-actions-chips
    :items="tags"
    :inline-count="inlineCount"
    item-class="c-entity-tags-chips__chip"
    item-text="text"
    item-value="text"
    row
  />
</template>

<script>
import { entitiesAlarmTagMixin } from '@/mixins/entities/alarm-tag';

export default {
  mixins: [entitiesAlarmTagMixin],
  props: {
    entity: {
      type: Object,
      required: true,
    },
    inlineCount: {
      type: [Number, String],
      default: 2,
    },
  },
  computed: {
    tags() {
      return (this.entity.tags ?? []).map(tag => ({
        text: tag,
        color: 'secondary',
      }));
    },
  },
};
</script>

<style lang="scss">
.c-entity-tags-chips__chip .v-chip__content {
  padding: 0 4px;
}
</style>
