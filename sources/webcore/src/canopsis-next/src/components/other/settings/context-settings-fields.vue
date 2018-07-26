<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.title")
      v-divider
      field-default-column-sort(v-model="settings.default_sort_column")
      v-divider
      field-context-entities-types-filter(v-model="settings.selected_types")
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed, fixed, right) {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

import FieldTitle from '@/components/other/settings/fields/title.vue';
import FieldDefaultColumnSort from '@/components/other/settings/fields/default-column-sort.vue';
import FieldContextEntitiesTypesFilter from '@/components/other/settings/fields/context-entities-types-filter.vue';

const storeSample = {
  default_sort_column: {
    direction: 'DESC',
    property: 'creation_date',
  },
  title: null,
  selected_types: [
    'connector',
    'watcher',
  ],
};

/**
* Component to regroup the entities list settings fields
*/
export default {
  components: {
    FieldTitle,
    FieldDefaultColumnSort,
    FieldContextEntitiesTypesFilter,
  },
  data() {
    return {
      settings: cloneDeep(storeSample),
    };
  },
  methods: {
    submit() {
      console.warn(this.settings);
    },
  },
};
</script>

<style scoped>
  .closeIcon:hover {
    cursor: pointer;
  }
</style>
