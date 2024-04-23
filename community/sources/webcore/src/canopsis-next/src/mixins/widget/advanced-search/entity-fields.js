import { ENTITY_FIELDS, ENTITY_ADVANCED_SEARCH_FIELDS, ENTITY_FIELDS_TO_LABELS_KEYS } from '@/constants';

import { widgetAdvancedSearchCommonMixin } from './common';

export const widgetAdvancedSearchEntityFieldsMixin = {
  mixins: [widgetAdvancedSearchCommonMixin],
  computed: {
    advancedSearchFields() {
      return ENTITY_ADVANCED_SEARCH_FIELDS.map((field) => {
        const text = this.$tc(ENTITY_FIELDS_TO_LABELS_KEYS[field], 2);
        const result = { text, value: field };

        if (field === ENTITY_FIELDS.infos) {
          result.items = this.entityInfosKeys.map(({ value }) => {
            const variableValue = `${ENTITY_FIELDS.infos}.${value}`;
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
