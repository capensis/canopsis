<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        v-list.widget-settings__list.py-0(expand)
          component(:is="formComponent", v-model="form")
      template(#actions="")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS, WIDGET_TYPES } from '@/constants';

import { alarmListChartToForm, formToAlarmListChart } from '@/helpers/forms/widgets/alarm';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import BarChartWidgetForm from '@/components/sidebars/settings/forms/widgets/bar-chart-widget-form.vue';
import LineChartWidgetForm from '@/components/sidebars/settings/forms/widgets/line-chart-widget-form.vue';
import NumbersWidgetForm from '@/components/sidebars/settings/forms/widgets/numbers-widget-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create an ack event
 */
export default {
  name: MODALS.createAlarmChart,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    BarChartWidgetForm,
    LineChartWidgetForm,
    NumbersWidgetForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: alarmListChartToForm(this.modal.config.chart),
    };
  },
  computed: {
    formComponent() {
      return {
        [WIDGET_TYPES.barChart]: 'bar-chart-widget-form',
        [WIDGET_TYPES.lineChart]: 'line-chart-widget-form',
        [WIDGET_TYPES.numbers]: 'numbers-widget-form',
      }[this.form.type] ?? 'div';
    },

    title() {
      return this.config.title ?? this.$t('modals.createAlarmChart.barChart.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config?.action(formToAlarmListChart(this.form));

        this.$modals.hide();
      }
    },
  },
};
</script>
