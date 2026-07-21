<template>
  <v-layout class="gap-2" column>
    <v-layout align-center justify-space-between>
      <c-name-field
        v-field="form.name"
        v-validate="nameRules"
        :error-messages="errors.collect(nameFieldName)"
        :name="nameFieldName"
        required
      />
      <v-flex v-if="removable" shrink>
        <div>
          <c-action-btn type="delete" @click="remove" />
        </div>
      </v-flex>
    </v-layout>
    <dynamic-info-information-type-field :value="form.type" @input="changeType" />
    <c-mixed-field
      v-if="isDefaultType"
      v-field="form.value"
      :label="$t('common.value')"
      :types="mixedFieldTypes"
      :name="valueFieldName"
      required
    />

    <c-payload-text-field
      v-else-if="isTemplateType"
      v-field="form.value"
      :label="$t('common.value')"
      :variables="isTemplateType ? variables : copyVariables"
      :name="valueFieldName"
      required
    />

    <c-select-field
      v-else
      v-field="form.value"
      :label="$t('common.value')"
      :items="copyVariables"
      :name="valueFieldName"
      :menu-props="comboboxMenuProps"
      :return-object="false"
      children-key="variables"
      required
      combobox
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { DYNAMIC_INFO_INFORMATION_TYPES, PATTERN_FIELD_TYPES } from '@/constants';

import { useValidator } from '@/hooks/validator/validator';
import { useModelField } from '@/hooks/form/model-field';

import DynamicInfoInformationTypeField from './dynamic-info-information-type-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DynamicInfoInformationTypeField,
  },
  model: {
    prop: 'form',
    event: 'input',
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
    existingNames: {
      type: Array,
      default: () => [],
    },
    initialName: {
      type: String,
      default: '',
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
  setup(props, { emit }) {
    const validator = useValidator();

    const mixedFieldTypes = [
      { value: PATTERN_FIELD_TYPES.string },
      { value: PATTERN_FIELD_TYPES.number },
      { value: PATTERN_FIELD_TYPES.boolean },
      { value: PATTERN_FIELD_TYPES.stringArray },
    ];

    const { updateModel } = useModelField(props, emit);

    const getFieldName = field => `${props.name}.${field}`;

    const nameFieldName = computed(() => getFieldName('name'));
    const valueFieldName = computed(() => getFieldName('value'));

    const nameRules = computed(() => ({
      required: true,
      unique: {
        values: props.existingNames,
        initialValue: props.initialName,
      },
    }));

    const isDefaultType = computed(() => props.form.type === DYNAMIC_INFO_INFORMATION_TYPES.setToInfo);
    const isTemplateType = computed(() => props.form.type === DYNAMIC_INFO_INFORMATION_TYPES.setToInfoFromTemplate);

    const comboboxMenuProps = computed(() => ({
      minWidth: 200,
    }));

    const changeType = (type) => {
      updateModel({
        ...props.form,

        type,
        value: '',
      });

      validator.errors.remove(valueFieldName.value);
    };

    const remove = () => emit('remove', props.form);

    return {
      mixedFieldTypes,

      nameFieldName,
      valueFieldName,
      nameRules,
      isDefaultType,
      isTemplateType,
      comboboxMenuProps,

      changeType,
      remove,
    };
  },
};
</script>
