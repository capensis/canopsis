import { MODALS } from '@/constants';

import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';

export default {
  mixins: [rightsTechnicalViewMixin],
  props: {
    view: {
      type: Object,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    viewLink() {
      const { view } = this;
      const link = {
        name: 'view',
        params: { id: view._id },
      };

      if (view.tabs && view.tabs.length) {
        link.query = { tabId: view.tabs[0]._id };
      }

      return link;
    },

    hasUpdateViewAccess() {
      return this.checkUpdateAccess(this.view._id) && this.hasUpdateAnyViewAccess;
    },

    hasDeleteViewAccess() {
      return this.checkDeleteAccess(this.view._id) && this.hasUpdateAnyViewAccess;
    },

    hasViewEditButtonAccess() {
      return (this.hasUpdateViewAccess || this.hasDeleteViewAccess) && this.isEditingMode;
    },
  },
  methods: {
    showEditViewModal() {
      this.$modals.show({
        name: MODALS.createView,
        config: {
          view: this.view,
        },
      });
    },

    showDuplicateViewModal() {
      this.$modals.show({
        name: MODALS.createView,
        config: {
          view: this.view,
          isDuplicating: true,
        },
      });
    },
  },
};
