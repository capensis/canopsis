import Vue from 'vue';
import Vuex from 'vuex';

import i18nModule from './modules/i18n';
import mFilterEditorModule from './modules/mfilter-editor';
import entitiesModule from './modules/entities';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    i18n: i18nModule,
    mFilterEditor: mFilterEditorModule,
    entities: entitiesModule,
  },
});
