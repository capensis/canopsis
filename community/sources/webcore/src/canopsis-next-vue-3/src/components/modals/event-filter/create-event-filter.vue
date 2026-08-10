<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <event-filter-form
          :rule-id="ruleId"
          :is-disabled-id-field="config.submittingIdField"
          :event-attributes="eventAttributes"
          :attributes-pending="pending"
        />
      </template>
      <template #actions="">
        <v-btn
          :disabled="submitting"
          depressed
          text
          @click="closeModal"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="submitting"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, inject, watch } from 'vue';

import EventFilterForm from '@/components/other/event-filter/form/event-filter-form.vue';
import { MODALS } from '@/constants';
import { eventFilterToForm, formToEventFilter } from '@/helpers/entities/event-filter/rule/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';

export default {
  name: MODALS.createEventFilter,
  components: {
    EventFilterForm,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const system = inject('$system');

    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});
    const { t } = useI18n();

    const form = useForm({
      initialValues: eventFilterToForm(config.value.rule, system.timezone),
    });
    const { values, handleSubmit, submitting, reset } = form;

    const { fetchEventFilterPatternFields } = usePatternsFields();

    const { pending, eventAttributes } = usePatternsFieldsFetching(fetchEventFilterPatternFields);

    const ruleId = computed(() => config.value.rule?._id);
    const title = computed(() => config.value.title ?? t('modals.createEventFilter.create.title'));
    const isNew = computed(() => !config.value.rule?._id || config.value.duplicating);
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToEventFilter(values.value, system.timezone));

      await config.value.afterSubmit?.(result);

      closeModal();

      return result;
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    watch(
      () => values.value.type,
      () => reset(values.value),
    );

    return {
      config,
      ruleId,
      title,
      pending,
      eventAttributes,
      submitting,
      submitLabel,
      submit,
      closeModal,
    };
  },
};
</script>
