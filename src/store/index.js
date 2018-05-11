import Vue from 'vue';
import Vuex from 'vuex';

import i18nModule from './modules/i18n';
<<<<<<< HEAD
<<<<<<< HEAD
import modalModule from './modules/modal';
import entitiesModule from './modules/entities';
=======
>>>>>>> parent of c7eef04... Add base modal structure and finish some modals
=======
>>>>>>> parent of c7eef04... Add base modal structure and finish some modals

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    i18n: i18nModule,
<<<<<<< HEAD
<<<<<<< HEAD
    modal: modalModule,
    entities: entitiesModule,
=======
>>>>>>> parent of c7eef04... Add base modal structure and finish some modals
=======
>>>>>>> parent of c7eef04... Add base modal structure and finish some modals
  },
});
