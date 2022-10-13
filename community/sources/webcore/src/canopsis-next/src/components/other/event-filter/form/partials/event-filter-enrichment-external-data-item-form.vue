<template lang="pug">
  v-card
    v-card-text
      v-layout(column)
        v-layout(row, align-center)
          v-select(
            :value="form.type",
            :items="types",
            :label="$t('common.type')",
            @change="updateType"
          )
          v-btn.mr-0(icon, @click="remove")
            v-icon(color="error") delete
        v-text-field(
          v-field="form.reference",
          v-validate="'required'",
          :label="$t('eventFilter.reference')",
          :error-messages="errors.collect(referenceFieldName)",
          :name="referenceFieldName"
        )
          template(#append="")
            v-tooltip(left)
              template(#activator="{ bind, on }")
                v-icon(v-bind="bind", v-on="on") help
              span(v-html="$t('eventFilter.tooltips.reference')")
        event-filter-enrichment-external-data-mongo-form(v-if="isMongoType", v-field="form", :name="form.key")
        request-form(v-else, v-field="form", :name="form.key")
</template>

<script>
import { EVENT_FILTER_EXTERNAL_DATA_TYPES } from '@/constants';

import { formMixin } from '@/mixins/form';

import RequestForm from '@/components/forms/request-form.vue';

import EventFilterEnrichmentExternalDataMongoForm from './event-filter-enrichment-external-data-mongo-form.vue';
import { eventFilterExternalDataItemToForm } from '@/helpers/forms/event-filter';

export default {
  inject: ['$validator'],
  components: {
    RequestForm,
    EventFilterEnrichmentExternalDataMongoForm,
  },
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
    types() {
      return Object.values(EVENT_FILTER_EXTERNAL_DATA_TYPES)
        .map(type => ({ text: this.$t(`eventFilter.externalDataTypes.${type}`), value: type }));
    },

    isMongoType() {
      return this.form.type === EVENT_FILTER_EXTERNAL_DATA_TYPES.mongo;
    },

    referenceFieldName() {
      return `${this.name}.reference`;
    },
  },
  methods: {
    updateType(type) {
      const { reference } = this.form;

      this.updateModel(eventFilterExternalDataItemToForm(reference, { type }));
    },

    remove() {
      this.$emit('remove', this.form.key);
    },
  },
};
</script>
