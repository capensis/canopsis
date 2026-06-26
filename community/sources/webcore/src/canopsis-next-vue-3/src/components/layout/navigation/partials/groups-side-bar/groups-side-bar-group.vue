<template>
  <group-panel
    :group="group"
    :order-changed="isGroupsOrderChanged"
    :is-editing="isNavigationEditingMode"
    :active="isActiveGroup"
    :editable="hasViewGroupEditAccess && isNavigationEditingMode"
    @change="showEditGroupModal"
  >
    <c-draggable-list-field
      :class="{ 'views-panel--empty': isGroupEmpty }"
      :model-value="group.views"
      :disabled="!isNavigationEditingMode"
      class="views-panel"
      group="views"
      @update:model-value="changeViewsOrdering"
    >
      <template #item="{ element: view }">
        <groups-side-bar-group-view
          :view="view"
          :is-groups-order-changed="isGroupsOrderChanged"
        />
      </template>
    </c-draggable-list-field>
  </group-panel>
</template>

<script>
import { layoutNavigationGroupsBarGroupMixin } from '@/mixins/layout/navigation/groups-bar-group';

import GroupsSideBarGroupView from './groups-side-bar-group-view.vue';
import GroupPanel from './group-panel.vue';

export default {
  components: { GroupPanel, GroupsSideBarGroupView },
  mixins: [layoutNavigationGroupsBarGroupMixin],
  props: {
    isGroupsOrderChanged: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:group'],
  computed: {
    isGroupEmpty() {
      return this.group.views && this.group.views.length === 0;
    },

    hasViewGroupEditAccess() {
      return this.group.is_private
        || this.hasDeleteAnyViewAccess
        || this.hasUpdateAnyViewAccess;
    },
  },
  methods: {
    changeViewsOrdering(views) {
      this.$emit('update:group', { ...this.group, views });
    },
  },
};
</script>

<style lang="scss" scoped>
.views-panel {
  background-color: var(--v-secondary-lighten1);

  &--empty {
    &:after {
      content: '';
      display: block;
      height: 48px;
      border: 4px dashed #4f6479;
      border-radius: 5px;
      position: relative;
    }
  }
}
</style>
