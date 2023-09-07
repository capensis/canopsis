<template lang="pug">
  v-card(:color="cardColor")
    v-card-text.panel-item-content
      v-layout(align-center, justify-space-between)
        v-flex(:class="{ 'panel-view-title--editing': isEditing }")
          v-layout(align-center)
            span.pl-2(:class="{ ellipsis: ellipsis }")
              slot(name="title") {{ view.title }}
        v-flex
          v-layout(v-if="allowEditing", justify-end)
            v-btn.ma-0(
              v-show="hasEditAccess",
              :disabled="isOrderChanged",
              :data-test="`editViewButton-view-${view._id}`",
              depressed,
              small,
              icon,
              @click.prevent="editHandler"
            )
              v-icon(small) edit
            v-btn.ma-0(
              v-show="isEditing",
              :disabled="isOrderChanged",
              :data-test="`copyViewButton-view-${view._id}`",
              depressed,
              small,
              icon,
              @click.prevent="duplicateHandler"
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
    allowEditing: {
      type: Boolean,
      default: false,
    },
    hasEditAccess: {
      type: Boolean,
      default: false,
    },
    isEditing: {
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
      return this.isViewActive ? 'secondary white--text lighten-3' : 'secondary white--text lighten-1';
    },

    ellipsis() {
      return !this.$slots.title && !this.$scopedSlots.title;
    },
  },
  methods: {
    editHandler() {
      this.$emit('change');
    },
    duplicateHandler() {
      this.$emit('duplicate');
    },
  },
};
</script>

<style lang="scss" scoped>
  .panel-item-content {
    display: -webkit-box;
    display: -ms-flexbox;
    display: flex;
    cursor: pointer;
    -webkit-box-align: center;
    -ms-flex-align: center;
    align-items: center;
    position: relative;
    padding: 12px 24px;
    height: 48px;

    & > div {
      max-width: 100%;
    }

    & ::v-deep .v-btn:not(:last-child) {
      margin-right: 0;
    }
  }

  .panel-view-title {
    &--editing {
      max-width: 73%;
    }
  }
</style>
