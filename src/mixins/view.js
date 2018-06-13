import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view');

export default {
  methods: {
    ...mapActions({
      fetchView: 'fetchItem',
    }),
  },
  computed: {
    ...mapGetters({
      getView: 'getItem',
      viewPending: 'pending',
    }),
  },
};
