<template lang="pug">
  v-layout(column)
    v-textarea(v-field="value.output", :label="$t('scenario.output')")
      template(#append="")
        c-help-icon(:text="$t('scenario.outputHelp')", icon="help", color="grey darken-1", left)
    action-author-field(v-field="value")
</template>

<script>
import ActionAuthorField from './partials/action-author-field.vue';

export default {
  inject: ['$validator'],
  components: { ActionAuthorField },
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
