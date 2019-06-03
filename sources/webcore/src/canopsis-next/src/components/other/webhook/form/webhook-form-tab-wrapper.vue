<template lang="pug">
  div
    v-tab
      .validation-header(:class="{ 'error--text': hasAnyError }") {{ title }}

</template>

<script>
import formValidatorErrorMixin from '@/mixins/form/validator-error';

export default {
  inject: ['$validator'],
  mixins: [formValidatorErrorMixin],
  props: {
    title: {
      type: String,
      required: true,
    },
  },
};
</script>
