import { createNamespacedHelpers } from 'vuex';

import { USERS_RIGHTS_MASKS, GROUPS_NAVIGATION_TYPES } from '@/constants';
import { checkUserAccess } from '@/helpers/right';

const { mapGetters, mapActions } = createNamespacedHelpers('auth');

export default {
  computed: {
    ...mapGetters(['isLoggedIn', 'currentUser']),
    ...mapGetters({
      currentUserPending: 'pending',
    }),

    checkAccess() {
      return (rightId, rightMask = 1) => checkUserAccess(this.currentUser, rightId, rightMask);
    },

    checkCreateAccess() {
      return rightId => this.checkAccess(rightId, USERS_RIGHTS_MASKS.create);
    },

    checkReadAccess() {
      return rightId => this.checkAccess(rightId, USERS_RIGHTS_MASKS.read);
    },

    checkUpdateAccess() {
      return rightId => this.checkAccess(rightId, USERS_RIGHTS_MASKS.update);
    },

    checkDeleteAccess() {
      return rightId => this.checkAccess(rightId, USERS_RIGHTS_MASKS.delete);
    },

    /**
     * Show groups side-bar only for groupsNavigationType='side-bar' or for mobile and tablet
     *
     * @returns {boolean|*}
     */
    isShownGroupsSideBar() {
      const isSelectedSideBar = this.currentUser.groupsNavigationType === GROUPS_NAVIGATION_TYPES.sideBar;
      const isMobileOrTablet = this.$options.filters.mq(this.$mq, { m: true, l: false });

      return isSelectedSideBar || isMobileOrTablet;
    },

    /**
     * Show groups top-bar only for groupsNavigationType='top-bar' only for laptop
     *
     * @returns {boolean|*}
     */
    isShownGroupsTopBar() {
      const isSelectedTopBar = this.currentUser.groupsNavigationType === GROUPS_NAVIGATION_TYPES.topBar;
      const isLaptop = this.$options.filters.mq(this.$mq, { l: true });

      return isSelectedTopBar && isLaptop;
    },
  },
  methods: {
    ...mapActions(['login', 'logout', 'fetchCurrentUser']),
  },
};
