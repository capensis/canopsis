<template lang="pug">
  div
    v-textarea(
      v-field="form.comment",
      v-validate="'required'",
      :label="$t('modals.createPause.comment')",
      :error-messages="errors.collect('comment')",
      name="comment"
    )
    c-pbehavior-reason-field(v-field="form.reason", required, return-object)
</template>

<script>
export default {
  inject: ['$validator'],
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
