<template>
  <widget-settings-item :title="$t('settings.filters')">
    <v-layout column>
      <filter-selector
        v-if="!hideSelector"
        v-field="value"
        :label="$t('filter.selector.defaultFilter')"
        :filters="filters"
        hide-multiply
      />
      <field-filters-list
        v-bind="listProps"
        :filters="filters"
        :addable="addable"
        :editable="editable"
        @input="updateFilters"
      />
    </v-layout>
  </widget-settings-item>
</template>

<script>
import { pick } from 'lodash';

import { formBaseMixin } from '@/mixins/form';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import FieldFiltersList from '@/components/sidebars/form/fields/filters-list.vue';

export default {
  components: {
    WidgetSettingsItem,
    FilterSelector,
    FieldFiltersList,
  },
  mixins: [formBaseMixin],
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
      type: String,
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
    withServiceWeather: {
      type: Boolean,
      default: false,
    },
    hideSelector: {
      type: Boolean,
      default: false,
    },
    entityTypes: {
      type: Array,
      required: false,
    },
    entityCountersType: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    listProps() {
      return pick(this, [
        'withAlarm',
        'withEntity',
        'withPbehavior',
        'withServiceWeather',
        'entityTypes',
        'entityCountersType',
      ]);
    },
  },
  methods: {
    updateFilters(filters) {
      if (this.value && !filters.some(filter => filter._id === this.value)) {
        this.updateModel(null);
      }

      this.$emit('update:filters', filters);
    },
  },
};
</script>
