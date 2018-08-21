<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title")
      v-divider
      field-default-column-sort(v-model="settings.widget.default_sort_column")
      v-divider
      field-context-entities-types-filter(v-model="settings.widget_preferences.selectedTypes")
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed, fixed, right) {{ $t('common.save') }}
</template>

<script>
import pick from 'lodash/pick';
import cloneDeep from 'lodash/cloneDeep';

import FieldTitle from '@/components/other/settings/fields/title.vue';
import FieldDefaultColumnSort from '@/components/other/settings/fields/default-column-sort.vue';
import FieldContextEntitiesTypesFilter from '@/components/other/settings/fields/context-entities-types-filter.vue';

import widgetSettingsMixin from '@/mixins/widget/settings';

/**
 * Component to regroup the entities list settings fields
 *
 * @prop {Object} widget - active widget
 * @prop {bool} isNew - is widget new
 *
 * @event closeSettings#click
 */
export default {
  components: {
    FieldTitle,
    FieldDefaultColumnSort,
    FieldContextEntitiesTypesFilter,
  },
  mixins: [
    widgetSettingsMixin,
  ],
  data() {
    return {
      settings: {
        widget: {
          title: this.widget.title,
          default_sort_column: cloneDeep(this.widget.default_sort_column),
        },
        widget_preferences: {
          selectedTypes: [],
        },
      },
    };
  },
  created() {
    this.settings.widget_preferences = pick(this.userPreference.widget_preferences, 'selectedTypes');
  },
};
</script>

<style scoped>
  .closeIcon:hover {
    cursor: pointer;
  }
</style>
