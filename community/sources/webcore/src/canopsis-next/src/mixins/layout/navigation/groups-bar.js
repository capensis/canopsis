import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { permissionsEntitiesGroupMixin } from '@/mixins/permissions/entities/group';

import layoutNavigationEditingModeMixin from '../../layout/navigation/editing-mode';

export default {
  mixins: [
    entitiesViewGroupMixin,
    permissionsEntitiesGroupMixin,
    layoutNavigationEditingModeMixin,
  ],
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
  mounted() {
    this.fetchAllGroupsListWithViewsWithCurrentUser();
  },
};
