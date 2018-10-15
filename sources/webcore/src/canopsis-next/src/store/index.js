import Vue from 'vue';
import Vuex from 'vuex';

import authModule from './modules/auth';
import i18nModule from './modules/i18n';
import modalModule from './modules/modal';
import popupModule from './modules/popup';
import eventModule from './modules/event';
import queryModule from './modules/query';
import sideBarModule from './modules/side-bar';
import entitiesModules from './modules/entities';

import entitiesPlugin from './plugins/entities';

Vue.use(Vuex);

/**
 * @typedef {Object} ActionContext
 * @property {function} commit
 * @property {function} dispatch
 * @property {Object} getters
 * @property {Object} rootGetters
 * @property {Object} state
 * @property {Object} rootState
 */

export default new Vuex.Store({
  modules: {
    auth: authModule,
    i18n: i18nModule,
    modal: modalModule,
    popup: popupModule,
    event: eventModule,
    query: queryModule,
    sideBar: sideBarModule,

    ...entitiesModules,
  },
  plugins: [entitiesPlugin],
});
