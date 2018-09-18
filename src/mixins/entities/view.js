import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view');

const MAX_COLUMNS = 12;

/**
 * @mixin Helpers for the view entity
 */
export default {
  computed: {
    ...mapGetters({
      view: 'item',
    }),

    availableRows() {
      return this.view.rows.map((row) => {
        const availableColumns = row.widgets.reduce((acc, widget) => {
          acc.sm -= widget.columnSM;
          acc.md -= widget.columnMD;
          acc.lg -= widget.columnLG;

          return acc;
        }, { sm: MAX_COLUMNS, md: MAX_COLUMNS, lg: MAX_COLUMNS });

        return {
          _id: row._id,
          title: row.title,

          availableColumns,
        };
      }).filter(({ availableColumns }) =>
        availableColumns.sm <= 9 && availableColumns.md <= 9 && availableColumns.lg <= 9);
    },
  },
  methods: {
    ...mapActions({
      fetchView: 'fetchItem',
      createView: 'create',
    }),
  },
};
