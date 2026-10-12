<template>
  <widget-settings-group :title="$t('settings.margin.title')">
    <field-slider
      v-model="model.top"
      :title="$t('settings.margin.top')"
      :min="min"
      :max="max"
    />
    <field-slider
      v-model="model.right"
      :title="$t('settings.margin.right')"
      :min="min"
      :max="max"
    />
    <field-slider
      v-model="model.bottom"
      :title="$t('settings.margin.bottom')"
      :min="min"
      :max="max"
    />
    <field-slider
      v-model="model.left"
      :title="$t('settings.margin.left')"
      :min="min"
      :max="max"
    />
  </widget-settings-group>
</template>

<script>
import FieldSlider from '@/components/sidebars/form/fields/slider.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';

export default {
  components: { WidgetSettingsGroup, FieldSlider },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    min: {
      type: Number,
      default: 0,
    },
    max: {
      type: Number,
      default: 5,
    },
  },

  emits: ['update:modelValue'],

  computed: {
    model: {
      get() {
        return this.modelValue;
      },
      set(value) {
        this.$emit('update:modelValue', value);
      },
    },
  },
};
</script>
