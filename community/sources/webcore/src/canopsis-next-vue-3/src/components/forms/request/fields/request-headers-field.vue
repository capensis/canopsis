<template>
  <c-layout gap="3" column>
    <c-label :label="$t('common.header', 2)" :help-text="$t('common.request.headersHelpText')" />

    <c-layout
      v-for="(item, index) in headers"
      :key="item.key"
      align="center"
      justify="space-between"
    >
      <request-header-field
        :header="headers[index]"
        :disabled="disabled"
        :name="`${name}.${item.key}`"
        :headers-hints="headersHints"
        :headers-variables="headersVariables"
        @update:header="updateItemInArray(index, $event)"
      />
      <c-action-btn v-if="!disabled" type="delete" @click="removeItemFromArray(index)" />
    </c-layout>

    <div>
      <v-btn color="primary" outlined @click="addItem">
        {{ $t('common.request.addHeader') }}
      </v-btn>
    </div>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { CONTENT_TYPES, HEADERS } from '@/constants';
import { textPairToForm } from '@/helpers/text-pairs';
import { useArrayModelField } from '@/hooks/form/array-model-field';

import RequestHeaderField from './request-header-field.vue';

export default {
  emits: ['update:headers'],
  components: { RequestHeaderField },
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
    const { addItemIntoArray, removeItemFromArray, updateItemInArray } = useArrayModelField(
      props,
      emit,
    );

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
      updateItemInArray,
      addItem,
    };
  },
};
</script>
