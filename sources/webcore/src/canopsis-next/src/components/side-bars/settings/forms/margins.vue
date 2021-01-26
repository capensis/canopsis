<template lang="pug">
  v-list-group(data-test="widgetMargin")
    v-list-tile(slot="activator") {{ $t('settings.margin.title') }}
    v-list.grey.lighten-4.px-2.py-0(expand)
      field-slider(
        v-field="form.top",
        :title="$t('settings.margin.top')",
        :min="min",
        :max="max",
        data-test="widget-margin-top"
      )
      v-divider
      field-slider(
        v-field="form.right",
        :title="$t('settings.margin.right')",
        :min="min",
        :max="max",
        data-test="widget-margin-right"
      )
      v-divider
      field-slider(
        v-field="form.bottom",
        :title="$t('settings.margin.bottom')",
        :min="min",
        :max="max",
        data-test="widget-margin-bottom"
      )
      v-divider
      field-slider(
        v-field="form.left",
        :title="$t('settings.margin.left')",
        :min="min",
        :max="max",
        data-test="widget-margin-left"
      )
</template>

<script>
import FieldSlider from '../fields/common/slider.vue';

export default {
  components: { FieldSlider },
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
