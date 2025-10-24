<template>
  <v-layout column>
    <c-alert
      :value="!links.length && !errors.has('links')"
      type="info"
    >
      {{ $t('linkRule.linksEmpty') }}
    </c-alert>
    <link-rule-link-form
      v-for="(link, index) in links"
      v-field="links[index]"
      :key="link.key"
      :name="link.key"
      :type="type"
      :template-vars="templateVars"
      class="mb-3"
      @remove="removeItemFromArray(index)"
    />
    <v-flex>
      <v-btn
        class="ml-0 my-0"
        color="primary"
        outlined
        @click="addItem"
      >
        {{ $t('linkRule.addLink') }}
      </v-btn>
    </v-flex>
  </v-layout>
</template>

<script>
import { LINK_RULE_TYPES } from '@/constants';

import { linkRuleLinkToForm } from '@/helpers/entities/link/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';

import LinkRuleLinkForm from './link-rule-link-form.vue';

export default {
  inject: ['$validator'],
  components: { LinkRuleLinkForm },
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
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const addItem = () => addItemIntoArray(linkRuleLinkToForm());

    return {
      addItem,
      removeItemFromArray,
    };
  },
};
</script>
