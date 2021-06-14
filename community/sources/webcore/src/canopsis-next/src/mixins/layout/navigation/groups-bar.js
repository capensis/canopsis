import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import rightsEntitiesGroupMixin from '@/mixins/rights/entities/group';

export default {
  mixins: [
    entitiesViewGroupMixin,
    rightsEntitiesGroupMixin,
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
