import Vue from 'vue';
import Vuex from 'vuex';

import i18nModule from './modules/i18n';
import appModule from './modules/app';
import entitiesModule from './modules/entities';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    app: appModule,
    i18n: i18nModule,
    entities: entitiesModule,
  },
});
