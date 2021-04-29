<template lang="pug">
  v-tabs(color="secondary lighten-1", dark, centered, slider-color="primary")
    v-tab {{ $tc('common.comment', 2) }}
    v-tab-item(lazy)
      v-layout.pa-3.secondary.lighten-2
        v-flex(offset-xs2, xs8)
          remediation-instructions-comments(:remediation-instruction="remediationInstruction")
</template>

<script>
import RemediationInstructionsComments from './remediation-instructions-comments.vue';

export default {
  components: { RemediationInstructionsComments },
  props: {
    remediationInstruction: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
