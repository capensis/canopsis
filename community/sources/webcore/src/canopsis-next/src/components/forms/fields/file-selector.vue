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
      :multiple="multiple"
      :accept="accept"
      :name="name"
      type="file"
      @change="change"
      @click.stop=""
    >
    <div
      class="mt-2"
      v-if="!hideDetails"
    >
      <v-messages
        :value="errors.collect(name)"
        color="error"
      />
    </div>
  </div>
</template>

<script>
import { union } from 'lodash';

export default {
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
    /**
     * File size in kilobytes (KB)
     */
    maxFileSize: {
      type: [String, Number],
      required: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      files: [],
      filesForValidator: [],
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

    errorMessages() {
      return this.$validator.errors.collect();
    },
  },
  mounted() {
    this.attachField();
  },
  beforeDestroy() {
    this.detachField();
  },
  methods: {
    attachField() {
      const rules = {
        required: this.required,
      };

      if (this.accept) {
        rules.mimes = this.accept.split(',');
      }

      if (this.maxFileSize) {
        rules.size = Number(this.maxFileSize);
      }

      this.$validator.attach({
        name: this.name,
        rules,
        getter: () => this.filesForValidator,
        vm: this,
      });
    },

    detachField() {
      this.$validator.detach(this.name);
    },

    selectFiles(event) {
      event?.stopPropagation();

      if (!this.fullDisabled) {
        this.$refs.fileInput.click();
      }
    },

    async dropFiles(event) {
      const { files } = event?.dataTransfer ?? {};

      if (files) {
        this.setFilesForValidator(files);

        const isValid = await this.$validator.validate(this.name);

        if (isValid) {
          this.setFiles(this.filesForValidator);
        }
      }
    },

    async change(event) {
      this.setFilesForValidator(event.target.files);

      const isValid = await this.$validator.validate(this.name);

      if (isValid) {
        this.setFiles(this.filesForValidator);
      }
    },

    setFilesForValidator(files) {
      this.filesForValidator = this.multiple
        ? union(this.files, [...files])
        : [...files];
    },

    setFiles(files) {
      this.files = [...files];
      this.$emit('change', this.files);
    },

    clear() {
      this.$refs.fileInput.value = null;
      this.files = [];
      this.filesForValidator = [];

      this.errors.remove(this.name);
    },

    internalClear() {
      this.clear();
      this.$emit('change', this.files);
    },

    async removeFileFromSelections(name) {
      this.setFilesForValidator(this.files.filter(file => file.name !== name));

      const isValid = await this.$validator.validate(this.name);

      if (!isValid) {
        return;
      }

      this.setFiles(this.filesForValidator);

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
