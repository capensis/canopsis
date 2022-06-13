<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile {{ $t('settings.filters') }}
    v-container
      v-layout(column)
        filter-selector(
          v-field="value",
          :label="$t('filterSelector.defaultFilter')",
          :filters="filters"
        )
        filters-list(
          :filters="filters",
          :addable="addable",
          :editable="editable"
        )
</template>

<script>
import { authMixin } from '@/mixins/auth';

import FilterSelector from '@/components/other/filter/filter-selector.vue';
import FiltersList from '@/components/other/filter/filters-list.vue';

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
    addable: {
      type: Boolean,
      default: false,
    },
    editable: {
      type: Boolean,
      default: false,
    },
    withAlarm: {
      type: Boolean,
      default: false,
    },
    withEntity: {
      type: Boolean,
      default: false,
    },
    withPbehavior: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
