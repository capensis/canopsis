import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';

import layoutNavigationEditingModeMixin from './editing-mode';

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
      return this.groups.filter(({ views = [] }) => views.length || this.isNavigationEditingMode);
    },
  },
  mounted() {
    this.fetchAllGroupsListWithWidgetsWithCurrentUser();
  },
};
