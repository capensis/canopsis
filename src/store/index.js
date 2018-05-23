import Vue from 'vue';
import Vuex from 'vuex';

import i18nModule from '@/store/modules/i18n';
import AuthModule from '@/store/modules/auth';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    i18n: i18nModule,
    auth: AuthModule,
  },
});
