<template>
  <v-layout
    justify-space-between
    align-center
  >
    <v-flex
      class="pr-2"
      xs6
    >
      <v-select
        v-field="request.method"
        v-validate="'required'"
        :items="availableMethods"
        :label="methodLabel || $t('common.method')"
        :error-messages="errors.collect(methodFieldName)"
        :name="methodFieldName"
        :disabled="disabled"
      />
    </v-flex>
    <v-flex
      class="pl-2"
      xs6
    >
      <c-payload-text-field
        v-field="request.url"
        :label="urlLabel || $t('common.url')"
        :name="urlFieldName"
        :variables="urlVariables"
        :disabled="disabled"
        :error-messages="errors.collect(urlFieldName)"
        required
      >
        <template
          v-if="helpText"
          #append=""
        >
          <c-help-icon
            :text="helpText"
            icon="help"
            color="grey darken-1"
            left
          />
        </template>
      </c-payload-text-field>
    </v-flex>
  </v-layout>
</template>

<script>
import { REQUEST_METHODS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'request',
    event: 'input',
  },
  props: {
    request: {
      type: Object,
      required: true,
    },
    methodLabel: {
      type: String,
      required: false,
    },
    urlLabel: {
      type: String,
      required: false,
    },
    helpText: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'request',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    urlVariables: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    availableMethods() {
      return Object.values(REQUEST_METHODS);
    },

    methodFieldName() {
      return `${this.name}.method`;
    },

    urlFieldName() {
      return `${this.name}.url`;
    },
  },
};
</script>
