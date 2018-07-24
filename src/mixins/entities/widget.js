import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('view/widget');

export default {
  computed: {
    ...mapGetters({
      widgets: 'items',
    }),
  },
  methods: {
    ...mapActions({
      createWidget: 'create',
      updateWidget: 'update',
    }),
  },
};
