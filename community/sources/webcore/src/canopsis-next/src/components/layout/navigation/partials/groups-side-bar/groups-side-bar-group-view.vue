<template lang="pug">
  router-link.panel-item-content-link(
    :class="{ editing: isNavigationEditingMode }",
    :event="routerLinkEvents",
    :title="view.title",
    :to="viewLink"
  )
    group-view-panel(
      :view="view",
      :is-order-changed="isGroupsOrderChanged",
      :is-view-active="isViewActive",
      :editable="(isViewPrivate || hasViewEditButtonAccess) && isNavigationEditingMode",
      :duplicable="(isViewPrivate || hasCreateAnyViewAccess) && isNavigationEditingMode",
      @duplicate="showDuplicateViewModal",
      @change="showEditViewModal"
    )
</template>

<script>
import layoutNavigationGroupsBarGroupViewMixin from '@/mixins/layout/navigation/groups-bar-group-view';

import GroupViewPanel from './group-view-panel.vue';

export default {
  components: { GroupViewPanel },
  mixins: [layoutNavigationGroupsBarGroupViewMixin],
  props: {
    isGroupsOrderChanged: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isViewActive() {
      return this.$route?.params?.id === this.view._id;
    },

    isViewPrivate() {
      return this.view.is_private;
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
</style>
