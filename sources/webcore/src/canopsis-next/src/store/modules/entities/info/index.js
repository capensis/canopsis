import { API_ROUTES } from '@/config';
import request from '@/services/request';

const types = {
  FETCH_LOGIN_INFOS: 'FETCH_LOGIN_INFOS',
  FETCH_LOGIN_INFOS_COMPLETED: 'FETCH_LOGIN_INFOS_COMPLETED',
  FETCH_LOGIN_INFOS_FAILED: 'FETCH_LOGIN_INFOS_FAILED',
};

export default {
  namespaced: true,
  state: {
    version: '',
    loginConfig: {},
    userInterface: {},
  },
  getters: {
    version: state => state.version,
  },
  mutations: {
    [types.FETCH_LOGIN_INFOS]() {

    },
    [types.FETCH_LOGIN_INFOS_COMPLETED](state, {
      version,
      userInterface,
      loginConfig,
    }) {
      state.version = version;
      state.userInterface = userInterface;
      state.loginConfig = loginConfig;
    },
    [types.FETCH_LOGIN_INFOS_FAILED]() {

    },
  },
  actions: {
    async fetchLoginInfos({ commit }) {
      commit(types.FETCH_LOGIN_INFOS);

      try {
        const {
          version,
          user_interface: userInterface,
          login_config: loginConfig,
        } = await request.get(API_ROUTES.infos.login);

        commit(types.FETCH_LOGIN_INFOS_COMPLETED, { version, userInterface, loginConfig });
      } catch (err) {
        commit(types.FETCH_LOGIN_INFOS_FAILED, err);
        console.warn(err);
      }
    },
  },
};
