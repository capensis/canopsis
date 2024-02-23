<template>
  <div>
    <div v-if="withFilesList">
      <div
        v-for="file in files"
        :key="file.name"
        class="ml-2 font-italic"
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
    <div
      :class="{ disabled: fullDisabled }"
      class="file-selector-button-wrapper"
      v-on="wrapperListeners"
    >
      <slot
        :disabled="fullDisabled"
        :loading="loading"
        :on="scopedActivatorSlotListeners"
        name="activator"
      >
        <v-btn
          :disabled="fullDisabled"
          :loading="loading"
          v-on="scopedActivatorSlotListeners"
        >
          <v-icon>cloud_upload</v-icon>
        </v-btn>
      </slot>
    </div>
    <input
      ref="fileInput"
      :multiple="multiple"
      :accept="accept"
      class="hidden"
      type="file"
      @change="change"
      @click.stop=""
    >
    <div
      v-if="!hideDetails"
      class="mt-2"
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
  },
  methods: {
    selectFiles(event) {
      event.stopPropagation();

      if (!this.fullDisabled) {
        this.$refs.fileInput.click();
      }
    },

    change(e) {
      const files = Object.values(e.target.files);
      this.files = this.multiple
        ? union(this.files, files)
        : files;

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
