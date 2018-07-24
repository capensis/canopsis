import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('view/widget');

/**
 * @mixin Helpers for the widget entity
 */
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
