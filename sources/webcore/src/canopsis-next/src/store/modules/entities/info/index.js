import { API_ROUTES } from '@/config';
import request from '@/services/request';
import { toSeconds } from '@/helpers/date/duration';
import { POPUP_TYPES } from '@/constants';

const types = {
  FETCH_LOGIN_INFOS: 'FETCH_LOGIN_INFOS',
  FETCH_APP_INFOS: 'FETCH_APP_INFOS',
};

export default {
  namespaced: true,
  state: {
    version: '',
    logo: '',
    appTitle: '',
    footer: '',
    edition: '',
    stack: '',
    description: '',
    language: '',
    isLDAPAuthEnabled: false,
    isCASAuthEnabled: false,
    allowChangeSeverityToInfo: false,
    casConfig: {},
    popupTimeout: undefined,
    timezone: undefined,
    jobExecutorFetchTimeoutSeconds: undefined,
  },
  getters: {
    version: state => state.version,
    logo: state => state.logo,
    appTitle: state => state.appTitle,
    popupTimeout: state => state.popupTimeout,
    allowChangeSeverityToInfo: state => state.allowChangeSeverityToInfo,
    footer: state => state.footer,
    edition: state => state.edition,
    stack: state => state.stack,
    description: state => state.description,
    language: state => state.language,
    isLDAPAuthEnabled: state => state.isLDAPAuthEnabled,
    isCASAuthEnabled: state => state.isCASAuthEnabled,
    casConfig: state => state.casConfig,
    timezone: state => state.timezone,
    jobExecutorFetchTimeoutSeconds: state => state.jobExecutorFetchTimeoutSeconds,
  },
  mutations: {
    [types.FETCH_LOGIN_INFOS](state, {
      version,
      userInterface = {},
      loginConfig = {},
    }) {
      state.version = version;
      state.logo = userInterface.logo;
      state.appTitle = userInterface.app_title;
      state.footer = userInterface.footer;
      state.description = userInterface.login_page_description;
      state.language = userInterface.language;
      state.popupTimeout = userInterface.popup_timeout || {};

      state.isLDAPAuthEnabled = loginConfig.ldapconfig ? loginConfig.ldapconfig.enable : false;
      state.isCASAuthEnabled = loginConfig.casconfig ? loginConfig.casconfig.enable : false;

      state.casConfig = loginConfig.casconfig;
    },
    [types.FETCH_APP_INFOS](state, {
      version,
      logo,
      appTitle,
      popupTimeout,
      allowChangeSeverityToInfo,
      edition,
      stack,
      language,
      timezone,
      jobExecutorFetchTimeoutSeconds,
    }) {
      state.version = version;
      state.logo = logo;
      state.appTitle = appTitle;
      state.popupTimeout = popupTimeout || {};
      state.allowChangeSeverityToInfo = allowChangeSeverityToInfo;
      state.edition = edition;
      state.stack = stack;
      state.language = language;
      state.timezone = timezone;
      state.jobExecutorFetchTimeoutSeconds = jobExecutorFetchTimeoutSeconds;
    },
  },
  actions: {
    async fetchLoginInfos({ commit, dispatch }) {
      try {
        const {
          version,
          user_interface: userInterface,
          login_config: loginConfig,
        } = await request.get(API_ROUTES.infos.login);

        const { language, popup_timeout: popupTimeout } = userInterface;

        commit(types.FETCH_LOGIN_INFOS, {
          version,
          userInterface: userInterface || {},
          loginConfig: loginConfig || {},
        });

        if (language) {
          dispatch('i18n/setGlobalLocale', language, { root: true });
        }

        if (popupTimeout) {
          dispatch('setPopupTimeouts', { popupTimeout });
        }
      } catch (err) {
        console.error(err);
      }
    },

    async fetchAppInfos({ commit, dispatch }) {
      try {
        const {
          version,
          logo,
          app_title: appTitle,
          popup_timeout: popupTimeout,
          allow_change_severity_to_info: allowChangeSeverityToInfo,
          jobexecutorfetchtimeoutseconds: jobExecutorFetchTimeoutSeconds,
          edition,
          stack,
          language,
          timezone,
        } = await request.get(API_ROUTES.infos.app);

        commit(
          types.FETCH_APP_INFOS,
          {
            version,
            logo,
            appTitle,
            edition,
            popupTimeout,
            allowChangeSeverityToInfo,
            stack,
            language,
            timezone,
            jobExecutorFetchTimeoutSeconds,
          },
        );

        if (language) {
          dispatch('i18n/setGlobalLocale', language, { root: true });
        }

        if (popupTimeout) {
          dispatch('setPopupTimeouts', { popupTimeout });
        }
      } catch (err) {
        console.error(err);
      }
    },

    updateUserInterface(context, { data } = {}) {
      return request.post(API_ROUTES.infos.userInterface, data);
    },

    setPopupTimeouts({ dispatch }, { popupTimeout = {} }) {
      const { interval: intervalInfo, unit: unitInfo } = popupTimeout.info;
      const { interval: intervalError, unit: unitError } = popupTimeout.error;

      const timeInfo = toSeconds(intervalInfo, unitInfo) * 1000;
      const timeError = toSeconds(intervalError, unitError) * 1000;

      dispatch('popups/setDefaultCloseTime', { type: POPUP_TYPES.success, time: timeInfo }, { root: true });
      dispatch('popups/setDefaultCloseTime', { type: POPUP_TYPES.info, time: timeInfo }, { root: true });
      dispatch('popups/setDefaultCloseTime', { type: POPUP_TYPES.warning, time: timeError }, { root: true });
      dispatch('popups/setDefaultCloseTime', { type: POPUP_TYPES.error, time: timeError }, { root: true });
    },
  },
};
