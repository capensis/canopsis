<template>
  <widget-settings-group :title="$t('settings.margin.title')">
    <field-slider
      v-field="form.top"
      :title="$t('settings.margin.top')"
      :min="min"
      :max="max"
    />
    <field-slider
      v-field="form.right"
      :title="$t('settings.margin.right')"
      :min="min"
      :max="max"
    />
    <field-slider
      v-field="form.bottom"
      :title="$t('settings.margin.bottom')"
      :min="min"
      :max="max"
    />
    <field-slider
      v-field="form.left"
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
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
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
};
</script>
