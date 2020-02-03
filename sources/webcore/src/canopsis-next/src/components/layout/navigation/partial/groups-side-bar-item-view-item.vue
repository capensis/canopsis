<template lang="pug">
  v-card(:color="cardColor")
    router-link.panel-item-content-link(
      :data-test="`linkView-view-${view._id}`",
      :title="view.title",
      :to="getViewLink(view)"
    )
      v-card-text.panel-item-content
        v-layout(align-center, justify-space-between)
          v-flex
            v-layout(align-center)
              span.pl-2 {{ view.title }}
          v-flex
            v-layout(justify-end)
              v-btn.ma-0(
                :data-test="`editViewButton-view-${view._id}`",
                v-show="checkViewEditButtonAccessById(view._id)",
                depressed,
                small,
                icon,
                @click.prevent="showEditViewModal(view)"
              )
                v-icon(small) edit
              v-btn.ma-0(
                :data-test="`copyViewButton-view-${view._id}`",
                v-show="isEditingMode",
                depressed,
                small,
                icon,
                @click.prevent="showDuplicateViewModal(view)"
              )
                v-icon(small) file_copy
    v-divider
</template>

<script>
import layoutNavigationViewItemMixin from '@/mixins/layout/navigation/view-item';

export default {
  mixins: [layoutNavigationViewItemMixin],
  props: {
    view: {
      type: Object,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    cardColor() {
      return this.isViewActive(this.view._id) ? 'secondary white--text lighten-3' : 'secondary white--text lighten-1';
    },
  },
  methods: {
    isViewActive(viewId) {
      return this.$route.params.id && this.$route.params.id === viewId;
    },
  },
};
</script>

<style lang="scss" scoped>
  a {
    color: inherit;
    text-decoration: none;
  }

  .panel-item-content-link {
    width: 100%;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    display: inline-block;
    vertical-align: middle;
  }

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

    & /deep/ .v-btn:not(:last-child) {
      margin-right: 0;
    }
  }
</style>
