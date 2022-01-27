<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile {{ $t('settings.filters') }}
    v-container
      filter-selector(
        v-field="value",
        :label="$t('filterSelector.defaultFilter')",
        :entities-type="entitiesType",
        :filters="filters",
        :condition="condition",
        :hide-select="hideSelect",
        :has-access-to-add-filter="addable",
        :has-access-to-edit-filter="editable",
        hide-select-icon,
        long,
        @update:condition="$emit('update:condition', $event)",
        @update:filters="updateFilters"
      )
</template>

<script>
import { isUndefined } from 'lodash';

import { FILTER_DEFAULT_VALUES, ENTITIES_TYPES } from '@/constants';

import { authMixin } from '@/mixins/auth';

import FilterSelector from '@/components/other/filter/filter-selector.vue';

export default {
  components: { FilterSelector },
  mixins: [authMixin],
  props: {
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
    hideSelect: {
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
    entitiesType: {
      type: String,
      default: ENTITIES_TYPES.alarm,
      validator: value => [ENTITIES_TYPES.alarm, ENTITIES_TYPES.entity, ENTITIES_TYPES.pbehavior].includes(value),
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
