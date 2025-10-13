<template>
  <c-alarm-actions-chips
    :items="entity.tag_colors"
    :inline-count="inlineCount"
    item-class="c-entity-tags-chips__chip"
    item-text="value"
    item-value="value"
    row
  />
</template>

<script>
export default {
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
};
</script>

<style lang="scss">
.c-entity-tags-chips__chip .v-chip__content {
  padding: 0 4px;

  * {
    color: white;
  }
}
</style>
