<template>
  <div>
    <c-form-block-row :label="$t('metaAlarmRule.threshold')">
      <c-number-field
        v-field="config.threshold_count"
        :min="0"
        :label="$t('metaAlarmRule.thresholdCount')"
        name="thresholdCount"
      >
        <template #append-outer>
          <c-help-icon
            :text="$t('metaAlarmRule.thresholdCountHelpText')"
            max-width="300"
            icon="help"
            left
          />
        </template>
      </c-number-field>
    </c-form-block-row>

    <meta-alarm-rule-time-based-field v-field="config.time_interval" />

    <meta-alarm-rule-child-inactive-delay-field v-field="config.child_inactive_delay" />

    <c-form-block-row :label="$t('common.correlation')" indented>
      <c-payload-text-field
        v-field="config.corel_id"
        :label="$t('metaAlarmRule.corelId')"
        :variables="templateVars.corel_id"
        name="corelId"
        required
      >
        <template #append-outer="">
          <c-help-icon
            :text="$t('metaAlarmRule.corelIdHelpText')"
            max-width="300"
            icon="help"
            left
          />
        </template>
      </c-payload-text-field>
      <c-payload-text-field
        v-field="config.corel_status"
        :label="$t('metaAlarmRule.corelStatus')"
        :variables="templateVars.corel_status"
        name="corelStatus"
        required
      >
        <template #append-outer="">
          <c-help-icon
            :text="$t('metaAlarmRule.corelStatusHelpText')"
            max-width="300"
            icon="help"
            left
          />
        </template>
      </c-payload-text-field>
      <v-text-field
        v-field="config.corel_parent"
        v-validate="'required'"
        :label="$t('metaAlarmRule.corelParent')"
        :error-messages="errors.collect('corelParent')"
        name="corelParent"
        required
      >
        <template #append-outer="">
          <c-help-icon
            :text="$t('metaAlarmRule.corelParentHelpText')"
            max-width="300"
            icon="help"
            left
          />
        </template>
      </v-text-field>
      <v-text-field
        v-field="config.corel_child"
        v-validate="'required'"
        :label="$t('metaAlarmRule.corelChild')"
        :error-messages="errors.collect('corelChild')"
        name="corelChild"
        required
      >
        <template #append-outer="">
          <c-help-icon
            :text="$t('metaAlarmRule.corelChildHelpText')"
            max-width="300"
            icon="help"
            left
          />
        </template>
      </v-text-field>
      <v-expand-transition>
        <v-layout v-show="!!sanitizedSummary" class="gap-2">
          <span class="text-subtitle-2">{{ $t('common.summary') }}: </span>
          <span v-html="sanitizedSummary" class="pre-wrap" />
        </v-layout>
      </v-expand-transition>
    </c-form-block-row>
  </div>
</template>

<script>
import { computed } from 'vue';

import { sanitizeHtml } from '@/helpers/html';

import { useI18n } from '@/hooks/i18n';

import MetaAlarmRuleChildInactiveDelayField from './fields/meta-alarm-rule-child-inactive-delay-field.vue';
import MetaAlarmRuleTimeBasedField from './fields/meta-alarm-rule-time-based-field.vue';

export default {
  inject: ['$validator'],
  components: {
    MetaAlarmRuleTimeBasedField,
    MetaAlarmRuleChildInactiveDelayField,
  },
  model: {
    prop: 'config',
    event: 'input',
  },
  props: {
    config: {
      type: Object,
      default: () => ({}),
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const { t } = useI18n();

    const sanitizedSummary = computed(() => {
      if (props.config.corel_id && props.config.corel_status && props.config.corel_child && props.config.corel_parent) {
        return sanitizeHtml(t('metaAlarmRule.corelGroupingSummary', {
          corelID: `{{ ${props.config.corel_id} }}`,
          corelChild: `{{ ${props.config.corel_child} }}`,
          corelParent: `{{ ${props.config.corel_parent} }}`,
          corel: `{{ ${props.config.corel_status} }}`,
        }));
      }

      return '';
    });

    return {
      sanitizedSummary,
    };
  },
};
</script>
