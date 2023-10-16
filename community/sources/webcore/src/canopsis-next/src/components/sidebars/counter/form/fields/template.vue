<template>
  <v-container
    class="pa-3"
    fluid="fluid"
  >
    <v-layout
      align-center="align-center"
      justify-space-between="justify-space-between"
    >
      <div class="subheading">
        {{ title }}
      </div>
      <v-btn
        class="primary"
        small="small"
        @click="showTextEditorModal"
      >
        {{ $t('common.show') }}/{{ $t('common.edit') }}
      </v-btn>
    </v-layout>
  </v-container>
</template>

<script>
import { MODALS } from '@/constants';

export default {
  props: {
    value: {
      type: String,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
    variables: {
      type: Array,
      required: false,
    },
  },
  methods: {
    showTextEditorModal() {
      this.$modals.show({
        name: MODALS.textEditor,
        config: {
          title: this.title,
          text: this.value,
          variables: this.variables,
          action: value => this.$emit('input', value),
        },
      });
    },
  },
};
</script>
