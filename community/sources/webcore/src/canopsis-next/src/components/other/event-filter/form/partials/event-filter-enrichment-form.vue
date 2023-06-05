<template lang="pug">
  div
    c-collapse-panel.mb-2(:title="$t('externalData.title')")
      external-data-form(v-field="form.external_data")
    c-collapse-panel.mb-2(:title="$t('eventFilter.editActions')")
      event-filter-enrichment-actions-form(v-field="form.config.actions")
    v-layout(row)
      v-select.mr-3(
        v-field="form.config.on_success",
        :label="$t('eventFilter.onSuccess')",
        :items="eventFilterAfterTypes"
      )
      v-select.ml-3(
        v-field="form.config.on_failure",
        :label="$t('eventFilter.onFailure')",
        :items="eventFilterAfterTypes"
      )
    v-alert(:value="errors.has(name)", type="error") {{ $t('eventFilter.actionsRequired') }}
</template>

<script>
import { EVENT_FILTER_ENRICHMENT_AFTER_TYPES } from '@/constants';

import { formMixin } from '@/mixins/form';

import ExternalDataForm from '@/components/forms/external-data/external-data-form.vue';

import EventFilterEnrichmentActionsForm from './event-filter-enrichment-actions-form.vue';

export default {
  inject: ['$validator'],
  components: {
    ExternalDataForm,
    EventFilterEnrichmentActionsForm,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'actions',
    },
  },
  computed: {
    eventFilterAfterTypes() {
      return Object.values(EVENT_FILTER_ENRICHMENT_AFTER_TYPES);
    },
  },
  watch: {
    'form.config.actions': function validateActions() {
      this.$validator.validate(this.name);
    },
  },
  created() {
    this.attachActionsRequiredRule();
  },
  beforeDestroy() {
    this.detachActionsRequiredRule();
  },
  methods: {
    attachActionsRequiredRule() {
      this.$validator.attach({
        name: this.name,
        rules: 'required:true',
        getter: () => this.form.config.actions,
      });
    },

    detachActionsRequiredRule() {
      this.$validator.detach(this.name);
    },
  },
};
</script>
