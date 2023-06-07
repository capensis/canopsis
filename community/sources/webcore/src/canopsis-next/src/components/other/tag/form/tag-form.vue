<template lang="pug">
  v-layout(column)
    c-name-field(v-field="form.name", :disabled="isImported", required)
    c-color-picker-field(v-field="form.color")
    c-patterns-field.mt-2(
      v-if="!isImported",
      v-field="form.patterns",
      with-alarm,
      with-entity,
      some-required
    )
</template>

<script>
import CColorPickerField from '@/components/forms/fields/color/c-color-picker-field.vue';

export default {
  components: { CColorPickerField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isImported: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
