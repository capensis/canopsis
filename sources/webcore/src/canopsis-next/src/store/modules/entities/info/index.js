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
    isLDAPAuthEnabled: false,
    isCASAuthEnabled: false,
    casConfig: {},
  },
  getters: {
    version: state => state.version,
    logo: state => state.logo,
    appTitle: state => state.appTitle,
    footer: state => state.footer,
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

      state.isLDAPAuthEnabled = loginConfig.ldapconfig ? loginConfig.ldapconfig.enable : false;
      state.isCASAuthEnabled = loginConfig.casconfig ? loginConfig.casconfig.enable : false;

      state.casConfig = loginConfig.casconfig;
    },
    [types.FETCH_APP_INFOS](state, {
      version,
      logo,
      appTitle,
    }) {
      state.version = version;
      state.logo = logo;
      state.appTitle = appTitle;
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

        commit(types.FETCH_LOGIN_INFOS, { version, userInterface, loginConfig });
      } catch (err) {
        console.error(err);
      }
    },

    async fetchAppInfos({ commit }) {
      try {
        const { version, logo, app_title: appTitle } = await request.get(API_ROUTES.infos.app);

        commit(types.FETCH_APP_INFOS, { version, logo, appTitle });
      } catch (err) {
        console.error(err);
      }
    },
  },
};
