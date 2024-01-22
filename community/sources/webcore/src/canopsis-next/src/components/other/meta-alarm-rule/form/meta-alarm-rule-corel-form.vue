<template>
  <v-layout column>
    <v-text-field
      v-field.number="config.corel_id"
      v-validate="'required'"
      :label="$t('metaAlarmRule.corelId')"
      :error-messages="errors.collect('corelId')"
      name="corelId"
    >
      <template #append="">
        <c-help-icon
          :text="$t('metaAlarmRule.corelIdHelp')"
          icon="help"
          left
        />
      </template>
    </v-text-field>
    <v-text-field
      v-field.number="config.corel_status"
      v-validate="'required'"
      :label="$t('metaAlarmRule.corelStatus')"
      :error-messages="errors.collect('corelStatus')"
      name="corelStatus"
    >
      <template #append="">
        <c-help-icon
          :text="$t('metaAlarmRule.corelStatusHelp')"
          icon="help"
          left
        />
      </template>
    </v-text-field>
    <v-text-field
      v-field.number="config.corel_parent"
      v-validate="'required'"
      :label="$t('metaAlarmRule.corelParent')"
      :error-messages="errors.collect('corelParent')"
      name="corelParent"
    />
    <v-text-field
      v-field.number="config.corel_child"
      v-validate="'required'"
      :label="$t('metaAlarmRule.corelChild')"
      :error-messages="errors.collect('corelChild')"
      name="corelChild"
    />
    <c-number-field
      v-field="config.threshold_count"
      :label="$t('metaAlarmRule.thresholdCount')"
      :min="0"
      name="thresholdCount"
      required
    />
  </v-layout>
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'config',
    event: 'input',
  },
  props: {
    config: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
