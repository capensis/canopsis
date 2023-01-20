import { ALARM_FIELDS } from '@/constants';

import { getInfosWidgetColumn, isLinksWidgetColumn } from '@/helpers/forms/shared/widget-column';

export const widgetColumnsMixin = {
  computed: {
    hasColumns() {
      return this.columns.length > 0;
    },
  },
  methods: {
    getColumnLabel(column, labelsMap) {
      if (column.label) {
        return column.label;
      }

      const infosColumn = getInfosWidgetColumn(column.value);

      if (infosColumn) {
        return this.$tc(labelsMap[infosColumn], 2);
      }

      if (isLinksWidgetColumn(column.value)) {
        return this.$tc(labelsMap[ALARM_FIELDS.links], 2);
      }

      return this.$tc(labelsMap[column.value], 2);
    },

    getSortable(column, unsortableFields) {
      return !unsortableFields.some(field => column.value.startsWith(field));
    },
  },
};
