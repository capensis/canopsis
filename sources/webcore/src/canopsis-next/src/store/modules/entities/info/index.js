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
  },
  getters: {
    version: state => state.version,
    logo: state => state.logo,
    appTitle: state => state.appTitle,
    footer: state => state.footer,
  },
  mutations: {
    [types.FETCH_LOGIN_INFOS](state, {
      version,
      userInterface,
    }) {
      state.version = version;
      state.logo = userInterface.logo || '';
      state.appTitle = userInterface.app_title || '';
      state.footer = userInterface.footer || '';
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
        } = await request.get(API_ROUTES.infos.login);

        commit(types.FETCH_LOGIN_INFOS, { version, userInterface });
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
