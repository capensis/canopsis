<template lang="pug">
  v-layout(column)
    v-layout(row)
      v-text-field(
        v-field="form.collection",
        :label="$t('eventFilter.collection')"
      )
    v-layout(v-for="(condition, index) in form.conditions", :key="condition.key", row)
      v-flex.pr-2(xs3)
        c-select-field(
          v-field="form.conditions[index].type",
          :items="conditionTypes"
        )
      v-flex.px-2(xs4)
        v-text-field(
          v-field="form.conditions[index].attribute",
          :label="$t('common.attribute')"
        )
      v-flex.pl-2(xs5)
        v-layout(row, align-center)
          v-text-field(
            v-field="form.conditions[index].value",
            :label="$t('common.value')"
          )
          v-btn(:disabled="hasOneCondition", icon, small, @click="removeCondition(condition.key)")
            v-icon(small) close
    v-flex
      v-btn.ml-0.mb-0(color="primary", outline, @click="addCondition") {{ $t('eventFilter.addMore') }}
</template>

<script>
import { EVENT_FILTER_EXTERNAL_DATA_TYPES, EVENT_FILTER_EXTERNAL_DATA_CONDITION_TYPES } from '@/constants';

import { eventFilterExternalDataConditionItemToForm } from '@/helpers/forms/event-filter';

import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
  },
  computed: {
    hasOneCondition() {
      return this.form.conditions.length === 1;
    },

    types() {
      return Object.values(EVENT_FILTER_EXTERNAL_DATA_TYPES)
        .map(type => ({ text: this.$t(`eventFilter.externalData.types.${type}`), value: type }));
    },

    conditionTypes() {
      return Object.values(EVENT_FILTER_EXTERNAL_DATA_CONDITION_TYPES)
        .map(type => ({ text: this.$t(`eventFilter.externalData.conditionTypes.${type}`), value: type }));
    },
  },
  methods: {
    addCondition() {
      this.updateField('conditions', [
        ...this.form.conditions,

        eventFilterExternalDataConditionItemToForm(),
      ]);
    },

    removeCondition(key) {
      this.updateField('conditions', this.form.conditions.filter(condition => condition.key !== key));
    },
  },
};
</script>
