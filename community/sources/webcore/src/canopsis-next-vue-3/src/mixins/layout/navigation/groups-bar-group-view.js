import { CRUD_ACTIONS, ROUTES_NAMES, VIEW_USER_PERMISSIONS_NAMES } from '@/constants';

import { checkUserAccess } from '@/helpers/permission';

import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';
import { entitiesViewMixin } from '@/mixins/entities/view';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
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
      return this.checkViewAccess(CRUD_ACTIONS.update) && this.hasUpdateAnyViewAccess;
    },

    hasDeleteViewAccess() {
      return this.checkViewAccess(CRUD_ACTIONS.delete) && this.hasUpdateAnyViewAccess;
    },

    hasViewEditButtonAccess() {
      return this.hasUpdateViewAccess || this.hasDeleteViewAccess;
    },
  },
  methods: {
    checkViewAccess(action) {
      return checkUserAccess(
        this.currentUserViewPermissionsByViewId[this.view._id]?.[VIEW_USER_PERMISSIONS_NAMES.actions],
        action,
      );
    },

    showEditViewModal() {
      // TODO: Enable after view modals are migrated.
    },

    showDuplicateViewModal() {
      // TODO: Enable after view modals are migrated.
    },
  },
};
