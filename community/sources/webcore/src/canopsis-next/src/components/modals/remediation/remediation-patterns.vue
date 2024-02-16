<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <c-patterns-field
          v-model="form"
          :alarm-attributes="alarmAttributes"
          :entity-attributes="entityAttributes"
          with-alarm
          with-entity
          both-counters
        />
        <c-collapse-panel
          class="mt-3"
          :title="$t('remediation.pattern.tabs.pbehaviorTypes.title')"
        >
          <remediation-patterns-pbehavior-types-form v-model="form" />
        </c-collapse-panel>
      </template>
      <template #actions="">
        <v-btn
          :disabled="submitting"
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          class="primary"
          :disabled="isDisabled"
          :loading="submitting"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import {
  ALARM_PATTERN_FIELDS,
  ENTITY_PATTERN_FIELDS,
  MODALS,
  PATTERNS_FIELDS,
  QUICK_RANGES,
  VALIDATION_DELAY,
} from '@/constants';

import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import RemediationPatternsPbehaviorTypesForm from '@/components/other/remediation/patterns/form/remediation-patterns-pbehavior-types-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.remediationPatterns,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    RemediationPatternsPbehaviorTypesForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { instruction } = this.modal.config;

    return {
      form: {
        ...filterPatternsToForm(
          instruction,
          [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity],
        ),
        active_on_pbh: instruction?.active_on_pbh ?? [],
        disabled_on_pbh: instruction?.disabled_on_pbh ?? [],
      },
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.patterns.title');
    },

    intervalOptions() {
      return {
        intervalRanges: [QUICK_RANGES.custom],
      };
    },

    alarmAttributes() {
      return [
        {
          value: ALARM_PATTERN_FIELDS.creationDate,
          options: this.intervalOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.ackAt,
          options: this.intervalOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.lastUpdateDate,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.resolved,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.activationDate,
          options: { disabled: true },
        },
      ];
    },

    entityAttributes() {
      return [
        {
          value: ENTITY_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
      ];
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action({
            ...formFilterToPatterns(this.form, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
            active_on_pbh: this.form.active_on_pbh,
            disabled_on_pbh: this.form.disabled_on_pbh,
          });
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
