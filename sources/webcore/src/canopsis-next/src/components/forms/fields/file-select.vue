<template lang="pug">
  v-tooltip(top)
    v-btn(
    slot="activator",
    v-bind="btnProps",
    @click.stop="selectFiles"
    )
      slot(name="label")
        v-icon cloud_upload
    input.hidden(
    ref="fileInput",
    type="file",
    :multiple="multiple",
    @change="$emit('change', $event)"
    )
    span {{ tooltip }}
</template>

<script>
export default {
  $_veeValidate: {
    value() {
      return this.$refs.fileInput.value;
    },

    name() {
      return this.name;
    },
  },
  inject: ['$validator'],
  props: {
    name: {
      type: String,
      default: null,
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    tooltip: {
      type: String,
      default: '',
    },
    btnProps: {
      type: Object,
      default: () => ({}),
    },
  },
  methods: {
    selectFiles() {
      this.$refs.fileInput.click();
    },

    clear() {
      this.$refs.fileInput.value = null;
    },
  },
};
</script>

<style scoped>
  .hidden {
    display: none;
  }
</style>
