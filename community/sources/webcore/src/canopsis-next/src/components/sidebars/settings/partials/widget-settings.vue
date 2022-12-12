<template lang="pug">
  v-form(@submit.prevent="$emit('submit')")
    v-list.pt-0(expand)
      slot
    v-btn.primary(
      :loading="submitting",
      :disabled="submitting",
      type="submit"
    ) {{ $t('common.save') }}
</template>

<script>
export default {
  props: {
    submitting: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
