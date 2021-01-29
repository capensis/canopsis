<template lang="pug">
  div
    v-textarea(
      v-field="value.message",
      :label="$t('modals.createAction.fields.message')",
      :name="`${name}.message`"
    )
    c-duration-field(v-field="value.duration", :name="`${name}.duration`", required)
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
    name: {
      type: String,
      required: 'parameters',
    },
  },
};
</script>
