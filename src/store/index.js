import Vue from 'vue';
import Vuex from 'vuex';

import i18nModule from './modules/i18n';
import popupModule from './modules/popup';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    i18n: i18nModule,
    popup: popupModule,
  },
});
