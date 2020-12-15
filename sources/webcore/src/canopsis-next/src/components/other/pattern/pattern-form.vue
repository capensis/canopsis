<template lang="pug">
  v-tabs(slider-color="primary", fixed-tabs)
    v-tab(:disabled="patternWasChanged") {{ $t('modals.eventFilterRule.simpleEditor') }}
    v-tab {{ $t('modals.eventFilterRule.advancedEditor') }}
    v-tab-item
      pattern-simple-editor(
        v-field="form",
        :operators="operators",
        :is-simple-pattern="isSimplePattern"
      )
    v-tab-item
      json-field(
        v-field="form",
        name="pattern",
        rows="15"
      )
</template>

<script>
import { get } from 'lodash';

import { EVENT_FILTER_RULE_OPERATORS } from '@/constants';

import PatternSimpleEditor from '@/components/modals/event-filter/pattern/partial/pattern-simple-editor.vue';
import JsonField from '@/components/forms/fields/json-field.vue';

export default {
  inject: ['$validator'],
  components: { PatternSimpleEditor, JsonField },
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
    isSimplePattern: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    patternWasChanged() {
      return get(this.fields, ['pattern', 'touched']);
    },
  },
};
</script>
