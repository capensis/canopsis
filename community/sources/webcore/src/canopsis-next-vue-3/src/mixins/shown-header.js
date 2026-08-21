import { ROUTES_NAMES } from '@/constants';

import { activeViewMixin } from '@/mixins/active-view';

export const shownHeaderMixin = {
  mixins: [activeViewMixin],

  computed: {
    isKioskMode() {
      return this.activeViewIsKioskScreenMode || this.$route?.name === ROUTES_NAMES.viewKiosk;
    },

    shownHeader() {
      return this.isKioskMode ? this.showHeaderOnKioskMode : !this.$route?.meta?.hideHeader;
    },
  },
};

export default shownHeaderMixin;
