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

import { getWidgetTemplateModalTitle } from '@/helpers/entities/widget/template/modal-title';

import SettingsButtonField from '@/components/sidebars/form/fields/button-field.vue';

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
    widgetTemplateType: {
      type: String,
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

  emits: ['update:modelValue', 'update:template'],

  computed: {
    currentValue() {
      return this.modelValue ?? this.value;
    },

    isValueEmpty() {
      return this.defaultValue === String(this.currentValue);
    },

    modalTitle() {
      if (!this.widgetTemplateType) {
        return undefined;
      }

      return getWidgetTemplateModalTitle({
        t: this.$t.bind(this),
        te: this.$te.bind(this),
        type: this.widgetTemplateType,
        isEdit: !this.isValueEmpty,
      });
    },
  },

  methods: {
    updateModel(text, template) {
      this.$emit('update:modelValue', text);
      this.$emit('update:template', template);
    },

    showTextEditorWithTemplateModal() {
      this.$modals.show({
        name: MODALS.textEditorWithTemplate,
        dialogProps: this.dialogProps,
        config: {
          text: this.currentValue,
          template: this.template,
          templates: this.templates,
          variables: this.variables,
          title: this.modalTitle,
          action: ({ text, template }) => this.updateModel(text, template),
        },
      });
    },

    showRemoveTextConfirmationModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.updateModel(this.defaultValue, CUSTOM_WIDGET_TEMPLATE),
        },
      });
    },
  },
};
</script>
