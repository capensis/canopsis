<template>
  <link-rule-links-form
    v-field="links"
    :type="type"
  />
</template>

<script>
import { LINK_RULE_TYPES } from '@/constants';

import { formValidationHeaderMixin } from '@/mixins/form/validation-header';

import LinkRuleLinksForm from './fields/link-rule-links-form.vue';

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
      type: Array,
      default: () => [],
    },
    type: {
      type: String,
      default: LINK_RULE_TYPES.alarm,
    },
  },
};
</script>
