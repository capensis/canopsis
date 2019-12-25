import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('alarmColumnFilters');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      alarmColumnFiltersPending: 'pending',
      alarmColumnFilters: 'alarmColumnFilters',
    }),

    columnFiltersMap() {
      return this.alarmColumnFilters.reduce((acc, { column, filter, attributes = [] }) => {
        acc[column] = this.getFilter(filter, attributes);

        return acc;
      }, {});
    },
  },
  mounted() {
    this.fetchAlarmColumnFilters();
  },
  methods: {
    ...mapActions(['fetchAlarmColumnFilters']),

    getFilter(filter, attributes = []) {
      return value => (this.$options.filters[filter]
        ? this.$options.filters[filter](value, ...attributes)
        : value);
    },
  },
};
