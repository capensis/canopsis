<template lang="pug">
  text-editor(
    v-field="value",
    :label="label",
    :buttons="buttons",
    :extraButtons="extraButtons",
    :config="config",
    :errorMessages="errorMessages"
  )
</template>

<script>
import TextEditor from '@/components/common/text-editor/text-editor.vue';

export default {
  $_veeValidate: {
    value() {
      return this.value;
    },

    name() {
      return this.name;
    },
  },
  inject: ['$validator'],
  components: {
    TextEditor,
  },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'text',
    },
    label: {
      type: String,
      default: '',
    },
    buttons: {
      type: Array,
      default: () => [],
    },
    extraButtons: {
      type: Array,
      default: () => [],
    },
    config: {
      type: Object,
      default: () => ({}),
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
