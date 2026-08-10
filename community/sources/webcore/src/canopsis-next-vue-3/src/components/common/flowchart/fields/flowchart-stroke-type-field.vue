<template>
  <c-layout justify-space-between align-center>
    <v-label>{{ label }}</v-label>
    <v-col cols="3">
      <v-select
        :model-value="value"
        :items="types"
        item-title="value"
        item-value="value"
        class="mt-0 pt-0"
        hide-details
        return-object
        @update:model-value="updateValue"
      >
        <template v-for="slotName in ['selection', 'item']" :key="slotName" #[slotName]="{ item }">
          <svg width="50" height="4">
            <line
              v-bind="getItemValue(item).properties"
              x1="0"
              y1="0"
              x2="100%"
              y2="0"
              stroke="currentColor"
              stroke-width="4"
            />
          </svg>
        </template>
      </v-select>
    </v-col>
  </c-layout>
</template>

<script>
import { STROKE_TYPES } from '@/constants';
import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  emits: ['input', 'update:value'],
  props: {
    value: {
      type: String,
      required: false,
    },
    label: {
      type: String,
      required: false,
    },
  },
  computed: {
    types() {
      return [
        {
          value: STROKE_TYPES.solid,
        },
        {
          value: STROKE_TYPES.dashed,
          properties: {
            'stroke-dasharray': '4 4',
          },
        },
      ];
    },
  },
  methods: {
    getItemValue(item) {
      return item?.raw ?? item?.value ?? item;
    },

    updateValue(item) {
      this.updateModel(item.value);
    },
  },
};
</script>
