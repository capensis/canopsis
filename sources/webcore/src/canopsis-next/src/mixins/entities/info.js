import { USER_RIGHTS_TO_EXPLOITATION_PAGES_RULES } from '@/constants';

import { createNamespacedHelpers } from 'vuex';
import { isMatch } from 'lodash';

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
  },
};
