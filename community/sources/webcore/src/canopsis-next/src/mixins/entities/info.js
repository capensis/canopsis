import { createNamespacedHelpers } from 'vuex';
import { isMatch, isArray } from 'lodash';

import { DEFAULT_APP_TITLE } from '@/config';
import { CANOPSIS_EDITION, ROUTES_NAMES, USER_PERMISSIONS_TO_PAGES_RULES } from '@/constants';

import { sanitizeHtml } from '@/helpers/html';
import { compile } from '@/helpers/handlebars';
import { groupedPermissionToPermission } from '@/helpers/permission';

const { mapGetters, mapActions } = createNamespacedHelpers('info');

export const entitiesInfoMixin = {
  computed: {
    ...mapGetters({
      appInfoPending: 'pending',
      appInfo: 'appInfo',
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
      showHeaderOnKioskMode: 'showHeaderOnKioskMode',
      requiredInstructionApprove: 'requiredInstructionApprove',
      isBasicAuthEnabled: 'isBasicAuthEnabled',
      isLDAPAuthEnabled: 'isLDAPAuthEnabled',
      isCASAuthEnabled: 'isCASAuthEnabled',
      isSAMLAuthEnabled: 'isSAMLAuthEnabled',
      isOauthAuthEnabled: 'isOauthAuthEnabled',
      casConfig: 'casConfig',
      samlConfig: 'samlConfig',
      oauthConfig: 'oauthConfig',
      timezone: 'timezone',
      fileUploadMaxSize: 'fileUploadMaxSize',
      remediationJobConfigTypes: 'remediationJobConfigTypes',
      maintenance: 'maintenance',
      defaultColorTheme: 'defaultColorTheme',
      eventsCountTriggerDefaultThreshold: 'eventsCountTriggerDefaultThreshold',
      disabledTransitions: 'disabledTransitions',
      autoSuggestPbehaviorName: 'autoSuggestPbehaviorName',
    }),

    isProVersion() {
      return this.edition === CANOPSIS_EDITION.pro;
    },

    shownHeader() {
      return this.$route?.name === ROUTES_NAMES.viewKiosk
        ? this.showHeaderOnKioskMode
        : !this.$route?.meta?.hideHeader;
    },
  },
  methods: {
    ...mapActions(['fetchAppInfo', 'updateUserInterface']),

    checkAppInfoAccessByPermission(permission) {
      const preparedPermission = isArray(permission) ? groupedPermissionToPermission(permission) : permission;
      const permissionAppInfoRules = USER_PERMISSIONS_TO_PAGES_RULES[preparedPermission];

      if (!permissionAppInfoRules) {
        return true;
      }

      const appInfo = {
        edition: this.edition,
        stack: this.stack,
      };

      return isMatch(appInfo, permissionAppInfoRules);
    },

    async setTitle() {
      document.title = this.appTitle
        ? await compile(sanitizeHtml(this.appTitle, {
          allowedTags: [],
          allowedAttributes: {},
        }))
        : DEFAULT_APP_TITLE;
    },
  },
};
