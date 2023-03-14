<template lang="pug">
  v-card
    v-card-text
      v-layout(column)
        v-layout(row, align-center)
          v-text-field.mr-2(
            v-field="form.reference",
            v-validate="'required'",
            :label="$t('externalData.fields.reference')",
            :error-messages="errors.collect(referenceFieldName)",
            :name="referenceFieldName",
            :disabled="disabled"
          )
            template(#append="")
              c-help-icon(
                :text="$t('externalData.tooltips.reference')",
                icon="help",
                left
              )
          v-select.ml-2(
            v-field="form.type",
            :items="types",
            :label="$t('common.type')",
            :disabled="disabled"
          )
          v-btn.mr-0(v-if="!disabled", icon, @click="remove")
            v-icon(color="error") delete
        external-data-mongo-form(
          v-if="isMongoType",
          v-field="form",
          :name="form.key",
          :disabled="disabled",
          :variables="variables"
        )
        request-form(
          v-else,
          v-field="form.request",
          :name="form.key",
          :disabled="disabled"
        )
</template>

<script>
import { EXTERNAL_DATA_TYPES } from '@/constants';

import { formMixin } from '@/mixins/form';

import RequestForm from '@/components/forms/request/request-form.vue';

import ExternalDataMongoForm from './external-data-mongo-form.vue';

export default {
  inject: ['$validator'],
  components: { RequestForm, ExternalDataMongoForm },
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
    types: {
      type: Array,
      default: () => [],
    },
    variables: {
      type: Array,
      default: () => ([]),
    },
  },
  computed: {
    availableTypes() {
      return this.types.length
        ? this.types
        : Object.values(EXTERNAL_DATA_TYPES)
          .map(type => ({ text: this.$t(`externalData.types.${type}`), value: type }));
    },

    isMongoType() {
      return this.form.type === EXTERNAL_DATA_TYPES.mongo;
    },

    referenceFieldName() {
      return `${this.name}.reference`;
    },
  },
  methods: {
    remove() {
      this.$emit('remove', this.form);
    },
  },
};
</script>
