import { API_ROUTES } from '@/config';
import { POPUP_TYPES } from '@/constants';

import request from '@/services/request';
import localStorage from '@/services/local-storage';

import { durationToSeconds } from '@/helpers/date/duration';

const types = {
  FETCH_APP_INFO: 'FETCH_APP_INFO',
  FETCH_APP_INFO_COMPLETED: 'FETCH_APP_INFO_COMPLETED',
  FETCH_APP_INFO_FAILED: 'FETCH_APP_INFO_FAILED',
};

export default {
  namespaced: true,
  state: {
    appInfo: {},
    pending: false,
  },
  getters: {
    pending: state => state.pending,
    appInfo: state => state.appInfo,
    version: state => state.appInfo.version,
    logo: state => state.appInfo.logo,
    appTitle: state => state.appInfo.app_title,
    maintenance: state => state.appInfo.maintenance,
    defaultColorTheme: state => state.appInfo.default_color_theme,
    popupTimeout: state => state.appInfo.popup_timeout || {},
    maxMatchedItems: state => state.appInfo.max_matched_items,
    checkCountRequestTimeout: state => state.appInfo.check_count_request_timeout,
    allowChangeSeverityToInfo: state => state.appInfo.allow_change_severity_to_info,
    showHeaderOnKioskMode: state => state.appInfo.show_header_on_kiosk_mode,
    requiredInstructionApprove: state => state.appInfo.required_instruction_approve,
    footer: state => state.appInfo.footer,
    edition: state => state.appInfo.edition,
    stack: state => state.appInfo.stack,
    description: state => state.appInfo.login_page_description,
    language: state => state.appInfo.language,
    timezone: state => state.appInfo.timezone,
    fileUploadMaxSize: state => state.appInfo.file_upload_max_size ?? 0,
    remediationJobConfigTypes: state => state.appInfo.remediation?.job_config_types ?? [],
    casConfig: state => state.appInfo?.login?.casconfig,
    samlConfig: state => state.appInfo?.login?.saml2config,
    oauthConfig: state => state.appInfo?.login?.oauth2config,
    isBasicAuthEnabled: state => !!state.appInfo?.login?.basic?.enable,
    isLDAPAuthEnabled: state => !!state.appInfo?.login?.ldapconfig?.enable,
    isCASAuthEnabled: state => !!state.appInfo?.login?.casconfig?.enable,
    isSAMLAuthEnabled: state => !!state.appInfo?.login?.saml2config?.enable,
    isOauthAuthEnabled: state => !!state.appInfo?.login?.oauth2config?.enable,
    eventsCountTriggerDefaultThreshold: state => state.appInfo?.events_count_trigger_default_threshold,
    disabledTransitions: state => state.appInfo?.disabled_transitions,
  },
  mutations: {
    [types.FETCH_APP_INFO](state) {
      state.pending = true;
    },
    [types.FETCH_APP_INFO_COMPLETED](state, { appInfo }) {
      state.appInfo = appInfo;
      state.pending = false;
    },
    [types.FETCH_APP_INFO_FAILED](state) {
      state.pending = false;
    },
  },
  actions: {
    async fetchAppInfo({ commit, dispatch }) {
      try {
        commit(types.FETCH_APP_INFO);

        const appInfo = await request.get(API_ROUTES.infos.app);
        const preparedAppInfo = {
          ...appInfo,
          disabled_transitions: String(localStorage.get('disabled_transitions')) === 'true',
        }; // TODO: remove it

        commit(types.FETCH_APP_INFO_COMPLETED, { appInfo: preparedAppInfo }); // TODO: remove it

        if (appInfo.language) {
          dispatch('i18n/setGlobalLocale', appInfo.language, { root: true });
        }

        if (appInfo.popup_timeout) {
          dispatch('setPopupTimeouts', { popupTimeout: appInfo.popup_timeout });
        }
      } catch (err) {
        commit(types.FETCH_APP_INFO_FAILED);

        throw err;
      }
    },

    updateUserInterface(context, { data } = {}) {
      localStorage.set('disabled_transitions', data.disabled_transitions);

      return request.post(API_ROUTES.infos.userInterface, data);
    },

    updateMaintenanceMode(context, { data } = {}) {
      return request.put(API_ROUTES.maintenance, data);
    },

    setPopupTimeouts({ dispatch }, { popupTimeout = {} }) {
      const timeInfo = durationToSeconds(popupTimeout.info) * 1000;
      const timeError = durationToSeconds(popupTimeout.error) * 1000;

      dispatch('popups/setDefaultCloseTime', { type: POPUP_TYPES.success, time: timeInfo }, { root: true });
      dispatch('popups/setDefaultCloseTime', { type: POPUP_TYPES.info, time: timeInfo }, { root: true });
      dispatch('popups/setDefaultCloseTime', { type: POPUP_TYPES.warning, time: timeError }, { root: true });
      dispatch('popups/setDefaultCloseTime', { type: POPUP_TYPES.error, time: timeError }, { root: true });
    },
  },
};
