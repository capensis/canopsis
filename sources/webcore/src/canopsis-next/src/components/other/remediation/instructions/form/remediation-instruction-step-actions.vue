<template lang="pug">
  v-layout(row)
    v-btn.accent.darken-1(@click="$emit('add-operation')") {{ $t('remediationInstructions.addOperation') }}
    v-btn.accent.darken-1(
      v-if="hasOperations",
      @click="$emit('add-endpoint')"
    ) {{ $t('remediationInstructions.addEndpoint') }}
</template>

<script>
export default {
  props: {
    addDisabled: {
      type: Boolean,
      default: false,
    },
    hasOperations: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
