<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.margin.title') }}
    v-list.grey.lighten-4.px-2.py-0(expand)
      field-slider(
        v-field="form.top",
        :title="$t('settings.margin.top')",
        :min="min",
        :max="max"
      )
      v-divider
      field-slider(
        v-field="form.right",
        :title="$t('settings.margin.right')",
        :min="min",
        :max="max"
      )
      v-divider
      field-slider(
        v-field="form.bottom",
        :title="$t('settings.margin.bottom')",
        :min="min",
        :max="max"
      )
      v-divider
      field-slider(
        v-field="form.left",
        :title="$t('settings.margin.left')",
        :min="min",
        :max="max"
      )
</template>

<script>
import FieldSlider from '@/components/sidebars/settings/fields/common/slider.vue';

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
