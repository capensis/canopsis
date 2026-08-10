<template>
  <c-layout gap="3" column>
    <c-label :label="$t('common.link', 2)" :error="!!errorMessages.length" required />
    <link-rule-link-form
      v-for="(link, index) in linksField"
      :key="link.key"
      :form="linksField[index]"
      :name="link.key"
      :type="type"
      :template-vars="templateVars"
      class="mb-3"
      @update:form="updateItem(index, $event)"
      @change="$emit('change')"
      @remove="removeItem(index)"
    />
    <div>
      <v-btn :color="hasErrors ? 'error' : 'primary'" variant="outlined" @click="addItem">
        {{ $t('linkRule.addLink') }}
      </v-btn>
    </div>
    <v-messages v-if="hasErrors" :messages="errorMessages" color="error" />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { LINK_RULE_TYPES } from '@/constants';
import { linkRuleLinkToForm } from '@/helpers/entities/link/form';
import { useFormArray, useFormContext } from '@/hooks/form';
import { useValidationHeader } from '@/hooks/validator/validation-header';

import LinkRuleLinkForm from './fields/link-rule-link-form.vue';

export default {
  components: { LinkRuleLinkForm },
  props: {
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

  emits: ['change'],

  setup(props, { emit }) {
    const form = useFormContext();
    const linksField = useFormArray('links');
    const { hasAnyError } = useValidationHeader();

    const hasErrors = computed(() => !!props.errorMessages.length);

    const addItem = () => {
      form.push('links', linkRuleLinkToForm());
      emit('change');
    };
    const updateItem = (index, value) => {
      form.replace('links', index, value);
      emit('change');
    };
    const removeItem = (index) => {
      form.remove('links', index);
      emit('change');
    };

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs
       */
      hasAnyError,
      linksField,
      hasErrors,
      addItem,
      updateItem,
      removeItem,
    };
  },
};
</script>
