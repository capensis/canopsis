<template>
  <v-layout
    :class="{ 'text-pairs__disabled': disabled }"
    class="text-pairs"
    wrap
  >
    <v-flex
      v-show="title"
      xs12
    >
      <h4 class="ml-1">
        {{ title }}
      </h4>
    </v-flex>
    <v-flex xs12>
      <slot
        v-if="!items.length"
        name="no-data"
      />
      <c-text-pair-field
        v-for="(item, index) in items"
        v-field="items[index]"
        :key="item[itemKey]"
        :disabled="disabled"
        :value-required="valueRequired"
        :text-required="textRequired"
        :text-label="textLabel"
        :value-label="valueLabel"
        :item-text="itemText"
        :item-value="itemValue"
        :name="item[itemKey]"
        @remove="removeItemFromArray(index)"
      >
        <template #append-value="">
          <slot
            :item="item"
            name="append-value"
          />
        </template>
      </c-text-pair-field>
    </v-flex>
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
          {{ addButtonLabel || $t('common.add') }}
        </v-btn>
      </v-layout>
    </v-flex>
  </v-layout>
</template>

<script>
import { textPairToForm } from '@/helpers/text-pairs';

import { formArrayMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formArrayMixin],
  model: {
    prop: 'items',
    event: 'input',
  },
  props: {
    title: {
      type: String,
      default: null,
    },
    items: {
      type: Array,
      default: () => [],
    },
    textLabel: {
      type: String,
      default: '',
    },
    valueLabel: {
      type: String,
      default: '',
    },
    itemText: {
      type: String,
      required: false,
    },
    itemValue: {
      type: String,
      required: false,
    },
    itemKey: {
      type: String,
      default: 'key',
    },
    name: {
      type: String,
      default: 'items',
    },
    textRequired: {
      type: Boolean,
      default: false,
    },
    valueRequired: {
      type: Boolean,
      default: false,
    },
    addButtonLabel: {
      type: String,
      required: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    addItem() {
      this.addItemIntoArray(textPairToForm());
    },
  },
};
</script>
