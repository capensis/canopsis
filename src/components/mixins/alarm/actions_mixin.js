import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('alarm/actions');

export default {
  methods: {
    ...mapActions({
      cancelAlarmConfirmation: 'cancel',
    }),
  },
};
