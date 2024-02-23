<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <v-list
          class="widget-settings widget-settings--divider widget-settings__list py-0"
          expand
        >
          <component
            :is="formComponent"
            v-model="form"
            :only-external="onlyExternal"
            required-title
          />
        </v-list>
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS, VALIDATION_DELAY, WIDGET_TYPES } from '@/constants';

import { alarmListChartToForm, formToAlarmListChart } from '@/helpers/entities/widget/forms/alarm';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import BarChartWidgetForm from '@/components/sidebars/chart/form/bar-chart-widget-form.vue';
import LineChartWidgetForm from '@/components/sidebars/chart/form/line-chart-widget-form.vue';
import NumbersWidgetForm from '@/components/sidebars/chart/form/numbers-widget-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createAlarmChart,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
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
      return this.config.title ?? this.$t(`modals.createAlarmChart.${WIDGET_TYPES.barChart}.create.title`);
    },

    onlyExternal() {
      return this.config.onlyExternal;
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
