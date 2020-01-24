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
  methods: {
    ...mapActions(['fetchAlarmColumnFilters']),

    getFilter(filter, attributes = []) {
      const filterFunc = this.$options.filters[filter];

      return value => (filterFunc ? filterFunc(value, ...attributes) : value);
    },
  },
};
