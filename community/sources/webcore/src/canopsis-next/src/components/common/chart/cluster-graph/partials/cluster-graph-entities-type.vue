<template lang="pug">
  div
    h6.my-2.title.text-xs-center {{ $t('common.type') }}
    v-radio-group.justify-center(v-field="impact", row, hide-details)
      v-radio(
        :value="false",
        :label="$t('context.treeOfDependencies')",
        color="primary"
      )
      v-radio(
        :value="true",
        :label="$t('context.impactChain')",
        color="primary"
      )
</template>

<script>
export default {
  model: {
    prop: 'impact',
    event: 'change',
  },
  props: {
    impact: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
