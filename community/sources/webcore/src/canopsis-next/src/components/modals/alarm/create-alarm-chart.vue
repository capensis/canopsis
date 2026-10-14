<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        {{ title }}
      </template>
      <template #text>
        <v-list class="widget-settings widget-settings--divider widget-settings__list py-0">
          <component
            :is="formComponent"
            v-model="values"
            :only-external="onlyExternal"
          />
        </v-list>
      </template>
      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn :loading="isSubmitting" color="primary" type="submit">
          {{ t('common.submit') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { WIDGET_TYPES } from '@/constants';

import { alarmListChartToForm, formToAlarmListChart } from '@/helpers/entities/widget/forms/alarm';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import BarChartWidgetForm from '@/components/sidebars/chart/form/bar-chart-widget-form.vue';
import LineChartWidgetForm from '@/components/sidebars/chart/form/line-chart-widget-form.vue';
import NumbersWidgetForm from '@/components/sidebars/chart/form/numbers-widget-form.vue';

export default {
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});
    const form = useForm({ initialValues: alarmListChartToForm(config.value.chart) });
    const { values, handleSubmit, isSubmitting } = form;
    const formComponent = computed(
      () => ({
          [WIDGET_TYPES.barChart]: BarChartWidgetForm,
          [WIDGET_TYPES.lineChart]: LineChartWidgetForm,
          [WIDGET_TYPES.numbers]: NumbersWidgetForm,
        })[values.value.type],
    );
    const title = computed(
      () => config.value.title ?? t('modals.createAlarmChart.BarChart.create.title'),
    );
    const onlyExternal = computed(() => config.value.onlyExternal ?? false);
    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      await config.value.action?.(formToAlarmListChart(values.value));
      closeModal();
    });

    useFormConfirmableCloseModal({ form: values, submit, close: closeModal });

    return {
      t,
      values,
      formComponent,
      title,
      onlyExternal,
      isSubmitting,
      closeModal,
      submit,
    };
  },
};
</script>
