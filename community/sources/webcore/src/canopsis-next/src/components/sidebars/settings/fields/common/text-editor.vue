<template lang="pug">
  settings-button-field(
    :is-empty="isValueEmpty",
    addable,
    removable,
    @create="showTextEditorModal",
    @edit="showTextEditorModal",
    @delete="showRemoveTextConfirmationModal"
  )
    template(#title="")
      div.subheading {{ title }}
</template>

<script>
import { MODALS } from '@/constants';

import { formMixin } from '@/mixins/form';

import SettingsButtonField from '@/components/sidebars/settings/fields/partials/button-field.vue';

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
  },
  computed: {
    isValueEmpty() {
      return !this.value || !this.value.length;
    },
  },
  methods: {
    showTextEditorModal() {
      this.$modals.show({
        name: MODALS.textEditor,
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
          action: () => this.updateModel(''),
        },
      });
    },
  },
};
</script>
