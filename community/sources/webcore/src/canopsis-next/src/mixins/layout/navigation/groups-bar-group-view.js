import { MODALS } from '@/constants';

import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';

import layoutNavigationEditingModeMixin from './editing-mode';

export default {
  mixins: [
    rightsTechnicalViewMixin,
    layoutNavigationEditingModeMixin,
  ],
  props: {
    view: {
      type: Object,
      required: true,
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
      return (this.hasUpdateViewAccess || this.hasDeleteViewAccess) && this.isNavigationEditingMode;
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
