<template>
  <div>
    <v-layout
      v-for="(value, index) in values"
      :key="index"
      align-center
    >
      <v-flex>
        <v-text-field
          v-field="values[index]"
          :disabled="disabled"
          :label="$t('common.value')"
        />
      </v-flex>
      <c-action-btn
        v-if="!disabled"
        type="delete"
        @click="removeItemFromArray(index)"
      />
    </v-layout>
    <v-messages
      :value="errorMessages"
      color="error"
    />
    <v-btn
      class="v-btn-legacy-m--y"
      :disabled="disabled"
      color="primary"
      outlined
      @click="addItem"
    >
      {{ $t('common.add') }}
    </v-btn>
  </div>
</template>

<script>
import { formArrayMixin } from '@/mixins/form/array';

export default {
  inject: ['$validator'],
  mixins: [
    formArrayMixin,
  ],
  model: {
    prop: 'values',
    event: 'change',
  },
  props: {
    values: {
      type: Array,
      default: () => [],
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    addItem() {
      this.addItemIntoArray('');
    },
  },
};
</script>
