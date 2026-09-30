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
import { computed } from 'vue';

import { CUSTOM_WIDGET_TEMPLATE, MODALS } from '@/constants';

import { getWidgetTemplateModalTitle } from '@/helpers/entities/widget/template/modal-title';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import SettingsButtonField from '@/components/sidebars/form/fields/button-field.vue';

export default {
  components: { SettingsButtonField },

  props: {
    modelValue: {
      type: String,
      default: undefined,
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

  setup(props, { emit }) {
    const { t, te } = useI18n();
    const modals = useModals();
    const { updateModel } = useModelField(props, emit);
    const isValueEmpty = computed(() => props.defaultValue === String(props.modelValue));
    const modalTitle = computed(() => {
      if (!props.widgetTemplateType) {
        return undefined;
      }

      return getWidgetTemplateModalTitle({
        t,
        te,
        type: props.widgetTemplateType,
        isEdit: !isValueEmpty.value,
      });
    });

    const updateText = (text, template) => {
      updateModel(text);
      emit('update:template', template, text);
    };

    const showTextEditorWithTemplateModal = () => {
      modals.show({
        name: MODALS.textEditorWithTemplate,
        dialogProps: props.dialogProps,
        config: {
          text: props.modelValue,
          template: props.template,
          templates: props.templates,
          variables: props.variables,
          title: modalTitle.value,
          action: ({ text, template }) => updateText(text, template),
        },
      });
    };

    const showRemoveTextConfirmationModal = () => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => updateText(props.defaultValue, CUSTOM_WIDGET_TEMPLATE),
        },
      });
    };

    return {
      isValueEmpty,
      modalTitle,
      updateText,
      showTextEditorWithTemplateModal,
      showRemoveTextConfirmationModal,
    };
  },
};
</script>
