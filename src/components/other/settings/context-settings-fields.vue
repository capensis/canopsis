<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.title")
      v-divider
      field-default-column-sort(v-model="settings.defaultSortColumn")
      v-divider
      field-context-entities-types-filter(v-model="settings.selectedTypes")
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed, fixed, right) {{ $t('common.save') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import cloneDeep from 'lodash/cloneDeep';
import find from 'lodash/find';

import FieldTitle from '@/components/other/settings/fields/title.vue';
import FieldDefaultColumnSort from '@/components/other/settings/fields/default-column-sort.vue';
import FieldContextEntitiesTypesFilter from '@/components/other/settings/fields/context-entities-types-filter.vue';

const { mapGetters } = createNamespacedHelpers('entities');

/**
* Component to regroup the entities list settings fields
*/
export default {
  components: {
    FieldTitle,
    FieldDefaultColumnSort,
    FieldContextEntitiesTypesFilter,
  },
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      settings: {
        title: this.widget.title,
        defaultSortColumn: cloneDeep(this.widget.default_sort_column),
        selectedTypes: [],
      },
    };
  },
  computed: {
    ...mapGetters(['getItem']),

    userPreference() {
      return this.getItem('userPreference', `${this.widget.id}_root`); // TODO: fix it
    },
  },
  created() {
    const filter = find(this.userPreference.widget_preferences.user_filters, { title: 'default_type_filter' });

    if (filter) {
      const filterData = JSON.parse(filter.filter);
      console.warn(filterData);
      this.settings.selectedTypes = filter.types;
    }
  },
  methods: {
    submit() {
      const widget = {
        ...this.widget,
        title: this.settings.title,
        default_sort_column: this.settings.defaultSortColumn,
      };

      const userPreference = { ...this.userPreference };
      const userFilters = userPreference.widget_preferences.user_filters || [];

      const defaultTypeFilterIndex = userFilters.findIndex(filter => filter.title === 'default_type_filter');

      if (defaultTypeFilterIndex > 0) {
        if (this.settings.selectedTypes.length) {
          userPreference.widget_preferences.user_filters[defaultTypeFilterIndex] = {
            title: 'default_type_filter',
            filter: JSON.stringify({
              $or: this.settings.selectedTypes.map(type => ({ type })),
            }),
          };
        } else {
          delete userPreference.widget_preferences.user_filters[defaultTypeFilterIndex];
        }
      } else if (this.settings.selectedTypes.length) {
        userPreference.widget_preferences.user_filters = [...userFilters, {
          title: 'default_type_filter',
          filter: JSON.stringify({
            $or: this.settings.selectedTypes.map(type => ({ type })),
          }),
        }];
      }

      console.warn(widget, userPreference);
    },
  },
};
</script>

<style scoped>
  .closeIcon:hover {
    cursor: pointer;
  }
</style>
