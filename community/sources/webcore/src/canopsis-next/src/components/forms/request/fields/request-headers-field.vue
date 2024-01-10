<template>
  <v-layout column>
    <c-alert
      v-if="!headers.length"
      type="info"
    >
      {{ $t('common.request.emptyHeaders') }}
    </c-alert>
    <v-layout
      v-for="(item, index) in headers"
      :key="item.key"
      align-center
      justify-space-between
    >
      <request-header-field
        v-field="headers[index]"
        :disabled="disabled"
        :name="`${name}.${item.key}`"
        :headers-hints="headersHints"
        :headers-variables="headersVariables"
      />
      <c-action-btn
        v-if="!disabled"
        type="delete"
        @click="removeItemFromArray(index)"
      />
    </v-layout>
    <v-flex
      v-if="!disabled"
      xs12
    >
      <v-layout>
        <v-btn
          class="ml-0"
          color="primary"
          outlined
          @click="addItem"
        >
          {{ $t('common.add') }}
        </v-btn>
      </v-layout>
    </v-flex>
  </v-layout>
</template>

<script>
import { CONTENT_TYPES, HEADERS } from '@/constants';

import { textPairToForm } from '@/helpers/text-pairs';

import { formArrayMixin } from '@/mixins/form';

import RequestHeaderField from './request-header-field.vue';

export default {
  inject: ['$validator'],
  components: { RequestHeaderField },
  mixins: [formArrayMixin],
  model: {
    prop: 'headers',
    event: 'input',
  },
  props: {
    title: {
      type: String,
      default: null,
    },
    headers: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'headers',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    headersVariables: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    headersHints() {
      return [
        {
          text: HEADERS.authorization,
        },
        {
          text: HEADERS.contentType,
          value: Object.values(CONTENT_TYPES),
        },
      ];
    },
  },
  methods: {
    addItem() {
      this.addItemIntoArray(textPairToForm());
    },
  },
};
</script>
