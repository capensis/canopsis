<template lang="pug">
  router-link.panel-item-content-link(
    :class="{ editing: isNavigationEditingMode }",
    :event="routerLinkEvents",
    :data-test="`linkView-view-${view._id}`",
    :title="view.title",
    :to="viewLink"
  )
    group-view-panel(
      :view="view",
      :is-editing="isNavigationEditingMode",
      :is-order-changed="isGroupsOrderChanged",
      :is-view-active="isViewActive",
      :has-edit-access="hasViewEditButtonAccess",
      allow-editing,
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
      return this.$route.params.id && this.$route.params.id === this.view._id;
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
