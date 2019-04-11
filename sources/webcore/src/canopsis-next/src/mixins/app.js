import { createNamespacedHelpers } from 'vuex';

import { GROUPS_NAVIGATION_TYPES } from '@/constants';

const { mapGetters, mapActions } = createNamespacedHelpers('app');

export default {
  computed: {
    ...mapGetters(['groupsNavigationType']),

    /**
     * Show groups side-bar only for groupsNavigationType='side-bar' or for mobile and tablet
     *
     * @returns {boolean|*}
     */
    isShownGroupsSideBar() {
      const isSelectedSideBar = this.groupsNavigationType === GROUPS_NAVIGATION_TYPES.sideBar;
      const isMobileOrTablet = this.$options.filters.mq(this.$mq, { m: true, l: false });

      return isSelectedSideBar || isMobileOrTablet;
    },

    /**
     * Show groups top-bar only for groupsNavigationType='top-bar' only for laptop
     *
     * @returns {boolean|*}
     */
    isShownGroupsTopBar() {
      const isSelectedTopBar = this.groupsNavigationType === GROUPS_NAVIGATION_TYPES.topBar;
      const isLaptop = this.$options.filters.mq(this.$mq, { l: true });

      return isSelectedTopBar && isLaptop;
    },
  },
  methods: {
    ...mapActions(['setGroupsNavigationType']),
  },
};
