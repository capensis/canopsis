<template>
  <widget-settings-group :title="$t('settings.counterLevels.title')">
    <widget-settings-item :title="$t('settings.counterLevels.fields.counter')">
      <c-select
        v-model="counterModel"
        :label="$t('settings.counterLevels.fields.counter')"
        :items="availableCounters"
        class="select"
        hide-details
        single-line
      />
    </widget-settings-item>

    <field-criticity-levels v-model="valuesModel" />

    <field-levels-colors-selector
      v-model="colorsModel"
      color-type="hex"
      hide-suffix
    />
  </widget-settings-group>
</template>

<script>
import { AVAILABLE_COUNTERS } from '@/constants';

import FieldCriticityLevels from '@/components/sidebars/stats/form/fields/criticity-levels.vue';
import FieldLevelsColorsSelector from '@/components/sidebars/stats/form/fields/levels-colors-selector.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: {
    WidgetSettingsItem,
    WidgetSettingsGroup,
    FieldCriticityLevels,
    FieldLevelsColorsSelector,
  },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:modelValue'],

  computed: {
    availableCounters() {
      return Object.values(AVAILABLE_COUNTERS);
    },

    counterModel: {
      get() {
        return this.modelValue.counter;
      },
      set(counter) {
        this.updateModel({ counter });
      },
    },

    valuesModel: {
      get() {
        return this.modelValue.values;
      },
      set(values) {
        this.updateModel({ values });
      },
    },

    colorsModel: {
      get() {
        return this.modelValue.colors;
      },
      set(colors) {
        this.updateModel({ colors });
      },
    },
  },

  methods: {
    updateModel(patch) {
      this.$emit('update:modelValue', {
        ...this.modelValue,
        ...patch,
      });
    },
  },
};
</script>
