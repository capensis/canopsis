import { createNamespacedHelpers } from 'vuex';
import { isMatch } from 'lodash';

import { DEFAULT_APP_TITLE } from '@/config';
import { CANOPSIS_EDITION, USER_PERMISSIONS_TO_PAGES_RULES } from '@/constants';

const { mapGetters, mapActions } = createNamespacedHelpers('info');

/**
 * @mixin
 */
export const entitiesInfoMixin = {
  computed: {
    ...mapGetters({
      version: 'version',
      logo: 'logo',
      appTitle: 'appTitle',
      popupTimeout: 'popupTimeout',
      maxMatchedItems: 'maxMatchedItems',
      checkCountRequestTimeout: 'checkCountRequestTimeout',
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
      fileUploadMaxSize: 'fileUploadMaxSize',
      remediationJobConfigTypes: 'remediationJobConfigTypes',
      remediationPauseManualInstructionIntervalSeconds: 'remediationPauseManualInstructionIntervalSeconds',
    }),

    isCatVersion() {
      return this.edition === CANOPSIS_EDITION.cat;
    },
  },
  methods: {
    ...mapActions(['fetchLoginInfos', 'fetchAppInfos', 'updateUserInterface']),

    checkAppInfoAccessByPermission(permission) {
      const permissionAppInfoRules = USER_PERMISSIONS_TO_PAGES_RULES[permission];

      if (!permissionAppInfoRules) {
        return true;
      }

      const appInfo = {
        edition: this.edition,
        stack: this.stack,
      };

      return isMatch(appInfo, USER_PERMISSIONS_TO_PAGES_RULES[permission]);
    },

    setTitle() {
      document.title = this.appTitle || DEFAULT_APP_TITLE;
    },
  },
};
