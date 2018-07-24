import Vue from 'vue';
import Vuex from 'vuex';

import appModule from './modules/app';
import i18nModule from './modules/i18n';
import AuthModule from './modules/auth';
import modalModule from './modules/modal';
import popupModule from './modules/popup';
import eventModule from './modules/event';
import mFilterEditorModule from './modules/mfilter-editor';
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
    app: appModule,
    i18n: i18nModule,
    modal: modalModule,
    popup: popupModule,
    event: eventModule,
    auth: AuthModule,
    mFilterEditor: mFilterEditorModule,

    ...entitiesModules,
  },
  plugins: [entitiesPlugin],
});
