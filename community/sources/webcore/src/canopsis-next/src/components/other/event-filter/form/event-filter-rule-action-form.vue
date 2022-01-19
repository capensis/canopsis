<template lang="pug">
  div
    v-select(
      v-model="form.type",
      :items="eventFilterRuleActionTypes",
      :label="$t('common.type')"
    )
    v-expand-transition
      event-filter-rule-action-form-type-info(v-if="form.type", :type="form.type")
    v-text-field(
      v-field="form.description",
      :label="$t('common.description')",
      key="description"
    )
    v-text-field(
      v-field="form.name",
      v-validate="'required'",
      :label="$t('common.name')",
      :error-messages="errors.collect('name')",
      name="name",
      key="name"
    )
    v-text-field(
      v-if="isCopyActionType",
      v-field="form.value",
      v-validate="'required'",
      :label="$t('common.value')",
      :error-messages="errors.collect('value')",
      key="from",
      name="value"
    )
      v-tooltip(slot="append", left)
        v-icon(slot="activator") help
        div(v-html="$t('eventFilter.tooltips.copyFromHelp')")
    c-mixed-field(
      v-else,
      v-field="form.value",
      :label="$t('common.value')",
      key="value"
    )
</template>

<script>
import { EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES } from '@/constants';

import EventFilterRuleActionFormTypeInfo from './partials/event-filter-rule-action-form-type-info.vue';

export default {
  inject: ['$validator'],
  components: { EventFilterRuleActionFormTypeInfo },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  computed: {
    eventFilterRuleActionTypes() {
      return Object.values(EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES).map(value => ({
        value,

        text: this.$t(`eventFilter.actionsTypes.${value}.text`),
      }));
    },

    isCopyActionType() {
      return this.form.type === EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy;
    },
  },
  watch: {
    'form.type': function typeWatcher() {
      this.errors.clear();
    },
  },
};
</script>
