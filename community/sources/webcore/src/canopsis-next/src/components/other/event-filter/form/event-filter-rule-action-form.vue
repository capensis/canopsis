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
      v-if="showingName",
      v-field="form.name",
      v-validate="'required'",
      :label="$t('common.name')",
      :error-messages="errors.collect('name')",
      name="name",
      key="name"
    )
    v-text-field(
      v-if="showingDescription",
      v-field="form.description",
      :label="$t('common.description')",
      key="description"
    )
    c-mixed-field(
      v-if="showingValue",
      v-field="form.value",
      :label="$t('common.value')",
      key="value"
    )
    v-text-field(
      v-if="showingFrom",
      v-field="form.value",
      v-validate="'required'",
      :label="$t('common.from')",
      :error-messages="errors.collect('value')",
      key="from",
      name="value"
    )
      v-tooltip(slot="append", left)
        v-icon(slot="activator") help
        div(v-html="$t('eventFilter.tooltips.copyFromHelp')")
    v-text-field(
      v-if="showingTo",
      v-field="form.name",
      v-validate="'required'",
      :label="$t('common.to')",
      :error-messages="errors.collect('name')",
      key="to",
      name="name"
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

    showingName() {
      return this.form.type !== EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy;
    },

    showingDescription() {
      return [
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfo,
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate,
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copyToEntityInfo,
      ].includes(this.form.type);
    },

    showingValue() {
      return [
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField,
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate,
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfo,
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate,
      ].includes(this.form.type);
    },

    showingFrom() {
      return [
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy,
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copyToEntityInfo,
      ].includes(this.form.type);
    },

    showingTo() {
      return EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy === this.form.type;
    },
  },
  watch: {
    'form.type': function typeWatcher() {
      this.errors.clear();
    },
  },
};
</script>
