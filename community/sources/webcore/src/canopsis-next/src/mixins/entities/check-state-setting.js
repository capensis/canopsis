import { pick } from 'lodash';
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
        if (!data.name) {
          return;
        }

        this.stateSettingPending = true;
        const response = await this.checkEntityStateSetting({
          data: pick(data, ['name', 'type', 'infos', 'impact_level']),
        });

        this.stateSetting = response?.title ? response : undefined;
      } catch (err) {
        console.error(err);
      } finally {
        this.stateSettingPending = false;
      }
    },
  },
};
