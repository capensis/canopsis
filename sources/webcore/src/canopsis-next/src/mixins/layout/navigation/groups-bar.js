import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import rightsEntitiesGroupMixin from '@/mixins/rights/entities/group';

import layoutNavigationEditingModeMixin from './editing-mode';

export default {
  mixins: [
    entitiesViewGroupMixin,
    rightsEntitiesGroupMixin,
    layoutNavigationEditingModeMixin,
  ],
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
  mounted() {
    this.fetchGroupsList();
  },
};
