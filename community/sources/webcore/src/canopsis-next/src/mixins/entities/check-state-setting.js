import { createNamespacedHelpers } from 'vuex';

const { mapActions: mapEntityActions } = createNamespacedHelpers('entity');

export const checkStateSettingMixin = {
  data() {
    return {
      stateSettingPending: true,
      stateSetting: undefined,
    };
  },
  methods: {
    ...mapEntityActions({
      checkEntityStateSetting: 'checkStateSetting',
    }),

    async checkStateSetting(data) {
      try {
        this.stateSettingPending = true;
        const response = await this.checkEntityStateSetting({ data });

        this.stateSetting = response?.title ? response : undefined;
      } catch (err) {
        console.error(err);
      } finally {
        this.stateSettingPending = false;
      }
    },
  },
};
