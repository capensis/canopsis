import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('stateSetting');

export const entitiesStateSettingMixin = {
  computed: {
    ...mapGetters({
      stateSettings: 'items',
      stateSettingsPending: 'pending',
    }),
  },
  methods: {
    ...mapActions({
      fetchStateSettingsList: 'fetchList',
      updateStateSetting: 'update',
    }),
  },
};
