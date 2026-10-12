<template>
  <slot />
</template>

<script>
import { FormContextKey as VEE_FORM_CONTEXT_KEY } from 'vee-validate';
import { provide } from 'vue';

import { FORM_CONTEXT_KEY } from '@/constants';

export default {
  props: {
    formContext: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    provide(FORM_CONTEXT_KEY, props.formContext);
    provide(VEE_FORM_CONTEXT_KEY, props.formContext.veeFormContext);
  },
};
</script>
