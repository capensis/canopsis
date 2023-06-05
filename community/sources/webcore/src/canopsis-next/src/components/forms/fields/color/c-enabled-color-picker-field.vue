<template lang="pug">
  v-layout(column)
    c-enabled-field(
      v-model="enabled",
      :label="label",
      color="primary",
      @input="updateColor"
    )
    c-color-picker-field(
      v-field="color",
      :disabled="!enabled",
      :required="enabled",
      :name="name"
    )
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
      required: true,
    },
    label: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'color',
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
