import { ALARM_ADVANCED_SEARCH_FIELDS, ALARM_FIELDS, ALARM_FIELDS_TO_LABELS_KEYS } from '@/constants';

import { widgetAdvancedSearchCommonMixin } from './common';

export const widgetAdvancedSearchAlarmFieldsMixin = {
  mixins: [widgetAdvancedSearchCommonMixin],
  computed: {
    advancedSearchFields() {
      return ALARM_ADVANCED_SEARCH_FIELDS.map((field) => {
        const text = this.$tc(ALARM_FIELDS_TO_LABELS_KEYS[field], 2);
        const result = { text, value: field };

        if (field === ALARM_FIELDS.entityInfos) {
          result.items = this.entityInfosKeys.map(({ value }) => {
            const variableValue = `${ALARM_FIELDS.entityInfos}.${value}`;
            const variableSelectedText = `${text}.${value}`;

            return {
              value: variableValue,
              text: value,
              selectorText: variableSelectedText,
              items: this.getDeepInfosItems(variableValue, variableSelectedText),
            };
          });
        }

        return result;
      });
    },
  },
};
