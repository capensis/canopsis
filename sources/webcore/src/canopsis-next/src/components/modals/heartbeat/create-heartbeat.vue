<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createHeartbeat.create.title') }}
    v-card-text
      v-form
        v-layout(row, wrap)
          v-flex(xs3)
            v-text-field(
            v-model="periodForm.periodValue",
            v-validate="'required'",
            :label="$t('modals.statsDateInterval.fields.periodValue')",
            :error-messages="errors.collect('periodValue')",
            type="number",
            name="periodValue"
            )
          v-flex
            v-select(
            v-model="periodForm.periodUnit",
            v-validate="'required'",
            :items="periodUnits",
            :label="$t('modals.statsDateInterval.fields.periodUnit')",
            :error-messages="errors.collect('periodUnit')",
            name="periodUnit"
            )
        v-layout
          v-btn(@click="showEditPatternModal") {{ $t('modals.eventFilterRule.editPattern') }}
    v-alert(:value="errors.has('pattern')", type="error") {{ $t('modals.createHeartbeat.patternRequired') }}
    v-divider
    v-layout.py-1(justify-end)
      v-btn(depressed, flat, @click="hideModal") {{ $t('common.cancel') }}
      v-btn(color="primary", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { isEmpty } from 'lodash';

import { MODALS, HEARTBEAT_DURATION_UNITS } from '@/constants';

import popupMixin from '@/mixins/popup';
import modalInnerMixin from '@/mixins/modal/inner';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createHeartbeat,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [popupMixin, modalInnerMixin],
  data() {
    return {
      periodForm: {
        periodValue: '',
        periodUnit: '',
      },
      form: {
        pattern: {},
      },
    };
  },
  computed: {
    periodUnits() {
      return [
        {
          text: this.$tc('common.times.minute'),
          value: HEARTBEAT_DURATION_UNITS.minute,
        },
        {
          text: this.$tc('common.times.hour'),
          value: HEARTBEAT_DURATION_UNITS.hour,
        },
      ];
    },
  },
  created() {
    this.$validator.attach({
      name: 'pattern',
      rules: 'required:true',
      getter: () => !isEmpty(this.form.pattern),
      context: () => this,
    });
  },
  methods: {
    showEditPatternModal() {
      this.showModal({
        name: MODALS.createEventFilterRulePattern,
        config: {
          isSimplePattern: true,
          pattern: this.form.pattern,
          action: (pattern) => {
            this.form.pattern = pattern;

            this.$validator.validate('pattern');
          },
        },
      });
    },

    async submit() {
      try {
        const isValid = await this.$validator.validateAll();

        if (isValid) {
          const { periodValue, periodUnit } = this.periodForm;
          const { pattern } = this.form;
          const data = {
            pattern,
            expected_interval: `${periodValue}${periodUnit}`,
          };

          if (this.config.action) {
            await this.config.action(data);
          }

          this.hideModal();
        }
      } catch (err) {
        this.addErrorPopup({ text: err.description });
      }
    },
  },
};
</script>
