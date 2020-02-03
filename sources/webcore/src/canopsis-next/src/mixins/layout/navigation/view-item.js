import { MODALS } from '@/constants';

import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';

export default {
  mixins: [rightsTechnicalViewMixin],
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
  methods: {
    showEditViewModal(view) {
      this.$modals.show({
        name: MODALS.createView,
        config: { view },
      });
    },

    showDuplicateViewModal(view) {
      this.$modals.show({
        name: MODALS.createView,
        config: {
          view,
          isDuplicating: true,
        },
      });
    },
  },
};
