import { API_ROUTES } from '@/config';
import request from '@/services/request';

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
    casConfig: {},
  },
  getters: {
    version: state => state.version,
    logo: state => state.logo,
    appTitle: state => state.appTitle,
    footer: state => state.footer,
    edition: state => state.edition,
    stack: state => state.stack,
    description: state => state.description,
    language: state => state.language,
    isLDAPAuthEnabled: state => state.isLDAPAuthEnabled,
    isCASAuthEnabled: state => state.isCASAuthEnabled,
    casConfig: state => state.casConfig,
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

      state.isLDAPAuthEnabled = loginConfig.ldapconfig ? loginConfig.ldapconfig.enable : false;
      state.isCASAuthEnabled = loginConfig.casconfig ? loginConfig.casconfig.enable : false;

      state.casConfig = loginConfig.casconfig;
    },
    [types.FETCH_APP_INFOS](state, {
      version,
      logo,
      language,
      appTitle,
      edition,
      stack,
    }) {
      state.version = version;
      state.logo = logo;
      state.appTitle = appTitle;
      state.edition = edition;
      state.stack = stack;
      state.language = language;
    },
  },
  actions: {
    async fetchLoginInfos({ commit }) {
      try {
        const {
          version,
          user_interface: userInterface,
          login_config: loginConfig,
        } = await request.get(API_ROUTES.infos.login);

        commit(types.FETCH_LOGIN_INFOS, {
          version,
          userInterface: userInterface || {},
          loginConfig: loginConfig || {},
        });
      } catch (err) {
        console.error(err);
      }
    },

    async fetchAppInfos({ commit, dispatch }) {
      try {
        const {
          version,
          logo,
          language,
          app_title: appTitle,
          edition,
          stack,
        } = await request.get(API_ROUTES.infos.app);

        commit(
          types.FETCH_APP_INFOS,
          {
            version,
            logo,
            appTitle,
            edition,
            language,
            stack,
          },
        );

        if (language) {
          dispatch('i18n/setGlobalLocale', language, { root: true });
        }
      } catch (err) {
        console.error(err);
      }
    },

    updateUserInterface(context, { data } = {}) {
      return request.post(API_ROUTES.infos.userInterface, data);
    },
  },
};
