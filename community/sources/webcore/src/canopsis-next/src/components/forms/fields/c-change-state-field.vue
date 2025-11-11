<template>
  <div class="mt-3">
    <v-layout>
      <state-criticity-field
        v-field="value.state"
        :state-values="availableStateValues"
        mandatory
      />
    </v-layout>
    <v-layout class="mt-4">
      <component
        :is="textareaComponent"
        v-field="value.output"
        v-validate="'required'"
        :label="label || $t('common.note')"
        :error-messages="errors.collect(outputFieldName)"
        :name="outputFieldName"
        :variables="variables"
        autofocus
      />
    </v-layout>
  </div>
</template>

<script>
import { omit } from 'lodash';
import { computed } from 'vue';

import { ALARM_STATES } from '@/constants';

import { useInfo } from '@/hooks/store/modules/info';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

export default {
  inject: ['$validator'],
  components: { StateCriticityField },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'changeState',
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const { allowChangeSeverityToInfo } = useInfo();

    const outputFieldName = computed(() => `${props.name}.output`);
    const availableStateValues = computed(() => (
      allowChangeSeverityToInfo.value ? ALARM_STATES : omit(ALARM_STATES, ['ok'])
    ));

    const textareaComponent = computed(() => (
      props.variables?.length ? 'c-payload-textarea-field' : 'v-textarea'
    ));

    return {
      outputFieldName,
      availableStateValues,
      textareaComponent,
    };
  },
};
</script>
