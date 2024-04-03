<template>
  <v-card>
    <v-card-text>
      <v-layout column>
        <v-layout align-center>
          <v-text-field
            v-field="form.reference"
            v-validate="'required'"
            :label="$t('externalData.fields.reference')"
            :error-messages="errors.collect(referenceFieldName)"
            :name="referenceFieldName"
            :disabled="disabled"
            class="mr-2"
          >
            <template #append="">
              <c-help-icon
                :text="$t('externalData.tooltips.reference')"
                icon="help"
                left
              />
            </template>
          </v-text-field>
          <v-select
            v-field="form.type"
            :items="availableTypes"
            :label="$t('common.type')"
            :disabled="disabled"
            class="ml-2"
          />
          <v-btn
            v-if="!disabled"
            class="mr-0"
            icon
            @click="remove"
          >
            <v-icon color="error">
              delete
            </v-icon>
          </v-btn>
        </v-layout>
        <external-data-mongo-form
          v-if="isMongoType"
          v-field="form"
          :name="name"
          :disabled="disabled"
          :variables="variables"
        />
        <request-form
          v-else
          v-field="form.request"
          :name="`${name}.request`"
          :disabled="disabled"
          :payload-variables="variables"
          :url-variables="variables"
        />
      </v-layout>
    </v-card-text>
  </v-card>
</template>

<script>
import { EXTERNAL_DATA_TYPES } from '@/constants';

import { isMongoExternalDataType } from '@/helpers/entities/shared/external-data/entity';

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
      return isMongoExternalDataType(this.form.type);
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
