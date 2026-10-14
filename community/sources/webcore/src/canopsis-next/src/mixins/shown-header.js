import { ROUTES_NAMES, VIEW_SCREEN_MODES } from '@/constants';

import { activeViewMixin } from '@/mixins/active-view';

export const shownHeaderMixin = {
  mixins: [activeViewMixin],

  computed: {
    isFullscreenMode() {
      return [VIEW_SCREEN_MODES.fullscreen, VIEW_SCREEN_MODES.kioskFullscreen].includes(
        this.activeViewScreenMode,
      );
    },

    isKioskMode() {
      return this.activeViewIsKioskScreenMode || this.$route?.name === ROUTES_NAMES.viewKiosk;
    },

    shownHeader() {
      if (this.isFullscreenMode) {
        return false;
      }

      return this.isKioskMode ? this.showHeaderOnKioskMode : !this.$route?.meta?.hideHeader;
    },
  },
};

export default shownHeaderMixin;
