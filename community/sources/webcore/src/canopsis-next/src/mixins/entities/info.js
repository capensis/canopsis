import { createNamespacedHelpers } from 'vuex';
import { isMatch } from 'lodash';

import { DEFAULT_APP_TITLE } from '@/config';
import { CANOPSIS_EDITION, USER_PERMISSIONS_TO_PAGES_RULES } from '@/constants';

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
      maxMatchedItems: 'maxMatchedItems',
      footer: 'footer',
      edition: 'edition',
      stack: 'stack',
      description: 'description',
      language: 'language',
      allowChangeSeverityToInfo: 'allowChangeSeverityToInfo',
      isLDAPAuthEnabled: 'isLDAPAuthEnabled',
      isCASAuthEnabled: 'isCASAuthEnabled',
      isSAMLAuthEnabled: 'isSAMLAuthEnabled',
      casConfig: 'casConfig',
      samlConfig: 'samlConfig',
      timezone: 'timezone',
      jobExecutorFetchTimeoutSeconds: 'jobExecutorFetchTimeoutSeconds',
    }),

    isCatVersion() {
      return this.edition === CANOPSIS_EDITION.cat;
    },
  },
  methods: {
    ...mapActions(['fetchLoginInfos', 'fetchAppInfos', 'updateUserInterface']),

    checkAppInfoAccessByRight(right) {
      const rightAppInfoRules = USER_PERMISSIONS_TO_PAGES_RULES[right];

      if (!rightAppInfoRules) {
        return true;
      }

      const appInfo = {
        edition: this.edition,
        stack: this.stack,
      };

      return isMatch(appInfo, USER_PERMISSIONS_TO_PAGES_RULES[right]);
    },

    setTitle() {
      document.title = this.appTitle || DEFAULT_APP_TITLE;
    },
  },
};
