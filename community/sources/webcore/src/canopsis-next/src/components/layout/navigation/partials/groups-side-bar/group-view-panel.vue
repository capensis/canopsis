<template lang="pug">
  v-card(:color="cardColor")
    v-card-text.panel-item-content
      div.panel-item-content__title.pl-2(:class="{ ellipsis: ellipsis }")
        slot(name="title") {{ view.title }}
      div.panel-item-content__actions
        v-icon.mr-2(v-if="view.is_private", small) lock
        template(v-if="editable || duplicable")
          v-btn.ma-0(
            v-if="editable",
            :disabled="isOrderChanged",
            depressed,
            small,
            icon,
            @click.prevent="$emit('change')"
          )
            v-icon(small) edit
          v-btn.ma-0(
            v-if="duplicable",
            :disabled="isOrderChanged",
            depressed,
            small,
            icon,
            @click.prevent="$emit('duplicate')"
          )
            v-icon(small) file_copy
    v-divider
</template>

<script>
export default {
  props: {
    view: {
      type: Object,
      required: true,
    },
    editable: {
      type: Boolean,
      default: false,
    },
    duplicable: {
      type: Boolean,
      default: false,
    },
    isOrderChanged: {
      type: Boolean,
      default: false,
    },
    isViewActive: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    cardColor() {
      return `secondary ${this.isViewActive ? 'lighten-3' : 'lighten-1'}`;
    },

    ellipsis() {
      return !this.$slots.title && !this.$scopedSlots.title;
    },
  },
};
</script>

<style lang="scss" scoped>
  .panel-item-content {
    display: flex;
    cursor: pointer;
    align-items: center;
    justify-content: space-between;
    position: relative;
    padding: 12px 24px;
    height: 48px;

    &__title {
      width: 100%;
    }

    &__actions {
      flex-shrink: 0;

      display: flex;
      align-items: center;
    }

    & ::v-deep .v-btn:not(:last-child) {
      margin-right: 0;
    }
  }
</style>
