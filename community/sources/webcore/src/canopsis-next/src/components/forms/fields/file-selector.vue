<template lang="pug">
  div
    div(v-if="withFilesList")
      div.ml-2.font-italic(v-for="file in files", :key="file.name") {{ file.name }}
        v-btn(icon, flat, small, @click="removeFileFromSelections(file.name)")
          v-icon(small) close
    div.file-selector-button-wrapper(
      v-on="wrapperListeners",
      :class="{ disabled: fullDisabled }"
    )
      slot(
        :disabled="fullDisabled",
        :loading="loading",
        :on="scopedActivatorSlotListeners",
        name="activator"
      )
        v-btn(
          v-on="scopedActivatorSlotListeners",
          :disabled="fullDisabled",
          :loading="loading"
        )
          v-icon cloud_upload
    input.hidden(
      ref="fileInput",
      type="file",
      :multiple="multiple",
      :accept="accept",
      @change="change",
      @click.stop=""
    )
    div.mt-2(v-if="!hideDetails")
      div.error--text(v-for="error in errorMessages", :key="error") {{ error }}
</template>

<script>
import { union } from 'lodash';

export default {
  $_veeValidate: {
    value() {
      return this.files;
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
    disabled: {
      type: Boolean,
      default: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    withFilesList: {
      type: Boolean,
      default: false,
    },
    accept: {
      type: String,
      default: null,
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      files: [],
    };
  },
  computed: {
    fullDisabled() {
      return this.loading || this.disabled || (!this.multiple && Boolean(this.files.length));
    },

    hasActivatorSlot() {
      return this.$slots.activator && Boolean(this.$slots.activator.length);
    },

    wrapperListeners() {
      if (this.hasActivatorSlot) {
        return this.scopedActivatorSlotListeners;
      }

      return {};
    },

    scopedActivatorSlotListeners() {
      return {
        click: this.selectFiles,
      };
    },
  },
  methods: {
    selectFiles(event) {
      event.stopPropagation();

      if (!this.fullDisabled) {
        this.$refs.fileInput.click();
      }
    },

    change(e) {
      this.files = union(this.files, Object.values(e.target.files));

      this.$emit('change', this.files);
    },

    internalClear() {
      this.$refs.fileInput.value = null;
      this.files = [];

      this.$emit('change', this.files);
    },

    clear() {
      this.$refs.fileInput.value = null;
      this.files = [];
    },

    removeFileFromSelections(name) {
      this.files = this.files.filter(file => file.name !== name);

      if (!this.files.length) {
        this.internalClear();
      } else {
        this.$emit('change', this.files);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  .hidden {
    display: none;
  }

  .file-selector-button-wrapper {
    -webkit-box-align: center;
    -ms-flex-align: center;
    align-items: center;
    cursor: pointer;
    display: -webkit-box;
    display: -ms-flexbox;
    display: flex;

    &.disabled {
      cursor: default;
    }
  }
</style>
