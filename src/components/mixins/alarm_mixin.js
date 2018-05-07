import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('alarm');

export default {
  methods: {
    ...mapActions([
      'cancelConfirmation',
    ]),
  },
  computed: {
    testComputed() {
      return true;
    },
  },
};
