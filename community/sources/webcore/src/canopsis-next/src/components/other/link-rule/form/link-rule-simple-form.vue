<template lang="pug">
  link-rule-links-form(v-field="links")
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form/validation-header';

import LinkRuleLinksForm from './partials/link-rule-links-form.vue';

export default {
  inject: ['$validator'],
  components: { LinkRuleLinksForm },
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'links',
    event: 'input',
  },
  props: {
    links: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
