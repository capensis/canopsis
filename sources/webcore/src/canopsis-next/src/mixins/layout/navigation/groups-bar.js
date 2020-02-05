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
  data() {
    return {
      isEditingMode: false,
    };
  },
  mounted() {
    this.fetchGroupsList();
  },
  methods: {
    toggleEditingMode() {
      this.isEditingMode = !this.isEditingMode;
    },
  },
};
