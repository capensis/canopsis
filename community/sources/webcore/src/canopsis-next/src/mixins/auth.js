import { mapActions, mapState, mapStores } from 'pinia';

import { MEDIA_QUERIES_BREAKPOINTS } from '@/config';
import {
  GROUPS_NAVIGATION_TYPES,
  CRUD_ACTIONS,
  ROUTES_NAMES,
  VIEW_USER_PERMISSIONS_NAMES,
} from '@/constants';

import { useAuthStore } from '@/stores/modules/auth';

import { checkUserAccess } from '@/helpers/permission';

export const authMixin = {
  computed: {
    ...mapStores(useAuthStore),
    ...mapState(useAuthStore, {
      isLoggedIn: 'isLoggedIn',
      currentUser: 'currentUser',
      currentUserPending: 'pending',
      currentUserPermissionsById: 'currentUserPermissionsById',
      currentUserViewPermissionsByViewId: 'currentUserViewPermissionsByViewId',
    }),

    isShownGroupsTopBar() {
      return this.currentUser?.ui_groups_navigation_type === GROUPS_NAVIGATION_TYPES.topBar;
    },

    isShownGroupsSideBar() {
      const groupsNavigationType = this.currentUser?.ui_groups_navigation_type;
      const isSelectedSideBar = groupsNavigationType === GROUPS_NAVIGATION_TYPES.sideBar;
      const isMobileOrTablet = this.$vuetify.display.width < MEDIA_QUERIES_BREAKPOINTS.t;

      return isSelectedSideBar || isMobileOrTablet || !this.isShownGroupsTopBar;
    },

    hasCurrentViewActionsAccess() {
      const { name, params } = this.$route;

      if (![ROUTES_NAMES.view, ROUTES_NAMES.viewKiosk].includes(name)) {
        return false;
      }

      return checkUserAccess(
        this.currentUserViewPermissionsByViewId[params.id]?.[VIEW_USER_PERMISSIONS_NAMES.actions],
        CRUD_ACTIONS.can,
      );
    },
  },

  methods: {
    ...mapActions(useAuthStore, [
      'login',
      'applyAccessToken',
      'logout',
      'fetchCurrentUser',
      'filesAccess',
    ]),

    checkAccess(permissionId, action = CRUD_ACTIONS.can) {
      return checkUserAccess(this.currentUserPermissionsById[permissionId], action);
    },

    checkCreateAccess(permissionId) {
      return this.checkAccess(permissionId, CRUD_ACTIONS.create);
    },

    checkReadAccess(permissionId) {
      return this.checkAccess(permissionId, CRUD_ACTIONS.read);
    },

    checkUpdateAccess(permissionId) {
      return this.checkAccess(permissionId, CRUD_ACTIONS.update);
    },

    checkDeleteAccess(permissionId) {
      return this.checkAccess(permissionId, CRUD_ACTIONS.delete);
    },
  },
};

export default authMixin;
