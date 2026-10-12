import { DEFAULT_APP_TITLE } from '@/config';

import { compile } from '@/helpers/handlebars';
import { sanitizeHtml } from '@/helpers/html';

import { useInfo } from '@/hooks/store/modules/info';

const infoComputed = [
  'appInfoPending',
  'appInfo',
  'version',
  'logo',
  'appTitle',
  'popupTimeout',
  'maxMatchedItems',
  'checkCountRequestTimeout',
  'footer',
  'edition',
  'description',
  'language',
  'allowChangeSeverityToInfo',
  'showHeaderOnKioskMode',
  'requiredInstructionApprove',
  'isBasicAuthEnabled',
  'isLDAPAuthEnabled',
  'isCASAuthEnabled',
  'isSAMLAuthEnabled',
  'isOauthAuthEnabled',
  'casConfig',
  'samlConfig',
  'oauthConfig',
  'timezone',
  'dateFormat',
  'fileUploadMaxSize',
  'remediationJobConfigTypes',
  'maintenance',
  'defaultColorTheme',
  'eventsCountTriggerDefaultThreshold',
  'disabledTransitions',
  'autoSuggestPbehaviorName',
  'userTimezones',
  'shownUserTimezone',
  'isProVersion',
].reduce((acc, key) => {
  acc[key] = () => useInfo()[key].value;

  return acc;
}, {});

export const entitiesInfoMixin = {
  computed: infoComputed,

  methods: {
    fetchAppInfo(...args) {
      return useInfo().fetchAppInfo(...args);
    },

    updateUserInterface(...args) {
      return useInfo().updateUserInterface(...args);
    },

    updateMaintenanceMode(...args) {
      return useInfo().updateMaintenanceMode(...args);
    },

    checkAppInfoAccessByPermission(...args) {
      return useInfo().checkAppInfoAccessByPermission(...args);
    },

    async setTitle() {
      document.title = this.appTitle
        ? await compile(
            sanitizeHtml(this.appTitle, {
              ALLOWED_TAGS: [],
              ALLOWED_ATTR: [],
            }),
          )
        : DEFAULT_APP_TITLE;
    },
  },
};

export default entitiesInfoMixin;
