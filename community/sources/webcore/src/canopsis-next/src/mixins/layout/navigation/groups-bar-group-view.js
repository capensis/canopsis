import { mapActions, mapState } from 'pinia';

import { MODALS, ROUTES_NAMES } from '@/constants';

import { useActiveViewStore } from '@/stores/modules/active-view';

import { entitiesViewMixin } from '@/mixins/entities/view';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';
import { viewRouterMixin } from '@/mixins/view/router';

import { layoutNavigationEditingModeMixin } from './editing-mode';

export default {
  mixins: [
    permissionsTechnicalViewMixin,
    layoutNavigationEditingModeMixin,
    entitiesViewMixin,
    entitiesViewGroupMixin,
    viewRouterMixin,
  ],
  props: {
    view: {
      type: Object,
      required: true,
    },
  },
  computed: {
    ...mapState(useActiveViewStore, {
      activeViewId: 'itemId',
    }),

    viewLink() {
      const { view } = this;
      const link = {
        name: ROUTES_NAMES.view,
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
      return this.checkDeleteAccess(this.view._id) && this.hasDeleteAnyViewAccess;
    },

    hasViewEditButtonAccess() {
      return this.hasUpdateViewAccess || this.hasDeleteViewAccess;
    },
  },
  methods: {
    ...mapActions(useActiveViewStore, {
      fetchActiveView: 'fetch',
    }),

    showEditViewModal() {
      this.$modals.show({
        name: MODALS.createView,
        config: {
          title: this.$t('modals.view.edit.title'),
          view: this.view,
          deletable: this.view.is_private || this.hasDeleteViewAccess,
          submittable: this.view.is_private || this.hasUpdateViewAccess,
          action: async (data) => {
            await this.updateViewWithPopup({ id: this.view._id, data });

            if (this.activeViewId === this.view._id) {
              await this.fetchActiveView();
            }

            return this.fetchAllGroupsListWithWidgetsWithCurrentUser();
          },
          remove: async () => {
            await this.removeViewWithPopup({ id: this.view._id });
            await this.fetchAllGroupsListWithWidgetsWithCurrentUser();

            return this.redirectToHomeIfCurrentRoute();
          },
        },
      });
    },

    showDuplicateViewModal() {
      this.$modals.show({
        name: MODALS.createView,
        config: {
          title: this.$t('modals.view.duplicate.title', { viewTitle: this.view.title }),
          duplicate: true,
          duplicableToAll: this.hasCreateAnyViewAccess,
          view: {
            ...this.view,

            name: '',
            title: '',
          },
          submittable: this.view.is_private || this.hasCreateAnyViewAccess,
          action: async (data) => {
            await this.copyViewWithPopup({ id: this.view._id, data });

            return this.fetchAllGroupsListWithWidgetsWithCurrentUser();
          },
        },
      });
    },
  },
};
