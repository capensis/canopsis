<template>
  <v-layout column>
    <v-layout>
      <v-flex xs6>
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
              top
            />
          </template>
        </c-number-field>
      </v-flex>
    </v-layout>
    <meta-alarm-rule-time-based-form v-field="config" with-child-inactive-delay />
    <c-payload-text-field
      v-field="config.corel_id"
      :label="$t('metaAlarmRule.corelId')"
      :variables="variables"
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
      :variables="variables"
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
    <c-payload-text-field
      v-field="config.corel_parent"
      :label="$t('metaAlarmRule.corelParent')"
      :variables="variables"
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
    </c-payload-text-field>
    <c-payload-text-field
      v-field="config.corel_child"
      :label="$t('metaAlarmRule.corelChild')"
      :variables="variables"
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
    </c-payload-text-field>
  </v-layout>
</template>

<script>
import MetaAlarmRuleTimeBasedForm from './meta-alarm-rule-time-based-form.vue';

export default {
  components: { MetaAlarmRuleTimeBasedForm },
  model: {
    prop: 'config',
    event: 'input',
  },
  props: {
    config: {
      type: Object,
      default: () => ({}),
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
