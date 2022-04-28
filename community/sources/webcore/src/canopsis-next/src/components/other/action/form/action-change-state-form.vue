<template lang="pug">
  v-layout(column)
    state-criticity-field(
      v-field="value.state",
      :state-values="availableStateValues"
    )
    v-textarea.mt-4(
      v-field="value.output",
      v-validate="'required'",
      :label="$t('common.output')",
      :error-messages="errors.collect('output')",
      name="output"
    )
    action-author-field(v-field="value")
</template>

<script>
import { omit } from 'lodash';

import { ENTITIES_STATES } from '@/constants';

import { entitiesInfoMixin } from '@/mixins/entities/info';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

import ActionAuthorField from './partials/action-author-field.vue';

export default {
  inject: ['$validator'],
  components: {
    StateCriticityField,
    ActionAuthorField,
  },
  mixins: [entitiesInfoMixin],
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
  },
  computed: {
    availableStateValues() {
      return this.allowChangeSeverityToInfo ? ENTITIES_STATES : omit(ENTITIES_STATES, ['ok']);
    },
  },
};
</script>
