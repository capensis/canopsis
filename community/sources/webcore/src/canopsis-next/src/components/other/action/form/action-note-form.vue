<template lang="pug">
  v-layout(row)
    v-textarea(v-field="value.output", :label="$t('scenario.output')")
      v-tooltip(slot="append", left)
        v-icon(slot="activator") help
        div(v-html="$t('scenario.outputHelp')")
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
