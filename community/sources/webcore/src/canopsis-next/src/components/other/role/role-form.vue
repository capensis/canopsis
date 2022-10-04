<template lang="pug">
  div
    c-name-field(v-field="form.name")
    v-layout
      v-text-field(
        v-field="form.description",
        :label="$t('common.description')",
        data-test="description"
      )
    v-layout
      view-selector(v-field="form.defaultview")
</template>

<script>
import ViewSelector from '@/components/forms/fields/view-selector.vue';

export default {
  inject: ['$validator'],
  components: { ViewSelector },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
