<template>
  <settings-button-field
    :title="title"
    :is-empty="isValueEmpty"
    :addable="addable"
    :removable="removable"
    @create="showTextEditorModal"
    @edit="showTextEditorModal"
    @delete="showRemoveTextConfirmationModal"
  />
</template>

<script>
import SettingsButtonField from '@/components/sidebars/form/fields/button-field.vue';
import { MODALS } from '@/constants';

export default {
  components: { SettingsButtonField },

  props: {
    modelValue: {
      type: String,
      default: undefined,
    },
    value: {
      type: String,
      default: '',
    },
    title: {
      type: String,
      default: '',
    },
    variables: {
      type: Array,
      required: false,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    removable: {
      type: Boolean,
      default: false,
    },
    defaultValue: {
      type: String,
      default: '',
    },
    dialogProps: {
      type: Object,
      required: false,
    },
    sanitizeOptions: {
      type: Object,
      required: false,
    },
  },

  emits: ['update:modelValue'],

  computed: {
    currentValue() {
      return this.modelValue ?? this.value;
    },

    isValueEmpty() {
      return this.defaultValue === String(this.currentValue);
    },
  },

  methods: {
    updateModel(value) {
      this.$emit('update:modelValue', value);
    },

    showTextEditorModal() {
      this.$modals.show({
        name: MODALS.textEditor,
        dialogProps: this.dialogProps,
        config: {
          text: this.currentValue,
          variables: this.variables,
          sanitizeOptions: this.sanitizeOptions,
          action: value => this.updateModel(value),
        },
      });
    },

    showRemoveTextConfirmationModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.updateModel(this.defaultValue),
        },
      });
    },
  },
};
</script>
