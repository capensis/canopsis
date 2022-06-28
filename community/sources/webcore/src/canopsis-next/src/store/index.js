import Vue from 'vue';
import Vuex from 'vuex';

import featuresService from '@/services/features';

import authModule from './modules/auth';
import i18nModule from './modules/i18n';
import eventModule from './modules/event';
import queryModule from './modules/query';
import navigationModule from './modules/navigation';
import entitiesModules from './modules/entities';
import activeViewModule from './modules/active-view';

import entitiesPlugin from './plugins/entities';
import watchOncePlugin from './plugins/watch-once';

Vue.use(Vuex);
/**
 * @typedef {Object} VuexActionContext
 * @property {Object} state
 * @property {Object} rootState
 * @property {Object} getters
 * @property {Object} rootGetters
 * @property {function} commit
 * @property {function} dispatch
 */

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
    event: eventModule,
    query: queryModule,
    navigation: navigationModule,
    activeView: activeViewModule,

    ...entitiesModules,
    ...featuresService.get('store.modules'),
  },
  plugins: [entitiesPlugin, watchOncePlugin],
});
