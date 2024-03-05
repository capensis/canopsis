import { createNamespacedHelpers } from 'vuex';

import { CRUD_ACTIONS, GROUPS_NAVIGATION_TYPES } from '@/constants';

import { checkUserAccess } from '@/helpers/entities/permissions/list';

const { mapGetters, mapActions } = createNamespacedHelpers('auth');

export const authMixin = {
  computed: {
    ...mapGetters(['isLoggedIn', 'currentUser', 'currentUserPermissionsById']),
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

    /**
     * Show groups top-bar only for ui_groups_navigation_type='top-bar' only for laptop
     *
     * @returns {boolean|*}
     */
    isShownGroupsTopBar() {
      const { ui_groups_navigation_type: groupsNavigationType } = this.currentUser;
      const isSelectedTopBar = groupsNavigationType === GROUPS_NAVIGATION_TYPES.topBar;
      const isLaptop = this.$options.filters.mq(this.$mq, { l: true });

      return isSelectedTopBar && isLaptop;
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
