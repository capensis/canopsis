<template lang="pug">
  v-layout(row)
    v-btn.ml-0.primary(
      :disabled="addDisabled",
      @click="$emit('add-step')"
    ) {{ $t('remediationInstructions.addStep') }}
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
