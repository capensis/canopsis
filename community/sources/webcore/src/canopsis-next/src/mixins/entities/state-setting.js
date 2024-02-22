import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('stateSetting');

export const entitiesStateSettingMixin = {
  computed: {
    ...mapGetters({
      stateSettings: 'items',
      stateSettingsPending: 'pending',
      stateSettingsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchStateSettingsList: 'fetchList',
      fetchStateSettingsListWithPreviousParams: 'fetchListWithPreviousParams',
      createStateSetting: 'create',
      updateStateSetting: 'update',
      removeStateSetting: 'remove',
    }),
  },
};
