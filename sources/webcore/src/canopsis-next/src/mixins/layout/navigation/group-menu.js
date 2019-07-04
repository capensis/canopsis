import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import rightsEntitiesGroupMixin from '@/mixins/rights/entities/group';

export default {
  mixins: [
    modalMixin,
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
  computed: {
    getViewLink() {
      return (view = {}) => {
        const link = {
          name: 'view',
          params: { id: view._id },
        };

        if (view.tabs && view.tabs.length) {
          link.query = { tabId: view.tabs[0]._id };
        }

        return link;
      };
    },

    checkUpdateViewAccessById() {
      return viewId => this.checkUpdateAccess(viewId) && this.hasUpdateAnyViewAccess;
    },

    checkDeleteViewAccessById() {
      return viewId => this.checkDeleteAccess(viewId) && this.hasDeleteAnyViewAccess;
    },

    checkViewEditButtonAccessById() {
      return id =>
        (this.checkUpdateViewAccessById(id) || this.checkDeleteViewAccessById(id)) && this.isEditingMode;
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
