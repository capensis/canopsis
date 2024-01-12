<template>
  <v-layout
    :column="column"
  >
    <v-flex
      :xs12="column"
      xs6
    >
      <c-enabled-field
        v-model="enabled"
        :label="label"
        color="primary"
        @input="updateColor"
      />
    </v-flex>
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
