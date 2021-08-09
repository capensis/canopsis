import { omit } from 'lodash';

import { MODALS, ROUTE_NAMES } from '@/constants';

import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';

import layoutNavigationEditingModeMixin from './editing-mode';

export default {
  mixins: [
    permissionsTechnicalViewMixin,
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
        name: ROUTE_NAMES.view,
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
          title: this.$t('modals.view.edit.title'),
          view: this.view,
        },
      });
    },

    showDuplicateViewModal() {
      this.$modals.show({
        name: MODALS.createView,
        config: {
          title: this.$t('modals.view.duplicate.title', { viewTitle: this.view.title }),
          view: {
            ...omit(this.view, ['_id', 'name', 'title']),

            name: '',
            title: '',
          },
        },
      });
    },
  },
};
