<template lang="pug">
  v-form(@submit.prevent="$emit('submit')")
    v-list.pt-0(expand)
      slot
    v-btn.primary(
      :loading="submitting",
      :disabled="submitting || errors.any()",
      type="submit"
    ) {{ $t('common.save') }}
</template>

<script>
export default {
  inject: ['$validator'],
  props: {
    submitting: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
