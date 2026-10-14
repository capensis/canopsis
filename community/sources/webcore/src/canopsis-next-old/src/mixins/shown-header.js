import { createNamespacedHelpers } from 'vuex';

const { mapGetters } = createNamespacedHelpers('info');
const { mapGetters: mapActiveViewGetters } = createNamespacedHelpers('activeView');

/**
 * Mixin to determine header visibility based on route and kiosk mode settings
 */
export const shownHeaderMixin = {
  computed: {
    ...mapGetters(['showHeaderOnKioskMode']),

    ...mapActiveViewGetters({
      activeViewIsKioskScreenMode: 'isKioskScreenMode',
    }),

    shownHeader() {
      return this.activeViewIsKioskScreenMode
        ? this.showHeaderOnKioskMode
        : !this.$route?.meta?.hideHeader;
    },
  },
};
