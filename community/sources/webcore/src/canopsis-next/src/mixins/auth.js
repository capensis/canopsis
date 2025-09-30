import { createNamespacedHelpers } from 'vuex';

import { CRUD_ACTIONS, GROUPS_NAVIGATION_TYPES, ROUTES_NAMES, VIEW_USER_PERMISSIONS_NAMES } from '@/constants';

import { checkUserAccess } from '@/helpers/entities/permissions/list';

const { mapGetters, mapActions } = createNamespacedHelpers('auth');

export const authMixin = {
  computed: {
    ...mapGetters(['isLoggedIn', 'currentUser', 'currentUserPermissionsById', 'currentUserViewPermissionsByViewId']),
    ...mapGetters({
      currentUserPending: 'pending',
    }),

    /**
     * Show groups side-bar only for ui_groups_navigation_type='side-bar' or for mobile and tablet
     *
     * @returns {boolean|*}
     */
    isShownGroupsSideBar() {
      const { ui_groups_navigation_type: groupsNavigationType } = this.currentUser;
      const isSelectedSideBar = groupsNavigationType === GROUPS_NAVIGATION_TYPES.sideBar;
      const isMobileOrTablet = this.$options.filters.mq(this.$mq, { m: true, l: false });

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
    ...mapActions([
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
