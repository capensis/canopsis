import { get } from 'lodash';

import { API_ROUTES, INSTRUCTION_EXECUTE_FETCHING_INTERVAL_SECONDS } from '@/config';
import { POPUP_TYPES } from '@/constants';

import request from '@/services/request';

const types = {
  FETCH_LOGIN_INFOS: 'FETCH_LOGIN_INFOS',
  FETCH_APP_INFOS: 'FETCH_APP_INFOS',
};

export default {
  namespaced: true,
  state: {
    version: '',
    userInterface: {},
    loginConfig: {},
  },
  getters: {
    version: state => state.version,
    userInterface: state => state.userInterface,
    logo: state => state.userInterface.logo,
    appTitle: state => state.userInterface.app_title,
    popupTimeout: state => state.userInterface.popup_timeout || {},
    maxMatchedItems: state => state.userInterface.max_matched_items,
    checkCountRequestTimeout: state => state.userInterface.check_count_request_timeout,
    allowChangeSeverityToInfo: state => state.userInterface.allow_change_severity_to_info,
    footer: state => state.userInterface.footer,
    edition: state => state.userInterface.edition,
    stack: state => state.userInterface.stack,
    description: state => state.userInterface.login_page_description,
    language: state => state.userInterface.language,
    timezone: state => state.userInterface.timezone,
    remediation: state => state.userInterface.remediation,
    remediationJobConfigTypes: state => get(state.userInterface.remediation, 'job_config_types', []),
    remediationPauseManualInstructionIntervalSeconds: state => get(
      state.userInterface.remediation,
      'pause_manual_instruction_interval.seconds',
      INSTRUCTION_EXECUTE_FETCHING_INTERVAL_SECONDS,
    ),

    isLDAPAuthEnabled: state => !!state.loginConfig.ldapconfig?.enable,
    isCASAuthEnabled: state => !!state.loginConfig.casconfig?.enable,
    isSAMLAuthEnabled: state => !!state.loginConfig.saml2config?.enable,
    casConfig: state => state.loginConfig.casconfig,
    samlConfig: state => state.loginConfig.saml2config,
  },
  mutations: {
    [types.FETCH_LOGIN_INFOS](state, { loginInfo = {} }) {
      const {
        version,
        user_interface: userInterface = {},
        login_config: loginConfig = {},
      } = loginInfo;

      state.userInterface = userInterface;
      state.loginConfig = loginConfig;
      state.version = version;
    },
    [types.FETCH_APP_INFOS](state, { userInterface }) {
      state.userInterface = userInterface;
    },
  },
  actions: {
    async fetchLoginInfos({ commit, dispatch }) {
      const loginInfo = await request.get(API_ROUTES.infos.login, { fullResponse: true });

      const { language, popup_timeout: popupTimeout } = loginInfo.user_interface;

      commit(types.FETCH_LOGIN_INFOS, { loginInfo });

      if (language) {
        dispatch('i18n/setGlobalLocale', language, { root: true });
      }

      if (popupTimeout) {
        dispatch('setPopupTimeouts', { popupTimeout });
      }
    },

    async fetchAppInfos({ commit, dispatch }) {
      try {
        const userInterface = await request.get(API_ROUTES.infos.app);

        commit(types.FETCH_APP_INFOS, { userInterface });

        if (userInterface.language) {
          dispatch('i18n/setGlobalLocale', userInterface.language, { root: true });
        }

        if (userInterface.popup_timeout) {
          dispatch('setPopupTimeouts', { popupTimeout: userInterface.popup_timeout });
        }
      } catch (err) {
        console.error(err);
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
