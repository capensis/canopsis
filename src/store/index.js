import Vue from 'vue';
import Vuex from 'vuex';

import i18nModule from './modules/i18n';
import modalModule from './modules/modal';
import eventModule from './modules/event';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    i18n: i18nModule,
    modal: modalModule,
    event: eventModule,
  },
});
