<template lang="pug">
  v-layout(column)
    c-enabled-field(
      v-field="value.forward_author",
      :label="$t('scenario.forwardAuthor')"
    )
    v-text-field(
      v-if="!value.forward_author",
      v-field="value.author",
      v-validate="'required'",
      :error-messages="errors.collect('author')",
      :label="$t('common.author')",
      name="author"
    )
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
};
</script>
