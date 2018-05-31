import Vue from 'vue';
import Vuex from 'vuex';

import appModule from '@/store/modules/app';
import i18nModule from '@/store/modules/i18n';
import AuthModule from '@/store/modules/auth';
import alarmsListSettingsModule from '@/store/modules/alarms-list-settings';
import modalModule from '@/store/modules/modal';
import eventModule from '@/store/modules/event';
import mFilterEditorModule from '@/store/modules/mfilter-editor';
import entitiesModules from '@/store/modules/entities';

import entitiesPlugin from '@/store/plugins/entities';


Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    app: appModule,
    i18n: i18nModule,
    alarmsListSettings: alarmsListSettingsModule,
    modal: modalModule,
    event: eventModule,
    auth: AuthModule,
    mFilterEditor: mFilterEditorModule,

    ...entitiesModules,
  },
  plugins: [entitiesPlugin],
});
