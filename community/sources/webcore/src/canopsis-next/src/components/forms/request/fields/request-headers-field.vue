<template>
  <v-layout class="gap-3" column>
    <c-label
      :label="$tc('common.header', 2)"
      :help-text="$t('common.request.headersHelpText')"
    />

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

    <div>
      <v-btn
        color="primary"
        outlined
        @click="addItem"
      >
        {{ $t('common.request.addHeader') }}
      </v-btn>
    </div>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { CONTENT_TYPES, HEADERS } from '@/constants';

import { textPairToForm } from '@/helpers/text-pairs';

import { useArrayModelField } from '@/hooks/form/array-model-field';

import RequestHeaderField from './request-header-field.vue';

export default {
  inject: ['$validator'],
  components: { RequestHeaderField },
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
  setup(props, { emit }) {
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const headersHints = computed(() => [
      {
        text: HEADERS.authorization,
      },
      {
        text: HEADERS.contentType,
        value: Object.values(CONTENT_TYPES),
      },
    ]);

    const addItem = () => addItemIntoArray(textPairToForm());

    return {
      headersHints,
      removeItemFromArray,
      addItem,
    };
  },
};
</script>
