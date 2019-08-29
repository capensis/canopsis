<template lang="pug">
  v-container.pa-3(fluid)
    v-layout(align-center, justify-space-between)
      div.subheading {{ title }}
      v-layout(justify-end)
        v-btn.primary(
        small,
        @click="openTextEditorModal"
        )
          span(data-test="createButton", v-show="isValueEmpty") {{ $t('common.create') }}
          span(data-test="editButton", v-show="!isValueEmpty") {{ $t('common.edit') }}
        v-btn.error(v-show="!isValueEmpty", small, @click="deleteMoreInfoTemplate")
          v-icon(data-test="deleteButton",) delete
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';

export default {
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

