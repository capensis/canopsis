import { createNamespacedHelpers } from 'vuex';

import { WIDGET_MAX_SIZE, WIDGET_MIN_SIZE } from '@/constants';

const { mapGetters, mapActions } = createNamespacedHelpers('view');

/**
 * @mixin Helpers for the view entity
 */
export default {
  computed: {
    ...mapGetters({
      view: 'item',
    }),

    getWidgetAvailableRows() {
      return widgetId => this.view.rows.map((row) => {
        const availableSize = row.widgets.reduce((acc, widget) => {
          if (widget._id !== widgetId) {
            acc.sm -= widget.size.sm;
            acc.md -= widget.size.md;
            acc.lg -= widget.size.lg;
          }

          return acc;
        }, { sm: WIDGET_MAX_SIZE, md: WIDGET_MAX_SIZE, lg: WIDGET_MAX_SIZE });

        return {
          _id: row._id,
          title: row.title,

          availableSize,
        };
      }).filter(({ availableSize }) =>
        availableSize.sm >= WIDGET_MIN_SIZE &&
        availableSize.md >= WIDGET_MIN_SIZE &&
        availableSize.lg >= WIDGET_MIN_SIZE);
    },
  },
  methods: {
    ...mapActions({
      fetchView: 'fetchItem',
      createView: 'create',
    }),
  },
};
