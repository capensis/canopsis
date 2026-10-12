<template>
  <v-layout class="gap-3" column>
    <c-label
      :label="$tc('common.link', 2)"
      :error="!!errorMessages.length"
      required
    />
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
        :color="hasErrors ? 'error' : 'primary'"
        outlined
        @click="addItem"
      >
        {{ $t('linkRule.addLink') }}
      </v-btn>
    </v-flex>
    <v-messages
      v-if="hasErrors"
      :value="errorMessages"
      color="error"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { LINK_RULE_TYPES } from '@/constants';

import { linkRuleLinkToForm } from '@/helpers/entities/link/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { formValidationHeaderMixin } from '@/mixins/form/validation-header';

import LinkRuleLinkForm from './fields/link-rule-link-form.vue';

export default {
  inject: ['$validator'],
  components: { LinkRuleLinkForm },
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
    templateVars: {
      type: Object,
      default: () => ({}),
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const hasErrors = computed(() => !!props.errorMessages.length);

    const addItem = () => addItemIntoArray(linkRuleLinkToForm());

    return {
      hasErrors,

      addItem,
      removeItemFromArray,
    };
  },
};
</script>
