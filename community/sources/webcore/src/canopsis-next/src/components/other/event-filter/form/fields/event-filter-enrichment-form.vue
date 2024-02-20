<template>
  <v-layout column>
    <c-collapse-panel
      class="mb-2"
      :title="$t('eventFilter.editActions')"
    >
      <event-filter-enrichment-actions-form
        v-field="form.config.actions"
        :variables="templateVariables"
        :name="name"
        :set-tags-items="setTagsItems"
      />
    </c-collapse-panel>
    <v-layout>
      <v-select
        class="mr-3"
        v-field="form.config.on_success"
        :label="$t('eventFilter.onSuccess')"
        :items="eventFilterAfterTypes"
      />
      <v-select
        class="ml-3"
        v-field="form.config.on_failure"
        :label="$t('eventFilter.onFailure')"
        :items="eventFilterAfterTypes"
      />
    </v-layout>
    <c-alert
      :value="errors.has(name)"
      type="error"
    >
      {{ $t('eventFilter.actionsRequired') }}
    </c-alert>
  </v-layout>
</template>

<script>
import { EVENT_FILTER_ENRICHMENT_AFTER_TYPES } from '@/constants';

import { formMixin } from '@/mixins/form';
import { validationAttachRequiredMixin } from '@/mixins/form/validation-attach-required';

import EventFilterEnrichmentActionsForm from './event-filter-enrichment-actions-form.vue';

export default {
  inject: ['$validator'],
  components: {
    EventFilterEnrichmentActionsForm,
  },
  mixins: [formMixin, validationAttachRequiredMixin],
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
      default: 'config.actions',
    },
    templateVariables: {
      type: Array,
      default: () => [],
    },
    setTagsItems: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    eventFilterAfterTypes() {
      return Object.values(EVENT_FILTER_ENRICHMENT_AFTER_TYPES);
    },
  },
  watch: {
    'form.config.actions': function validateActions() {
      this.validateRequiredRule();
    },
  },
  created() {
    this.attachRequiredRule(() => this.form.config.actions);
  },
  beforeDestroy() {
    this.detachRequiredRule();
  },
};
</script>
