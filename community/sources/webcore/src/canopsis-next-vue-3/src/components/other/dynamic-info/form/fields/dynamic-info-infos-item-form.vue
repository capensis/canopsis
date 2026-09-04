<template>
  <c-layout gap="2" column>
    <c-layout align="center" justify="space-between" gap="2">
      <c-name-field
        :model-value="form.name"
        :name="nameFieldName"
        :unique="uniqueNames"
        required
        @update:model-value="updateField('name', $event)"
      />
      <c-action-btn v-if="removable" type="delete" @click="remove" />
    </c-layout>

    <dynamic-info-information-type-field
      :model-value="form.type"
      @update:model-value="changeType"
    />

    <c-mixed-field
      v-if="isDefaultType"
      :value="form.value"
      :label="$t('common.value')"
      :types="mixedFieldTypes"
      :name="valueFieldName"
      required
      @update:value="updateField('value', $event)"
    />

    <c-payload-text-field
      v-else-if="isTemplateType"
      :value="form.value"
      :label="$t('common.value')"
      :variables="variables"
      :name="valueFieldName"
      required
      @update:value="updateField('value', $event)"
    />

    <c-combobox
      v-else
      :model-value="form.value"
      :label="$t('common.value')"
      :items="copyVariables"
      :name="valueFieldName"
      :menu-props="comboboxMenuProps"
      children-key="variables"
      :return-object="false"
      required
      @update:model-value="updateCopyValue"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { DYNAMIC_INFO_INFORMATION_TYPES, PATTERN_FIELD_TYPES } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';

import DynamicInfoInformationTypeField from './dynamic-info-information-type-field.vue';

export default {
  components: {
    DynamicInfoInformationTypeField,
  },

  props: {
    form: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: '',
    },
    uniqueNames: {
      type: Array,
      default: () => [],
    },
    variables: {
      type: Array,
      default: () => [],
    },
    copyVariables: {
      type: Array,
      default: () => [],
    },
    removable: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:form', 'remove'],

  setup(props, { emit }) {
    const mixedFieldTypes = [
      { value: PATTERN_FIELD_TYPES.string },
      { value: PATTERN_FIELD_TYPES.number },
      { value: PATTERN_FIELD_TYPES.boolean },
      { value: PATTERN_FIELD_TYPES.stringArray },
    ];

    const { updateField, updateModel } = useModelField(props, emit);

    const nameFieldName = computed(() => `${props.name}.name`);
    const valueFieldName = computed(() => `${props.name}.value`);

    const isDefaultType = computed(
      () => props.form.type === DYNAMIC_INFO_INFORMATION_TYPES.setToInfo,
    );
    const isTemplateType = computed(
      () => props.form.type === DYNAMIC_INFO_INFORMATION_TYPES.setToInfoFromTemplate,
    );

    const comboboxMenuProps = computed(() => ({
      minWidth: 200,
    }));

    const changeType = (type) => {
      updateModel({
        ...props.form,
        type,
        value: '',
      });
    };

    const updateCopyValue = newValue => updateField('value', newValue?.value || newValue);

    const remove = () => emit('remove', props.form);

    return {
      mixedFieldTypes,
      nameFieldName,
      valueFieldName,
      isDefaultType,
      isTemplateType,
      comboboxMenuProps,
      updateField,
      changeType,
      updateCopyValue,
      remove,
    };
  },
};
</script>
