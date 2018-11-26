import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal/modal';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';

export default {
  mixins: [
    modalMixin,
    entitiesViewGroupMixin,
    rightsTechnicalViewMixin,
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
  computed: {
    checkUpdateViewAccessById() {
      return viewId => this.checkUpdateAccess(viewId) && this.hasUpdateAnyViewAccess;
    },

    checkDeleteViewAccessById() {
      return viewId => this.checkDeleteAccess(viewId) && this.hasDeleteAnyViewAccess;
    },

    getAvailableViewsForGroup() {
      return group => group.views.filter(view => this.checkReadAccess(view._id));
    },
  },
  mounted() {
    this.fetchGroupsList();
  },
  methods: {
    toggleEditingMode() {
      this.isEditingMode = !this.isEditingMode;
    },

    showEditGroupModal(group) {
      this.showModal({
        name: MODALS.createGroup,
        config: { group },
      });
    },

    showEditViewModal(view) {
      this.showModal({
        name: MODALS.createView,
        config: { view },
      });
    },

    showDuplicateViewModal(view) {
      this.showModal({
        name: MODALS.createView,
        config: {
          view,
          isDuplicating: true,
        },
      });
    },
  },
};
