import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view');
const { mapGetters: entitiesMapGetters } = createNamespacedHelpers('entities');

const MAX_COLUMNS = 12;

/**
 * @mixin Helpers for the view entity
 */
export default {
  computed: {
    ...mapGetters({
      view: 'item',
    }),

    ...entitiesMapGetters([
      'getItem',
    ]),

    getWidgetAvailableRows() {
      return widgetId => this.view.rows.map((row) => {
        const availableColumns = row.widgets.reduce((acc, widget) => {
          if (widget._id !== widgetId) {
            acc.sm -= widget.columnSM;
            acc.md -= widget.columnMD;
            acc.lg -= widget.columnLG;
          }

          return acc;
        }, { sm: MAX_COLUMNS, md: MAX_COLUMNS, lg: MAX_COLUMNS });

        return {
          _id: row._id,
          title: row.title,

          availableColumns,
        };
      }).filter(({ availableColumns }) =>
        availableColumns.sm >= 3 && availableColumns.md >= 3 && availableColumns.lg >= 3);
    },
  },
  methods: {
    ...mapActions({
      fetchView: 'fetchItem',
      createView: 'create',
    }),
  },
};
