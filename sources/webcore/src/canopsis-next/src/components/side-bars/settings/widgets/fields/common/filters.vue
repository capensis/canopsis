<template lang="pug">
  v-list-group(data-test="filters")
    v-list-tile(slot="activator") {{ $t('settings.filters') }}
    v-container
      filter-selector(
        :label="$t('filterSelector.defaultFilter')",
        :entitiesType="this.entitiesType",
        :value="value",
        :filters="filters",
        :condition="condition",
        :hideSelect="hideSelect",
        :hasAccessToAddFilter="hasAccessToAddFilter",
        :hasAccessToEditFilter="hasAccessToEditFilter",
        hideSelectIcon,
        long,
        @input="$emit('input', $event)",
        @update:condition="$emit('update:condition', $event)",
        @update:filters="updateFilters"
      )
</template>

<script>
import { isUndefined } from 'lodash';

import { FILTER_DEFAULT_VALUES, ENTITIES_TYPES } from '@/constants';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';

import FilterSelector from '@/components/other/filter/selector/filter-selector.vue';

export default {
  components: { FilterSelector },
  mixins: [authMixin, modalMixin],
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
    hasAccessToAddFilter: {
      type: Boolean,
      default: true,
    },
    hasAccessToEditFilter: {
      type: Boolean,
      default: true,
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
