<template>
  <c-layout gap="3" column>
    <slot v-if="!exdates.length" name="no-data" />
    <c-form-block-array-field
      :form="exdates"
      :item-to-form="pbehaviorExdateItemToForm"
      :add-button-label="addButtonLabel || $t('modals.createPbehaviorException.addDate')"
      :disabled="disabled"
      :required="required"
      :error-messages="errors.collect(name)"
      :required-error-message="resolvedRequiredErrorMessage"
      @update:form="updateModel"
    >
      <template #item="{ index, remove }">
        <pbehavior-exception-field
          :value="exdates[index]"
          :disabled="disabled"
          :with-type="withExdateType"
          @update:value="updateItemInArray(index, $event)"
          @delete="remove"
        />
      </template>
    </c-form-block-array-field>
  </c-layout>
</template>

<script>
import { computed, nextTick, onBeforeUnmount, onMounted, watch } from 'vue';

import PbehaviorExceptionField from '@/components/other/pbehavior/exceptions/fields/pbehavior-exception-field.vue';
import { pbehaviorExdateItemToForm } from '@/helpers/entities/pbehavior/form';
import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useI18n } from '@/hooks/i18n';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';
import { useValidationAttachMinValue } from '@/hooks/validator/validation-attach-min-value';

export default {
  emits: ['update:exdates'],
  components: { PbehaviorExceptionField },
  props: {
    exdates: {
      type: Array,
      default: () => [],
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    withExdateType: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'exdates',
    },
    requiredErrorMessage: {
      type: String,
      default: '',
    },
    addButtonLabel: {
      type: String,
      default: '',
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel, updateItemInArray } = useArrayModelField(props, emit);
    const { fetchPbehaviorTypesFieldList } = usePbehaviorType();

    const { validator, attachMinValueRule, detachMinValueRule, validateMinValueRule } =
      useValidationAttachMinValue(props.name);

    const resolvedRequiredErrorMessage = computed(
      () => props.requiredErrorMessage || t('common.addAtLeastOneItem'),
    );

    /**
     * Returns the current number of exception dates for min_value validation.
     *
     * @returns {number}
     */
    const exdatesLengthGetter = () => props.exdates?.length ?? 0;

    /**
     * Re-validates the exdates min_value rule on the next tick after the array changes.
     */
    const asyncValidateMinValueRule = () => nextTick(validateMinValueRule);

    watch(
      () => props.exdates,
      () => {
        if (props.required) {
          asyncValidateMinValueRule();
        }
      },
    );

    onMounted(() => {
      fetchPbehaviorTypesFieldList();

      if (!props.required) {
        return;
      }

      attachMinValueRule(exdatesLengthGetter);
    });

    onBeforeUnmount(detachMinValueRule);

    return {
      errors: validator.errors,
      pbehaviorExdateItemToForm,
      resolvedRequiredErrorMessage,
      updateModel,
      updateItemInArray,
    };
  },
};
</script>
