<template lang="pug">
  v-tabs(slider-color="primary", fixed-tabs)
    v-tab(:disabled="patternWasChanged") {{ $t('modals.eventFilterRule.simpleEditor') }}
    v-tab {{ $t('modals.eventFilterRule.advancedEditor') }}
    v-tab-item
      pattern-simple-form(
        v-field="form",
        :operators="operators",
        :only-simple-rule="onlySimpleRule"
      )
    v-tab-item
      json-field(
        v-field="form",
        name="pattern",
        rows="15",
        validate-on="button"
      )
</template>

<script>
import { get } from 'lodash';

import { EVENT_FILTER_RULE_OPERATORS } from '@/constants';

import JsonField from '@/components/forms/fields/json-field.vue';

import PatternSimpleForm from './pattern-simple-form.vue';

export default {
  inject: ['$validator'],
  components: { PatternSimpleForm, JsonField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    operators: {
      type: Array,
      default: () => EVENT_FILTER_RULE_OPERATORS,
    },
    onlySimpleRule: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    patternWasChanged() {
      return get(this.fields, ['pattern', 'changed']);
    },
  },
};
</script>
