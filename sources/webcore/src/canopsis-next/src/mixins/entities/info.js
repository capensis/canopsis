import { POPUP_TYPES, USER_RIGHTS_TO_EXPLOITATION_PAGES_RULES } from '@/constants';

import { createNamespacedHelpers } from 'vuex';
import { isMatch } from 'lodash';
import { getSecondsByUnit } from '@/helpers/time';

const { mapGetters, mapActions } = createNamespacedHelpers('info');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      version: 'version',
      logo: 'logo',
      appTitle: 'appTitle',
      popupTimeout: 'popupTimeout',
      footer: 'footer',
      edition: 'edition',
      stack: 'stack',
      description: 'description',
      language: 'language',
      isLDAPAuthEnabled: 'isLDAPAuthEnabled',
      isCASAuthEnabled: 'isCASAuthEnabled',
      casConfig: 'casConfig',
    }),
  },
  methods: {
    ...mapActions({
      fetchLoginInfos: 'fetchLoginInfos',
      fetchAppInfos: 'fetchAppInfos',
      updateUserInterface: 'updateUserInterface',
    }),

    checkAppInfoAccessByRight(right) {
      const rightAppInfoRules = USER_RIGHTS_TO_EXPLOITATION_PAGES_RULES[right];

      if (!rightAppInfoRules) {
        return true;
      }

      const appInfo = {
        edition: this.edition,
        stack: this.stack,
      };

      return isMatch(appInfo, USER_RIGHTS_TO_EXPLOITATION_PAGES_RULES[right]);
    },

    setErrorPopupTime() {
      const { interval, unit } = this.popupTimeout.error;
      const delay = getSecondsByUnit(interval, unit) * 1000;

      this.$popups.setDefaultCloseTime(POPUP_TYPES.error, delay);
      this.$popups.setDefaultCloseTime(POPUP_TYPES.warning, delay);
    },

    setInfoPopupTime() {
      const { interval, unit } = this.popupTimeout.info;
      const delay = getSecondsByUnit(interval, unit) * 1000;

      this.$popups.setDefaultCloseTime(POPUP_TYPES.infos, delay);
      this.$popups.setDefaultCloseTime(POPUP_TYPES.success, delay);
    },

    setPopupTimeout() {
      if (!this.popupTimeout) {
        return;
      }

      if (this.popupTimeout.error) {
        this.setErrorPopupTime();
      }

      if (this.popupTimeout.info) {
        this.setInfoPopupTime();
      }
    },
  },
};
