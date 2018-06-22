<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.title")
      v-divider
      field-default-column-sort(v-model="settings.default_sort_column")
      v-divider
      field-columns
      v-divider
      field-periodic-refresh
      v-divider
      field-default-elements-per-page
      v-divider
      field-opened-resolved-filter
      v-divider
      field-filters
      v-divider
      field-info-popup
      v-divider
      field-more-info
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed, fixed, right) {{ $t('common.save') }}
</template>

<script>
// LIBS
import cloneDeep from 'lodash/cloneDeep';
// COMPONENTS
import FieldTitle from '@/components/other/settings/fields/title.vue';
import FieldDefaultColumnSort from '@/components/other/settings/fields/default-column-sort.vue';
import FieldColumns from '@/components/other/settings/fields/columns.vue';
import FieldPeriodicRefresh from '@/components/other/settings/fields/periodic-refresh.vue';
import FieldDefaultElementsPerPage from '@/components/other/settings/fields/default-elements-per-page.vue';
import FieldOpenedResolvedFilter from '@/components/other/settings/fields/opened-resolved-filter.vue';
import FieldFilters from '@/components/other/settings/fields/filters.vue';
import FieldInfoPopup from '@/components/other/settings/fields/info-popup.vue';
import FieldMoreInfo from '@/components/other/settings/fields/more-info.vue';

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

export default {
  components: {
    FieldTitle,
    FieldDefaultColumnSort,
    FieldColumns,
    FieldPeriodicRefresh,
    FieldDefaultElementsPerPage,
    FieldOpenedResolvedFilter,
    FieldFilters,
    FieldInfoPopup,
    FieldMoreInfo,
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
