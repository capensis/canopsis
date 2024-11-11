import { isMatch } from 'lodash';
import { computed } from 'vue';
import { useRoute } from 'vue-router/composables';

import { DEFAULT_APP_TITLE } from '@/config';
import { CANOPSIS_EDITION, ROUTES_NAMES, USER_PERMISSIONS_TO_PAGES_RULES } from '@/constants';

import { compile } from '@/helpers/handlebars';
import { sanitizeHtml } from '@/helpers/html';

import { useStoreModuleHooks } from '@/hooks/store';

export const useInfoStoreModule = () => useStoreModuleHooks('info');

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
    userTimezones: 'userTimezones',
  });

  const isProVersion = computed(() => getters.edition.value === CANOPSIS_EDITION.pro);

  const { shownHeader } = useShownHeader();

  const actions = useActions({
    fetchAppInfo: 'fetchAppInfo',
    updateUserInterface: 'updateUserInterface',
  });

  const checkAppInfoAccessByPermission = (permission) => {
    const permissionAppInfoRules = USER_PERMISSIONS_TO_PAGES_RULES[permission];

    if (!permissionAppInfoRules) {
      return true;
    }

    const appInfo = {
      edition: getters.edition.value,
      stack: getters.stack.value,
    };

    return isMatch(appInfo, USER_PERMISSIONS_TO_PAGES_RULES[permission]);
  };

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
