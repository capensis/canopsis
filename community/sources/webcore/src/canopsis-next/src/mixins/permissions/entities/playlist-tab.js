import { createNamespacedHelpers } from 'vuex';

import { SCHEMA_EMBEDDED_KEY } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { authMixin } from '@/mixins/auth';

const { mapGetters: mapEntitiesGetters } = createNamespacedHelpers('entities');

export const permissionsEntitiesPlaylistTabMixin = {
  mixins: [authMixin],
  computed: {
    ...mapEntitiesGetters({
      getEntitiesList: 'getList',
    }),
  },
  methods: {
    getAvailableTabsByIds(tabsIds) {
      const tabs = this.getEntitiesList(ENTITIES_TYPES.viewTab, tabsIds, true);

      return tabs.filter(tab => tab[SCHEMA_EMBEDDED_KEY].parents.some(parent => this.checkReadAccess(parent.id)));
    },
  },
};
