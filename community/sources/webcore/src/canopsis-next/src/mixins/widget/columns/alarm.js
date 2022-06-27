import { formToColumnValue } from '@/helpers/forms/widgets/alarm';

import { widgetColumnsMixin } from './common';

export const widgetColumnsAlarmMixin = {
  mixins: [widgetColumnsMixin],
  computed: {
    infoPopupsMap() {
      return (this.widget.parameters?.infoPopups ?? []).reduce((acc, { column, template }) => {
        acc[formToColumnValue(column)] = template;

        return acc;
      }, {});
    },
  },
  methods: {
    mapColumnEntity({ label, value, ...column }) {
      return {
        ...column,

        value,
        text: label,
        popupTemplate: this.infoPopupsMap[value],
      };
    },
  },
};
