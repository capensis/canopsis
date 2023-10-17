<template>
  <v-layout
    justify-space-between
    align-center
  >
    <v-label>{{ label }}</v-label>
    <v-flex xs3>
      <v-select
        class="mt-0 pt-0"
        v-field="value"
        :items="types"
        hide-details
      >
        <template
          v-for="slotName in ['selection', 'item']"
          #[slotName]="{ item }"
        >
          <svg
            :key="slotName"
            width="50"
            height="4"
          >
            <line
              v-bind="item.properties"
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
    </v-flex>
  </v-layout>
</template>

<script>
import { STROKE_TYPES } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
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
};
</script>
