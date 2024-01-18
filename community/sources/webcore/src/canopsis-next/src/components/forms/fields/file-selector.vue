<template>
  <div>
    <div v-if="withFilesList">
      <div
        class="ml-2 font-italic"
        v-for="file in files"
        :key="file.name"
      >
        {{ file.name }}
        <v-btn
          icon
          text
          small
          @click="removeFileFromSelections(file.name)"
        >
          <v-icon small>
            close
          </v-icon>
        </v-btn>
      </div>
    </div>
    <v-layout
      v-on="wrapperListeners"
      align-center
    >
      <slot
        :disabled="fullDisabled"
        :loading="loading"
        :on="scopedActivatorSlotListeners"
        :clear="internalClear"
        :drop="dropFiles"
        :files="files"
        name="activator"
      >
        <v-btn
          v-on="scopedActivatorSlotListeners"
          :disabled="fullDisabled"
          :loading="loading"
        >
          <v-icon>cloud_upload</v-icon>
        </v-btn>
      </slot>
    </v-layout>
    <input
      class="hidden"
      ref="fileInput"
      type="file"
      :multiple="multiple"
      :accept="accept"
      @change="change"
      @click.stop=""
    >
    <div
      class="mt-2"
      v-if="!hideDetails"
    >
      <v-messages
        :value="errorMessages"
        color="error"
      />
    </div>
  </div>
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
      return this.loading || this.disabled;
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

    mimeType() {
      return !this.accept ? '' : new RegExp(this.accept.replace('*', '.*'));
    },
  },
  methods: {
    selectFiles(event) {
      event?.stopPropagation();

      if (!this.fullDisabled) {
        this.$refs.fileInput.click();
      }
    },

    dropFiles(event) {
      event?.preventDefault();

      if (event?.dataTransfer.items) {
        const files = [...event.dataTransfer.items]
          .filter(({ type }) => this.mimeType && this.mimeType.test(type))
          .map(item => item.getAsFile());

        if (!files.length) {
          this.errors.add({
            field: this.name,
            msg: this.$t('common.fileSelector.dragAndDrop.fileTypeError', { accept: this.accept }),
          });

          return;
        }

        this.files = this.multiple
          ? union(this.files, files)
          : files;

        this.errors.remove(this.name);
        this.$emit('change', files);
      }
    },

    change(event) {
      const files = Object.values(event.target.files);
      this.files = this.multiple
        ? union(this.files, files)
        : files;

      this.errors.remove(this.name);
      this.$emit('change', this.files);
    },

    clear() {
      this.$refs.fileInput.value = null;
      this.files = [];
    },

    internalClear() {
      this.clear();
      this.$emit('change', this.files);
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
</style>
