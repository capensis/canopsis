<template>
  <div>
    <div v-if="withFilesList">
      <div v-for="file in selectedFiles" :key="file.name" class="ml-2 font-italic">
        {{ file.name }}
        <v-btn
          icon="close"
          variant="text"
          size="small"
          @click="removeFileFromSelections(file.name)"
        />
      </div>
    </div>

    <slot
      :on="{ click: selectFile }"
      :attrs="{ disabled: fullDisabled }"
      :disabled="fullDisabled"
      :loading="loading"
      :files="selectedFiles"
      :clear="clear"
      :drop="dropFiles"
      name="activator"
    >
      <v-btn :disabled="fullDisabled" :loading="loading" @click="selectFile">
        {{ label || t('common.import') }}
      </v-btn>
    </slot>

    <input
      ref="inputElement"
      :accept="accept"
      :multiple="multiple"
      :disabled="fullDisabled"
      class="d-none"
      type="file"
      @change="change"
      @click.stop
    />

    <v-messages
      v-if="!hideDetails && errorMessages.length"
      :messages="errorMessages"
      color="error"
    />
  </div>
</template>

<script>
import { computed, ref, watch } from 'vue';
import { useField } from 'vee-validate';
import { unionBy } from 'lodash';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: [File, Array],
      default: undefined,
    },
    label: {
      type: String,
      default: '',
    },
    accept: {
      type: String,
      default: '',
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    /**
     * File size in kilobytes (KB)
     */
    maxFileSize: {
      type: [Number, String],
      required: false,
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
    hideDetails: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'file',
    },
  },

  emits: ['update:modelValue', 'change'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const inputElement = ref(null);
    const selectedFiles = ref([]);
    const fieldName = computed(() => props.name || 'file');
    const fullDisabled = computed(() => props.disabled || props.loading);
    const fieldRules = computed(() => ({
      required: props.required,
      size: props.maxFileSize ? Number(props.maxFileSize) : false,
    }));
    const { errors, setValue } = useField(fieldName, fieldRules, {
      initialValue: props.modelValue,
      validateOnMount: false,
      validateOnValueUpdate: false,
    });

    const selectFile = (event) => {
      event?.stopPropagation();

      if (!fullDisabled.value) {
        inputElement.value?.click();
      }
    };

    const updateFiles = (files) => {
      const value = props.multiple ? files : files[0];

      selectedFiles.value = files;
      setValue(value, false);
      emit('update:modelValue', value);
      emit('change', value);
    };

    const change = ({ target }) => {
      const files = Array.from(target.files ?? []);
      const nextFiles = props.multiple ? unionBy(selectedFiles.value, files, 'name') : files;

      updateFiles(nextFiles);
    };

    const dropFiles = (event) => {
      const files = Array.from(event?.dataTransfer?.files ?? []);
      const nextFiles = props.multiple ? unionBy(selectedFiles.value, files, 'name') : files;

      updateFiles(nextFiles);
    };

    const clear = () => {
      if (inputElement.value) {
        inputElement.value.value = '';
      }

      updateFiles([]);
    };

    const removeFileFromSelections = (fileName) => {
      updateFiles(selectedFiles.value.filter((file) => file.name !== fileName));
    };

    watch(
      () => props.modelValue,
      (nextValue) => {
        selectedFiles.value = Array.isArray(nextValue) ? nextValue : [nextValue].filter(Boolean);
        setValue(nextValue, false);
      },
    );

    return {
      t,
      inputElement,
      selectedFiles,
      errorMessages: errors,
      fullDisabled,
      selectFile,
      change,
      dropFiles,
      clear,
      removeFileFromSelections,
    };
  },
};
</script>
