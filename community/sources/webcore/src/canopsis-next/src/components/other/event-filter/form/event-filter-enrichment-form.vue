<template lang="pug">
  div
    c-pattern-panel.mb-2(:title="$t('eventFilter.externalData')")
      event-filter-enrichment-external-data-form(v-field="form.external_data")
    c-pattern-panel.mb-2(:title="$t('eventFilter.editActions')")
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
import { EVENT_FILTER_ENRICHMENT_AFTER_TYPES, MODALS } from '@/constants';

import { formMixin } from '@/mixins/form';

import EventFilterEnrichmentActionsForm from './event-filter-enrichment-actions-form.vue';
import EventFilterEnrichmentExternalDataForm from './event-filter-enrichment-external-data-form.vue';

export default {
  inject: ['$validator'],
  components: {
    EventFilterEnrichmentActionsForm,
    EventFilterEnrichmentExternalDataForm,
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
  created() {
    this.attachActionsRequiredRule();
  },
  beforeDestroy() {
    this.detachActionsRequiredRule();
  },
  methods: {
    showEditActionsModal() {
      this.$modals.show({
        name: MODALS.eventFilterActions,
        config: {
          actions: this.form.config.actions,
          action: (updatedActions) => {
            this.updateField('actions', updatedActions);
            this.$nextTick(() => this.$validator.validate('actions'));
          },
        },
      });
    },

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
