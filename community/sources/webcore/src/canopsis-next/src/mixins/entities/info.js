import { CANOPSIS_EDITION, USER_RIGHTS_TO_PAGES_RULES } from '@/constants';

import { createNamespacedHelpers } from 'vuex';
import { isMatch } from 'lodash';
import { DEFAULT_APP_TITLE } from '@/config';

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
      allowChangeSeverityToInfo: 'allowChangeSeverityToInfo',
      isCASAuthEnabled: 'isCASAuthEnabled',
      casConfig: 'casConfig',
      timezone: 'timezone',
      jobExecutorFetchTimeoutSeconds: 'jobExecutorFetchTimeoutSeconds',
    }),

    isCatVersion() {
      return this.edition === CANOPSIS_EDITION.cat;
    },
  },
  methods: {
    ...mapActions({
      fetchLoginInfos: 'fetchLoginInfos',
      fetchAppInfos: 'fetchAppInfos',
      updateUserInterface: 'updateUserInterface',
    }),

    checkAppInfoAccessByRight(right) {
      const rightAppInfoRules = USER_RIGHTS_TO_PAGES_RULES[right];

      if (!rightAppInfoRules) {
        return true;
      }

      const appInfo = {
        edition: this.edition,
        stack: this.stack,
      };

      return isMatch(appInfo, USER_RIGHTS_TO_PAGES_RULES[right]);
    },

    setTitle() {
      document.title = this.appTitle || DEFAULT_APP_TITLE;
    },
  },
};
