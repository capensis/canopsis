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
import { MODALS } from '@/constants';

import { formMixin } from '@/mixins/form';

import SettingsButtonField from '@/components/sidebars/form/fields/button-field.vue';

export default {
  components: { SettingsButtonField },
  mixins: [formMixin],
  props: {
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
  },
  computed: {
    isValueEmpty() {
      return this.defaultValue === String(this.value);
    },
  },
  methods: {
    showTextEditorModal() {
      this.$modals.show({
        name: MODALS.textEditor,
        dialogProps: this.dialogProps,
        config: {
          text: this.value,
          variables: this.variables,
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
