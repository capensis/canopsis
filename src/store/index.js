import Vue from 'vue';
import Vuex from 'vuex';

import i18nModule from '@/store/modules/i18n';
import appModule from '@/store/modules/app';
import entitiesModule from '@/store/modules/entities';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    app: appModule,
    i18n: i18nModule,
    entities: entitiesModule,
  },
});
