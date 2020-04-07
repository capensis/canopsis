import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('alarm');

export default {
  methods: {
    ...mapActions({
      fetchAlarmItem: 'fetchItem',
    }),

    fetchItemWithSteps() {
      const params = {
        sort_key: 't',
        sort_dir: 'DESC',
        limit: 1,
        with_steps: true,
      };

      if (this.alarm.v.resolved) {
        params.resolved = true;
      }

      return this.fetchAlarmItem({
        id: this.alarm._id,
        params,
      });
    },
  },
};
