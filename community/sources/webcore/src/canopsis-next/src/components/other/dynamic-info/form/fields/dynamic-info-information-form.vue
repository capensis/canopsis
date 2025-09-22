<template>
  <v-layout class="gap-2" column>
    <v-text-field
      v-field="form.name"
      v-validate="nameRules"
      :label="$t('common.name')"
      :error-messages="errors.collect('name')"
      name="name"
      required
    />
    <dynamic-info-information-type-field :value="form.type" @input="changeType" />
    <c-mixed-field
      v-if="isDefaultType"
      v-field="form.value"
      :label="$t('common.value')"
      name="value"
      required
    />

    <c-payload-text-field
      v-else-if="isTemplateType"
      v-field="form.value"
      :label="$t('common.value')"
      :variables="isTemplateType ? variables : copyVariables"
      name="value"
      required
    />

    <v-combobox
      v-else
      v-field="form.value"
      v-validate="'required'"
      :label="$t('common.value')"
      :items="copyVariables"
      :menu-props="comboboxMenuProps"
      :error-messages="errors.collect('value')"
      :return-object="false"
      children-key="variables"
      name="value"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { DYNAMIC_INFO_INFORMATION_TYPES } from '@/constants';

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
  },
  setup(props, { emit }) {
    const validator = useValidator();

    const { updateModel } = useModelField(props, emit);

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

      validator.errors.remove('value');
    };

    return {
      nameRules,
      isDefaultType,
      isTemplateType,
      comboboxMenuProps,

      changeType,
    };
  },
};
</script>
