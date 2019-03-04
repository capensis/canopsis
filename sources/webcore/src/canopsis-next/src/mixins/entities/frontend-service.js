import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('frontendService');

export default {
  computed: {
    ...mapGetters({
      frontendServiceItem: 'item',
    }),
  },
  methods: {
    ...mapActions({
      fetchFrontendService: 'fetch',
      updateFrontendService: 'update',
    }),
  },
};
