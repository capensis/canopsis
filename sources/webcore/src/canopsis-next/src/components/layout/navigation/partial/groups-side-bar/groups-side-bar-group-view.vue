<template lang="pug">
  v-card(:color="cardColor")
    router-link.panel-item-content-link(
      :class="{ editing: isNavigationEditingMode }",
      :event="routerLinkEvents",
      :data-test="`linkView-view-${view._id}`",
      :title="view.title",
      :to="viewLink"
    )
      v-card-text.panel-item-content
        v-layout(align-center, justify-space-between)
          v-flex
            v-layout(align-center)
              span.pl-2 {{ view.title }}
          v-flex
            v-layout(justify-end)
              v-btn.ma-0(
                v-show="hasViewEditButtonAccess",
                :disabled="isGroupsOrderChanged",
                :data-test="`editViewButton-view-${view._id}`",
                depressed,
                small,
                icon,
                @click.prevent="showEditViewModal"
              )
                v-icon(small) edit
              v-btn.ma-0(
                v-show="isNavigationEditingMode",
                :disabled="isGroupsOrderChanged",
                :data-test="`copyViewButton-view-${view._id}`",
                depressed,
                small,
                icon,
                @click.prevent="showDuplicateViewModal"
              )
                v-icon(small) file_copy
    v-divider
</template>

<script>
import layoutNavigationGroupsBarGroupViewMixin from '@/mixins/layout/navigation/groups-bar-group-view';

export default {
  mixins: [layoutNavigationGroupsBarGroupViewMixin],
  props: {
    isGroupsOrderChanged: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isViewActive() {
      return this.$route.params.id && this.$route.params.id === this.view._id;
    },

    cardColor() {
      return this.isViewActive ? 'secondary white--text lighten-3' : 'secondary white--text lighten-1';
    },

    routerLinkEvents() {
      return this.isGroupsOrderChanged ? [] : ['click'];
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

    &.editing {
      cursor: move;

      .panel-item-content {
        cursor: inherit;
      }
    }
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
