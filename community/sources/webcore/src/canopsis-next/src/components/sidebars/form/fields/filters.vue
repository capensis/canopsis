<template>
  <widget-settings-item :title="$t('settings.filters')">
    <c-layout gap="2" column>
      <filter-selector
        v-if="!hideSelector"
        :model-value="modelValue"
        :label="$t('filter.selector.defaultFilter')"
        :filters="filters"
        hide-multiply
        @update:model-value="$emit('update:modelValue', $event)"
      />

      <field-filters-list
        v-bind="listProps"
        :filters="filters"
        :addable="addable"
        :editable="editable"
        @update:filters="updateFilters"
      />
    </c-layout>
  </widget-settings-item>
</template>

<script>
import { pick } from 'lodash';

import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

import FieldFiltersList from './filters-list.vue';

export default {
  components: {
    WidgetSettingsItem,
    FilterSelector,
    FieldFiltersList,
  },

  props: {
    widgetId: {
      type: String,
      required: false,
    },
    widgetType: {
      type: String,
      required: false,
    },
    filters: {
      type: Array,
      default: () => [],
    },
    modelValue: {
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

  emits: ['update:modelValue', 'update:filters'],

  computed: {
    listProps() {
      return pick(this, [
        'withAlarm',
        'withEntity',
        'withPbehavior',
        'withServiceWeather',
        'entityTypes',
        'entityCountersType',
        'widgetId',
        'widgetType',
      ]);
    },
  },

  methods: {
    updateFilters(filters) {
      if (this.modelValue && !filters.some(filter => filter._id === this.modelValue)) {
        this.$emit('update:modelValue', null);
      }

      this.$emit('update:filters', filters);
    },
  },
};
</script>
