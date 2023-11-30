import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';

import { layoutNavigationEditingModeMixin } from './editing-mode';

export const layoutNavigationGroupsBarMixin = {
  mixins: [
    entitiesViewGroupMixin,
    permissionsTechnicalViewMixin,
    layoutNavigationEditingModeMixin,
  ],
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    availableGroups() {
      if (this.isNavigationEditingMode) {
        return this.groups;
      }

      return this.groups.filter(({ views = [] }) => views.length);
    },
  },
  mounted() {
    this.fetchAllGroupsListWithWidgetsWithCurrentUser();
  },
};
