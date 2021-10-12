import { API_ROUTES, INSTRUCTION_EXECUTE_FETCHING_INTERVAL_SECONDS } from '@/config';
import { POPUP_TYPES } from '@/constants';

import request from '@/services/request';

const types = {
  FETCH_APP_INFO: 'FETCH_APP_INFO',
};

export default {
  namespaced: true,
  state: {
    appInfo: {},
  },
  getters: {
    appInfo: state => state.appInfo,
    version: state => state.appInfo.version,
    logo: state => state.appInfo.logo,
    appTitle: state => state.appInfo.app_title,
    popupTimeout: state => state.appInfo.popup_timeout || {},
    maxMatchedItems: state => state.appInfo.max_matched_items,
    checkCountRequestTimeout: state => state.appInfo.check_count_request_timeout,
    allowChangeSeverityToInfo: state => state.appInfo.allow_change_severity_to_info,
    footer: state => state.appInfo.footer,
    edition: state => state.appInfo.edition,
    stack: state => state.appInfo.stack,
    description: state => state.appInfo.login_page_description,
    language: state => state.appInfo.language,
    timezone: state => state.appInfo.timezone,
    fileUploadMaxSize: state => state.appInfo.file_upload_max_size ?? 0,
    remediationJobConfigTypes: state => state.appInfo.remediation?.job_config_types ?? [],
    remediationPauseManualInstructionIntervalSeconds:
        state => state.appInfo.remediation?.pause_manual_instruction_interval.seconds
          ?? INSTRUCTION_EXECUTE_FETCHING_INTERVAL_SECONDS,

    casConfig: state => state.appInfo?.login?.casconfig,
    samlConfig: state => state.appInfo?.login?.saml2config,
    isLDAPAuthEnabled: state => !!state.appInfo?.login?.ldapconfig?.enable,
    isCASAuthEnabled: state => !!state.appInfo?.login?.casconfig?.enable,
    isSAMLAuthEnabled: state => !!state.appInfo?.login?.saml2config?.enable,
  },
  mutations: {
    [types.FETCH_APP_INFO](state, { appInfo }) {
      state.appInfo = appInfo;
    },
  },
  actions: {
    async fetchAppInfo({ commit, dispatch }) {
      const appInfo = await request.get(API_ROUTES.infos.app);

      commit(types.FETCH_APP_INFO, { appInfo });

      if (appInfo.language) {
        dispatch('i18n/setGlobalLocale', appInfo.language, { root: true });
      }

      if (appInfo.popup_timeout) {
        dispatch('setPopupTimeouts', { popupTimeout: appInfo.popup_timeout });
      }
    },

    updateUserInterface(context, { data } = {}) {
      return request.post(API_ROUTES.infos.userInterface, data);
    },

    setPopupTimeouts({ dispatch }, { popupTimeout = {} }) {
      const { seconds: intervalInfo } = popupTimeout.info;
      const { seconds: intervalError } = popupTimeout.error;

      const timeInfo = intervalInfo * 1000;
      const timeError = intervalError * 1000;

      dispatch('popups/setDefaultCloseTime', { type: POPUP_TYPES.success, time: timeInfo }, { root: true });
      dispatch('popups/setDefaultCloseTime', { type: POPUP_TYPES.info, time: timeInfo }, { root: true });
      dispatch('popups/setDefaultCloseTime', { type: POPUP_TYPES.warning, time: timeError }, { root: true });
      dispatch('popups/setDefaultCloseTime', { type: POPUP_TYPES.error, time: timeError }, { root: true });
    },
  },
};
