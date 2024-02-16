<template>
  <group-panel
    :group="group"
    :order-changed="isGroupsOrderChanged"
    :is-editing="isNavigationEditingMode"
    :editable="hasViewGroupEditAccess && isNavigationEditingMode"
    @change="showEditGroupModal"
  >
    <c-draggable-list-field
      class="views-panel secondary lighten-1"
      :class="{ 'views-panel--empty': isGroupEmpty }"
      :value="group.views"
      :disabled="!isNavigationEditingMode"
      group="views"
      @input="changeViewsOrdering"
    >
      <groups-side-bar-group-view
        v-for="view in group.views"
        :key="view._id"
        :view="view"
        :is-groups-order-changed="isGroupsOrderChanged"
      />
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
