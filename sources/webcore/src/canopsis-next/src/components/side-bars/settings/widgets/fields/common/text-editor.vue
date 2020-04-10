<template lang="pug">
  settings-button-field(
    :isEmpty="isValueEmpty",
    @create="openTextEditorModal",
    @edit="openTextEditorModal",
    @delete="deleteMoreInfoTemplate"
  )
    .subheading(slot="title") {{ title }}
</template>

<script>
import { MODALS } from '@/constants';

import SettingsButtonField from '../partials/button-field.vue';

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
    extraButtons: {
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
    openTextEditorModal() {
      this.$modals.show({
        name: MODALS.textEditor,
        config: {
          text: this.value,
          action: value => this.$emit('input', value),
          extraButtons: this.extraButtons,
        },
      });
    },

    deleteMoreInfoTemplate() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.$emit('input', ''),
        },
      });
    },
  },
};
</script>

