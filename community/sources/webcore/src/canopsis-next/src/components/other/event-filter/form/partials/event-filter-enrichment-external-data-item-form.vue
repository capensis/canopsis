<template lang="pug">
  v-card
    v-card-text
      v-layout(column)
        v-layout(row, align-center)
          v-text-field.mr-2(
            v-field="form.reference",
            v-validate="'required'",
            :label="$t('eventFilter.reference')",
            :error-messages="errors.collect(referenceFieldName)",
            :name="referenceFieldName",
            :disabled="disabled"
          )
            template(#append="")
              c-help-icon(
                :text="$t('eventFilter.tooltips.reference')",
                icon="help",
                color="grey darken-1",
                left
              )
          v-select.ml-2(
            :value="form.type",
            :items="types",
            :label="$t('common.type')",
            :disabled="disabled",
            @change="updateType"
          )
          v-btn.mr-0(v-if="!disabled", icon, @click="remove")
            v-icon(color="error") delete
        event-filter-enrichment-external-data-mongo-form(
          v-if="isMongoType",
          v-field="form",
          :name="form.key",
          :disabled="disabled"
        )
        request-form(
          v-else,
          v-field="form",
          :name="form.key",
          :disabled="disabled"
        )
</template>

<script>
import { EVENT_FILTER_EXTERNAL_DATA_TYPES } from '@/constants';

import { formMixin } from '@/mixins/form';

import { eventFilterExternalDataItemToForm } from '@/helpers/forms/event-filter';

import RequestForm from '@/components/forms/request-form.vue';

import EventFilterEnrichmentExternalDataMongoForm from './event-filter-enrichment-external-data-mongo-form.vue';

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
    disabled: {
      type: Boolean,
      default: false,
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
      this.$emit('remove', this.form);
    },
  },
};
</script>
