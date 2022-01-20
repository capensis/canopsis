import { permissionsTechnicalViewMixin } from '../technical/view';
import layoutNavigationEditingModeMixin from '../../layout/navigation/editing-mode';

export const permissionsEntitiesGroupMixin = {
  mixins: [
    permissionsTechnicalViewMixin,
    layoutNavigationEditingModeMixin,
  ],
  computed: {
    availableGroups() {
      return this.groups.filter(({ views = [] }) => views.length || this.isNavigationEditingMode);
    },
  },
};
