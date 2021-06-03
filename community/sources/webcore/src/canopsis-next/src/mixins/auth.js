import { get } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { CRUD_ACTIONS, GROUPS_NAVIGATION_TYPES } from '@/constants';
import { checkUserAccess } from '@/helpers/permission';

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
    ...mapActions(['login', 'logout', 'fetchCurrentUser']),

    checkIsTourEnabled(tour) {
      return !get(this.currentUser, ['ui_tours', tour]);
    },

    checkAccess(rightId, action = CRUD_ACTIONS.can) {
      return checkUserAccess(this.currentUserPermissionsById, rightId, action);
    },

    checkCreateAccess(rightId) {
      return this.checkAccess(rightId, CRUD_ACTIONS.create);
    },

    checkReadAccess(rightId) {
      return this.checkAccess(rightId, CRUD_ACTIONS.read);
    },

    checkUpdateAccess(rightId) {
      return this.checkAccess(rightId, CRUD_ACTIONS.update);
    },

    checkDeleteAccess(rightId) {
      return this.checkAccess(rightId, CRUD_ACTIONS.delete);
    },
  },
};

export default authMixin;