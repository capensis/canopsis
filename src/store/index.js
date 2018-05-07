import Vue from 'vue';
import Vuex from 'vuex';

import i18nModule from './modules/i18n';
import modalModule from './modules/modal';
import alarmsModule from './modules/alarm';

Vue.use(Vuex);

const store = new Vuex.Store({
  modules: {
    i18n: i18nModule,
    modal: modalModule,
    alarm: alarmsModule,
  },
});

export default store;
