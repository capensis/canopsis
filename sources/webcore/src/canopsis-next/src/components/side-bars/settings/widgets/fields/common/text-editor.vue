<template lang="pug">
  settings-button-field(
  :isEmpty="isValueEmpty",
  @create="openTextEditorModal",
  @edit="openTextEditorModal",
  @delete="deleteMoreInfoTemplate",
  )
    .subheading(slot="title") {{ title }}
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';

import SettingsButtonField from '../partials/button-field.vue';

export default {
  components: { SettingsButtonField },
  mixins: [modalMixin],
  props: {
    value: {
      type: String,
      default: '',
    },
    title: {
      type: String,
      default: '',
    },
  },
  computed: {
    isValueEmpty() {
      return !this.value || !this.value.length;
    },
  },
  methods: {
    openTextEditorModal() {
      this.showModal({
        name: MODALS.textEditor,
        config: {
          text: this.value,
          action: value => this.$emit('input', value),
        },
      });
    },

    deleteMoreInfoTemplate() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.$emit('input', ''),
        },
      });
    },
  },
};
</script>

