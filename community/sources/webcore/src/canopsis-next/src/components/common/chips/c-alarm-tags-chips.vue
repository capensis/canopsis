<template>
  <c-alarm-actions-chips
    :items="preparedTags"
    :active-item="selectedTag"
    :small="small"
    :inline-count="inlineCount"
    :closable-active="closableActive"
    item-class="c-alarm-tags-chips__chip"
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
    closableActive: {
      type: Boolean,
      default: false,
    },
    nameFilter: {
      type: String,
      default: '',
    },
    regexFilter: {
      type: String,
      default: '',
    },
  },
  computed: {
    filteredTags() {
      let { tags } = this.alarm;

      if (this.nameFilter) {
        const regex = new RegExp(this.nameFilter);

        tags = tags.filter(tag => tag.match(regex));
      }

      if (this.regexFilter) {
        const regex = new RegExp(this.regexFilter, 'i');

        tags = tags.filter(tag => tag.match(regex));
      }

      return tags;
    },

    preparedTags() {
      return this.filteredTags.map(tag => ({
        text: tag,
        color: this.getTagColor(tag),
      }));
    },
  },
};
</script>

<style lang="scss">
.c-alarm-tags-chips__chip .v-chip__content {
  padding: 0 4px;
}
</style>
