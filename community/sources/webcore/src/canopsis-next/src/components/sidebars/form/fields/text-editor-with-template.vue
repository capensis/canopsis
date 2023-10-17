<template>
  <settings-button-field
    :title="title"
    :is-empty="isValueEmpty"
    :addable="addable"
    :removable="removable"
    @create="showTextEditorWithTemplateModal"
    @edit="showTextEditorWithTemplateModal"
    @delete="showRemoveTextConfirmationModal"
  />
</template>

<script>
import { CUSTOM_WIDGET_TEMPLATE, MODALS } from '@/constants';

import SettingsButtonField from '@/components/sidebars/form/fields/button-field.vue';

export default {
  components: { SettingsButtonField },
  props: {
    value: {
      type: String,
      default: '',
    },
    title: {
      type: String,
      default: '',
    },
    template: {
      type: String,
      required: false,
    },
    templates: {
      type: Array,
      required: false,
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
  },
  computed: {
    isValueEmpty() {
      return this.defaultValue === String(this.value);
    },
  },
  methods: {
    showTextEditorWithTemplateModal() {
      this.$modals.show({
        name: MODALS.textEditorWithTemplate,
        dialogProps: this.dialogProps,
        config: {
          text: this.value,
          template: this.template,
          templates: this.templates,
          variables: this.variables,
          action: ({ text, template }) => this.$emit('input', text, template),
        },
      });
    },

    showRemoveTextConfirmationModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.$emit('input', this.defaultValue, CUSTOM_WIDGET_TEMPLATE),
        },
      });
    },
  },
};
</script>
