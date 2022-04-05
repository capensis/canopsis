<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile {{ $t('settings.filters') }}
    v-container
      v-layout(column)
        filter-selector(
          v-field="value",
          :label="$t('filterSelector.defaultFilter')",
          :filters="filters",
          :condition="condition",
          :hide-prepend="hidePrepend",
          @update:condition="$emit('update:condition', $event)"
        )
        filters-list(
          v-if="widgetId",
          :widget-id="widgetId",
          :addable="addable",
          :editable="editable",
          with-alarm
        )
</template>

<script>
import { isUndefined } from 'lodash';

import { FILTER_DEFAULT_VALUES } from '@/constants';

import { authMixin } from '@/mixins/auth';

import FilterSelector from '@/components/other/filter/filter-selector.vue';
import FiltersList from '@/components/other/filter/filters-list.vue';

// TODO: add withAlarm and etc properties
export default {
  components: { FilterSelector, FiltersList },
  mixins: [authMixin],
  props: {
    widgetId: {
      type: String,
      required: false,
    },
    filters: {
      type: Array,
      default: () => [],
    },
    value: {
      type: [Object, Array],
      default: null,
    },
    condition: {
      type: String,
      default: FILTER_DEFAULT_VALUES.condition,
    },
    hidePrepend: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    editable: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    updateFilters(filters, value) {
      this.$emit('update:filters', filters);

      if (!isUndefined(value)) {
        this.$emit('input', value);
      }
    },
  },
};
</script>
