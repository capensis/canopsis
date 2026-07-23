export const shownHeaderMixin = {
  computed: {
    activeViewIsKioskScreenMode() {
      return this.$route?.name === 'view-kiosk';
    },

    shownHeader() {
      return this.activeViewIsKioskScreenMode
        ? this.showHeaderOnKioskMode
        : !this.$route?.meta?.hideHeader;
    },
  },
};

export default shownHeaderMixin;
