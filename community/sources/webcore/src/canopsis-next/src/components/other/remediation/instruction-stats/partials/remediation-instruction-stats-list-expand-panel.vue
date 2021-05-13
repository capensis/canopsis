<template lang="pug">
  pre {{ remediationInstruction | json }}
</template>

<script>
export default {
  props: {
    remediationInstruction: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
