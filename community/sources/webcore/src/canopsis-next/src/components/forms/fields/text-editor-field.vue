<template lang="pug">
  text-editor(
    v-field="value",
    :label="label",
    :public="public",
    :buttons="buttons",
    :config="config",
    :extra-buttons="extraButtons",
    :error-messages="errorMessages",
    :max-file-size="maxFileSize"
  )
</template>

<script>
import { entitiesInfoMixin } from '@/mixins/entities/info';

const TextEditor = () => import(/* webpackChunkName: "TextEditor" */ '@/components/common/text-editor/text-editor.vue');

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
  mixins: [entitiesInfoMixin],
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
    public: {
      type: Boolean,
      default: false,
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
    maxFileSize: {
      type: Number,
      default() {
        return this.fileUploadMaxSize;
      },
    },
  },
};
</script>
