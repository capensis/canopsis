<template lang="pug">
  v-layout(column)
    v-text-field(
      v-field="form.name",
      v-validate="'required'",
      :label="$t('common.name')",
      :error-messages="errors.collect('name')",
      name="name"
    )
    c-enabled-field(v-field="form.enabled")
</template>

<script>
export default {
  inject: ['$validator'],
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
