<template>
  <v-row>
    <v-col v-show="title" cols="12">
      <c-layout>
        <h4 class="text-subtitle-1 text-grey-darken-2">
          {{ title }}
        </h4>
        <c-help-icon
          v-if="helpText"
          :text="helpText"
          icon-class="ml-2 storage-help-tooltip"
          max-width="250"
          top
        />
      </c-layout>
    </v-col>
    <v-col cols="12">
      <c-layout column>
        <c-text-field
          :model-value="modelValue"
          :label="label"
          :placeholder="placeholder"
          :disabled="disabled"
          :name="name"
          :required="required"
          @update:model-value="$emit('update:modelValue', $event)"
        />
      </c-layout>
    </v-col>
  </v-row>
</template>

<script>
export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    title: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    placeholder: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'mask',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    helpText: {
      type: String,
      required: false,
    },
  },

  emits: ['update:modelValue'],
};
</script>
