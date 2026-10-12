<template>
  <v-layout
    class="gap-3"
    column
  >
    <slot
      v-if="!exdates.length"
      name="no-data"
    />
    <c-form-block-array-field
      v-field="exdates"
      :item-to-form="pbehaviorExdateItemToForm"
      :add-button-label="addButtonLabel || $t('modals.createPbehaviorException.addDate')"
      :disabled="disabled"
      :required="required"
      :error-messages="errors.collect(name)"
      :required-error-message="resolvedRequiredErrorMessage"
    >
      <template #item="{ index, remove }">
        <pbehavior-exception-field
          v-field="exdates[index]"
          :disabled="disabled"
          :with-type="withExdateType"
          @delete="remove"
        />
      </template>
    </c-form-block-array-field>
  </v-layout>
</template>

<script>
import {
  computed,
  nextTick,
  onBeforeUnmount,
  onMounted,
  watch,
} from 'vue';

import { pbehaviorExdateItemToForm } from '@/helpers/entities/pbehavior/form';

import { useI18n } from '@/hooks/i18n';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';
import { useValidationAttachMinValue } from '@/hooks/validator/validation-attach-min-value';

import PbehaviorExceptionField from '@/components/other/pbehavior/exceptions/fields/pbehavior-exception-field.vue';

export default {
  inject: ['$validator'],
  components: { PbehaviorExceptionField },
  model: {
    prop: 'exdates',
    event: 'input',
  },
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
  setup(props) {
    const { t } = useI18n();
    const { fetchPbehaviorTypesFieldList } = usePbehaviorType();

    const {
      validator,
      attachMinValueRule,
      detachMinValueRule,
      validateMinValueRule,
    } = useValidationAttachMinValue(props.name);

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

    watch(() => props.exdates, () => {
      if (props.required) {
        asyncValidateMinValueRule();
      }
    });

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
    };
  },
};
</script>
