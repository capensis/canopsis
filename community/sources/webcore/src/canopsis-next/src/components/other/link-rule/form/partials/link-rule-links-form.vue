<template lang="pug">
  v-layout(column)
    v-flex.mb-3(xs12)
      v-alert(:value="!links.length", type="info") {{ $t('linkRule.linksEmpty') }}
    link-rule-link-form.mb-3(
      v-for="(link, index) in links",
      v-field="links[index]",
      :key="link.key",
      :name="link.key",
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
  },
  methods: {
    addItem() {
      this.addItemIntoArray(linkRuleLinkToForm());
    },
  },
};
</script>
