<template lang="pug">
  div
    c-id-field(
      v-field="form._id",
      :disabled="isDisabledIdField",
      :help-text="$t('eventFilter.idHelp')"
    )
    c-event-filter-type-field(v-field="form.type")
    c-description-field(v-field="form.description", required)
    c-priority-field(v-field="form.priority")
    c-enabled-field(v-field="form.enabled")
    c-patterns-field(v-field="form.patterns", with-entity, with-event, some-required)

    template(v-if="isChangeEntityType || isEnrichmentType")
      v-divider.my-3
      c-information-block(:title="$t('eventFilter.configuration')")
        template(v-if="isChangeEntityType")
          event-filter-change-entity-form(v-field="form.config")
        template(v-if="isEnrichmentType")
          event-filter-enrichment-form(v-field="form.config")
</template>

<script>
import { EVENT_FILTER_TYPES } from '@/constants';

import EventFilterEnrichmentForm from '@/components/other/event-filter/form/event-filter-enrichment-form.vue';
import EventFilterChangeEntityForm from '@/components/other/event-filter/form/event-filter-change-entity-form.vue';

export default {
  inject: ['$validator'],
  components: { EventFilterEnrichmentForm, EventFilterChangeEntityForm },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isEnrichmentType() {
      return this.form.type === EVENT_FILTER_TYPES.enrichment;
    },

    isChangeEntityType() {
      return this.form.type === EVENT_FILTER_TYPES.changeEntity;
    },
  },
};
</script>
