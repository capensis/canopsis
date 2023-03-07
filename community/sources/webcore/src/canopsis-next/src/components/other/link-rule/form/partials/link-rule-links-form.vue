<template lang="pug">
  v-layout(column)
    c-alert(:value="!links.length && !errors.has('links')", type="info") {{ $t('linkRule.linksEmpty') }}
    link-rule-link-form.mb-3(
      v-for="(link, index) in links",
      v-field="links[index]",
      :key="link.key",
      :name="link.key",
      :type="type",
      @remove="removeItemFromArray(index)"
    )
    v-flex
      v-btn.ml-0.my-0(
        color="primary",
        outline,
        @click="addItem"
      ) {{ $t('linkRule.addLink') }}
</template>

<script>
import { LINK_RULE_TYPES } from '@/constants';

import { linkRuleLinkToForm } from '@/helpers/forms/link-rule';

import { formArrayMixin } from '@/mixins/form';

import LinkRuleLinkForm from './link-rule-link-form.vue';

export default {
  inject: ['$validator'],
  components: { LinkRuleLinkForm },
  mixins: [formArrayMixin],
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
  methods: {
    addItem() {
      this.addItemIntoArray(linkRuleLinkToForm());
    },
  },
};
</script>
