<template>
  <router-link v-slot="{ href, navigate }" :to="viewLink" custom>
    <a
      :class="{ editing: isNavigationEditingMode }"
      :href="href"
      :title="view.title"
      class="panel-item-content-link"
      @click="handleClick($event, navigate)"
    >
      <group-view-panel
        :view="view"
        :is-order-changed="isGroupsOrderChanged"
        :is-view-active="isViewActive"
        :editable="(isViewPrivate || hasViewEditButtonAccess) && isNavigationEditingMode"
        :duplicable="(isViewPrivate || hasCreateAnyViewAccess) && isNavigationEditingMode"
        @duplicate="showDuplicateViewModal"
        @change="showEditViewModal"
      />
    </a>
  </router-link>
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
  },
  methods: {
    handleClick(event, navigate) {
      if (this.isGroupsOrderChanged) {
        event.preventDefault();

        return;
      }

      navigate(event);
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

    :deep(.panel-item-content) {
      cursor: inherit;
    }
  }
}
</style>
