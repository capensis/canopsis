<template>
  <text-editor
    v-field="value"
    :label="label"
    :public="isPublic"
    :buttons="buttons"
    :config="config"
    :extra-buttons="extraButtons"
    :error-messages="errorMessages"
    :max-file-size="maxFileSize"
    :variables="variables"
    :dark="$system.dark"
  />
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
  inject: ['$validator', '$system'],
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
      required: false,
    },
    extraButtons: {
      type: Array,
      required: false,
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
    variables: {
      type: Array,
      required: false,
    },
  },
  computed: {
    isPublic() {
      return this.public;
    },
  },
};
</script>
