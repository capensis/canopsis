import { ALARM_ENTITY_FIELDS } from '@/constants';

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
        sortable: ALARM_ENTITY_FIELDS.extraDetails !== value,
        popupTemplate: this.infoPopupsMap[value],
      };
    },
  },
};
