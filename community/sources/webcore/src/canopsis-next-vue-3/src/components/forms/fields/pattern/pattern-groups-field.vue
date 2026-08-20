<template>
  <c-layout gap="3" column>
    <v-alert v-if="!groups.length" type="info">
      {{ disabled || readonly ? $t('pattern.noDataDisabled') : $t('pattern.noData') }}
    </v-alert>
    <c-layout v-for="(group, index) in groups" :key="group.key" wrap>
      <pattern-group-field
        :group="groups[index]"
        :attributes="attributes"
        :disabled="disabled"
        :readonly="readonly"
        @update:group="updateItemInArray(index, $event)"
        @remove="removeItemFromArray(index)"
      />
      <c-layout v-show="index !== groups.length - 1" justify-center>
        <c-pattern-operator-chip>{{ $t('common.or') }}</c-pattern-operator-chip>
      </c-layout>
    </c-layout>

    <div class="align-self-start">
      <c-btn-with-error
        v-if="!readonly"
        :error="hasGroupsErrors ? $t('pattern.errors.groupRequired') : ''"
        :disabled="disabled"
        class="mt-2"
        @click="addFilterGroup"
      >
        {{ $t('pattern.addGroup') }}
      </c-btn-with-error>
    </div>
  </c-layout>
</template>

<script>
import { computed, nextTick, onBeforeUnmount, toRef, watch } from 'vue';

import { patternRulesToGroup } from '@/helpers/entities/pattern/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useValidationAttachMinValue } from '@/hooks/validator/validation-attach-min-value';

import PatternGroupField from './pattern-group-field.vue';

export default {
  components: { PatternGroupField },
  emits: ['update:groups'],
  props: {
    groups: {
      type: Array,
      required: true,
    },
    attributes: {
      type: Array,
      default: () => [],
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'groups',
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { addItemIntoArray, updateItemInArray, removeItemFromArray } = useArrayModelField(
      props,
      emit,
    );
    const fieldName = toRef(props, 'name');
    const { validator, attachMinValueRule, detachMinValueRule, validateMinValueRule } =
      useValidationAttachMinValue(fieldName);

    const hasGroupsErrors = computed(() => validator.errors.has(fieldName.value));
    const groupsLengthGetter = () => props.groups?.length ?? 0;

    /**
     * Attaches min_value:1 when the groups field is required.
     */
    const attachRequiredMinValueRule = () => attachMinValueRule(groupsLengthGetter);

    const validateGroupsField = async () => {
      if (!props.required) {
        validator.errors.clear(fieldName.value);

        return true;
      }

      return validateMinValueRule();
    };

    const addFilterGroup = () => {
      const [firstAttribute] = props.attributes;

      addItemIntoArray(patternRulesToGroup([{ field: firstAttribute?.value }]));
    };

    watch(
      () => props.groups,
      () => {
        if (props.required) {
          nextTick(validateMinValueRule);
        }
      },
    );

    watch(
      () => props.required,
      (required) => {
        if (required) {
          attachRequiredMinValueRule();
        } else {
          detachMinValueRule();
        }
      },
      { immediate: true },
    );

    onBeforeUnmount(detachMinValueRule);

    return {
      hasGroupsErrors,
      validateGroupsField,
      addFilterGroup,
      updateItemInArray,
      removeItemFromArray,
    };
  },
};
</script>
