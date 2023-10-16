<template>
  <v-layout
    :row="row"
    :column="column"
  >
    <c-enabled-field
      v-model="enabled"
      :label="label"
      color="primary"
      @input="updateColor"
    />
    <c-color-picker-field
      v-field="color"
      :disabled="!enabled"
      :required="enabled"
      :name="name"
    />
  </v-layout>
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formBaseMixin],
  model: {
    prop: 'color',
    event: 'input',
  },
  props: {
    color: {
      type: String,
      required: false,
    },
    label: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'color',
    },
    row: {
      type: Boolean,
      default: false,
    },
    column: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      enabled: !!this.color,
    };
  },
  methods: {
    updateColor(value) {
      if (!value) {
        this.updateModel('');
      }
    },
  },
};
</script>
