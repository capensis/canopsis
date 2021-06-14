<template lang="pug">
  v-tabs(color="secondary lighten-1", dark, centered, slider-color="primary")
    v-tab {{ $tc('common.comment', remediationInstruction.comments.length) }}
    v-tab-item
      v-layout.pa-3.secondary.lighten-2
        v-flex(offset-xs2, xs8)
          remediation-instructions-comments(:comments="remediationInstruction.comments")
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
