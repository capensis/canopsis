<template>
  <div>
    <c-form-block-row :label="t('metaAlarmRule.threshold')">
      <c-number-field
        :value="config.threshold_count"
        :min="0"
        :label="t('metaAlarmRule.thresholdCount')"
        name="threshold_count"
        @update:value="updateField('threshold_count', $event)"
      >
        <template #append>
          <c-help-icon
            :text="t('metaAlarmRule.thresholdCountHelpText')"
            max-width="300"
            icon="help"
            left
          />
        </template>
      </c-number-field>
    </c-form-block-row>

    <meta-alarm-rule-time-based-field
      :model-value="config.time_interval"
      @update:model-value="updateField('time_interval', $event)"
    />

    <meta-alarm-rule-child-inactive-delay-field
      :model-value="config.child_inactive_delay"
      @update:model-value="updateField('child_inactive_delay', $event)"
    />

    <c-form-block-row :label="t('common.correlation')" indented>
      <c-layout gap="3" column>
        <c-payload-text-field
          :value="config.corel_id"
          :label="t('metaAlarmRule.corelId')"
          :variables="templateVars.corel_id"
          name="corel_id"
          required
          @update:value="updateField('corel_id', $event)"
        >
          <template #append-outer>
            <c-help-icon
              :text="t('metaAlarmRule.corelIdHelpText')"
              max-width="300"
              icon="help"
              left
            />
          </template>
        </c-payload-text-field>

        <c-payload-text-field
          :value="config.corel_status"
          :label="t('metaAlarmRule.corelStatus')"
          :variables="templateVars.corel_status"
          name="corel_status"
          required
          @update:value="updateField('corel_status', $event)"
        >
          <template #append-outer>
            <c-help-icon
              :text="t('metaAlarmRule.corelStatusHelpText')"
              max-width="300"
              icon="help"
              left
            />
          </template>
        </c-payload-text-field>

        <c-text-field
          :model-value="config.corel_parent"
          :label="t('metaAlarmRule.corelParent')"
          name="corel_parent"
          required
          @update:model-value="updateField('corel_parent', $event)"
        >
          <template #append>
            <c-help-icon
              :text="t('metaAlarmRule.corelParentHelpText')"
              max-width="300"
              icon="help"
              left
            />
          </template>
        </c-text-field>

        <c-text-field
          :model-value="config.corel_child"
          :label="t('metaAlarmRule.corelChild')"
          name="corel_child"
          required
          @update:model-value="updateField('corel_child', $event)"
        >
          <template #append>
            <c-help-icon
              :text="t('metaAlarmRule.corelChildHelpText')"
              max-width="300"
              icon="help"
              left
            />
          </template>
        </c-text-field>

        <v-expand-transition>
          <c-layout v-if="sanitizedSummary" gap="2">
            <span class="text-subtitle-2">{{ t('common.summary') }}:</span>
            <span v-html="sanitizedSummary" class="pre-wrap" />
          </c-layout>
        </v-expand-transition>
      </c-layout>
    </c-form-block-row>
  </div>
</template>

<script>
import { computed } from 'vue';

import { sanitizeHtml } from '@/helpers/html';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

import MetaAlarmRuleChildInactiveDelayField from './fields/meta-alarm-rule-child-inactive-delay-field.vue';
import MetaAlarmRuleTimeBasedField from './fields/meta-alarm-rule-time-based-field.vue';

export default {
  components: {
    MetaAlarmRuleChildInactiveDelayField,
    MetaAlarmRuleTimeBasedField,
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

  emits: ['update:config'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    const sanitizedSummary = computed(() => {
      if (
        props.config.corel_id &&
        props.config.corel_status &&
        props.config.corel_child &&
        props.config.corel_parent
      ) {
        return sanitizeHtml(
          t('metaAlarmRule.corelGroupingSummary', {
            corelID: `{{ ${props.config.corel_id} }}`,
            corelChild: `{{ ${props.config.corel_child} }}`,
            corelParent: `{{ ${props.config.corel_parent} }}`,
            corel: `{{ ${props.config.corel_status} }}`,
          }),
        );
      }

      return '';
    });

    return {
      t,
      sanitizedSummary,
      updateField,
    };
  },
};
</script>
