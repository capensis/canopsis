<template>
  <c-alarm-actions-chips
    :items="alarm.tag_colors"
    :active-item="selectedTag"
    :small="small"
    :inline-count="inlineCount"
    :closable-active="closableActive"
    item-class="c-alarm-tags-chips__chip"
    item-text="value"
    item-value="value"
    row
    v-on="$listeners"
  />
</template>

<script>
export default {
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
  },
};
</script>

<style lang="scss">
.c-alarm-tags-chips__chip .v-chip__content {
  padding: 0 4px;
}
</style>
