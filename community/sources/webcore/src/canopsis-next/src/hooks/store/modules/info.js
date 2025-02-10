import { isMatch } from 'lodash';
import { computed } from 'vue';
import { useRoute } from 'vue-router/composables';

import { DEFAULT_APP_TITLE } from '@/config';
import { CANOPSIS_EDITION, ROUTES_NAMES, USER_PERMISSIONS_TO_PAGES_RULES } from '@/constants';

import { compile } from '@/helpers/handlebars';
import { sanitizeHtml } from '@/helpers/html';

import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates an instance of the info store module hooks.
 *
 * @returns {Object} Store module hooks for the 'info' namespace
 */
export const useInfoStoreModule = () => useStoreModuleHooks('info');

/**
 * Hook to determine header visibility based on route and kiosk mode settings.
 */
export const useShownHeader = () => {
  const { useGetters } = useInfoStoreModule();

  const getters = useGetters(['showHeaderOnKioskMode']);

  const route = useRoute();
  const shownHeader = computed(() => (
    route?.name === ROUTES_NAMES.viewKiosk
      ? getters.showHeaderOnKioskMode.value
      : !route?.meta?.hideHeader
  ));

  return {
    ...getters,
    shownHeader,
  };
};

/**
 * Hook for accessing and managing application information and settings.
 * Provides access to app configuration, authentication settings, UI preferences,
 * and various application-wide features.
 *
 * @typedef {Object} InfoGetters
 * @property {Ref<boolean>} appInfoPending - Loading state of app info
 * @property {Ref<Object>} appInfo - Raw application information
 * @property {Ref<string>} version - Application version
 * @property {Ref<string>} logo - Application logo URL
 * @property {Ref<string>} appTitle - Application title
 * @property {Ref<Object>} popupTimeout - Popup timeout settings
 * @property {Ref<number>} maxMatchedItems - Maximum number of matched items
 * @property {Ref<number>} checkCountRequestTimeout - Timeout for count requests
 * @property {Ref<Object>} footer - Footer configuration
 * @property {Ref<string>} edition - Application edition (pro/community)
 * @property {Ref<string>} description - Application description
 * @property {Ref<string>} language - Application language
 * @property {Ref<boolean>} allowChangeSeverityToInfo - Permission flag for severity changes
 * @property {Ref<boolean>} showHeaderOnKioskMode - Header visibility in kiosk mode
 * @property {Ref<boolean>} requiredInstructionApprove - Instruction approval requirement
 * @property {Ref<boolean>} isBasicAuthEnabled - Basic auth status
 * @property {Ref<boolean>} isLDAPAuthEnabled - LDAP auth status
 * @property {Ref<boolean>} isCASAuthEnabled - CAS auth status
 * @property {Ref<boolean>} isSAMLAuthEnabled - SAML auth status
 * @property {Ref<boolean>} isOauthAuthEnabled - OAuth auth status
 * @property {Ref<Object>} casConfig - CAS configuration
 * @property {Ref<Object>} samlConfig - SAML configuration
 * @property {Ref<Object>} oauthConfig - OAuth configuration
 * @property {Ref<string>} timezone - Application timezone
 * @property {Ref<number>} fileUploadMaxSize - Maximum file upload size
 * @property {Ref<Array>} remediationJobConfigTypes - Remediation job types
 * @property {Ref<Object>} maintenance - Maintenance settings
 * @property {Ref<string>} defaultColorTheme - Default UI theme
 * @property {Ref<number>} eventsCountTriggerDefaultThreshold - Default event count threshold
 * @property {Ref<Array>} disabledTransitions - Disabled state transitions
 * @property {Ref<boolean>} autoSuggestPbehaviorName - Auto-suggestion settings
 * @property {Ref<Array>} userTimezones - Available user timezones
 * @property {Ref<boolean>} shownUserTimezone - User timezone display preference
 * @property {Ref<string>} serialName - Build serial name
 * @property {Ref<number>} versionUpdated - Version updated timestamp
 * @property {Ref<string>} versionDescription - Version description template
 *
 * @returns {Object} Hook return object
 * @property {InfoGetters} getters - All available getters
 * @property {ComputedRef<boolean>} isProVersion - Whether running pro edition
 * @property {ComputedRef<boolean>} shownHeader - Header visibility state
 * @property {Function} fetchAppInfo - Fetches application information
 * @property {Function} updateUserInterface - Updates UI settings
 * @property {Function} checkAppInfoAccessByPermission - Checks permission access
 * @property {Function} setTitle - Sets document title
 */
export const useInfo = () => {
  const { useGetters, useActions } = useInfoStoreModule();

  const getters = useGetters({
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
    userTimezones: 'userTimezones',
    shownUserTimezone: 'shownUserTimezone',
    serialName: 'serialName',
    versionUpdated: 'versionUpdated',
    versionDescription: 'versionDescription',
  });

  const isProVersion = computed(() => getters.edition.value === CANOPSIS_EDITION.pro);

  const { shownHeader } = useShownHeader();

  const actions = useActions({
    fetchAppInfo: 'fetchAppInfo',
    updateUserInterface: 'updateUserInterface',
  });

  /**
   * Checks if a user has access to certain app info based on their permission.
   *
   * @param {string} permission - Permission to check
   * @returns {boolean} Whether access is allowed
   */
  const checkAppInfoAccessByPermission = (permission) => {
    const permissionAppInfoRules = USER_PERMISSIONS_TO_PAGES_RULES[permission];

    if (!permissionAppInfoRules) {
      return true;
    }

    const appInfo = { edition: getters.edition.value };

    return isMatch(appInfo, USER_PERMISSIONS_TO_PAGES_RULES[permission]);
  };

  /**
   * Sets the document title using the configured app title or default.
   * Sanitizes and compiles the title with Handlebars if templates are used.
   *
   * @returns {Promise<void>}
   */
  const setTitle = async () => {
    document.title = getters.appTitle.value
      ? await compile(sanitizeHtml(getters.appTitle.value, { allowedTags: [], allowedAttributes: {} }))
      : DEFAULT_APP_TITLE;
  };

  return {
    ...getters,
    isProVersion,
    shownHeader,

    ...actions,
    checkAppInfoAccessByPermission,
    setTitle,
  };
};
