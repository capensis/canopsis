<template>
  <div class="meta-alarm-rule-type-field">
    <v-tabs
      v-model="activeTab"
      bg-color="transparent"
      color="primary"
      class="meta-alarm-rule-type-field__tabs"
      fixed-tabs
      hide-slider
    >
      <v-tab :value="TYPE_TABS.existing" class="meta-alarm-rule-type-field__tab">
        <v-btn class="meta-alarm-rule-type-field__tab-btn" variant="flat" block>
          {{ t('metaAlarmRule.groupingTabs.existing') }}
        </v-btn>
      </v-tab>
      <v-tab :value="TYPE_TABS.createNew" class="meta-alarm-rule-type-field__tab">
        <v-btn class="meta-alarm-rule-type-field__tab-btn" variant="flat" block>
          {{ t('metaAlarmRule.groupingTabs.createNew') }}
        </v-btn>
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab" touchless>
      <v-window-item :value="TYPE_TABS.existing">
        <c-layout class="pa-4" gap="2" column>
          <c-label :label="t('metaAlarmRule.groupingLabels.groupUnder')" />
          <c-radio-group
            :model-value="modelValue"
            name="meta_alarm_rule_type_existing"
            hide-details
            @update:model-value="updateModel"
          >
            <v-radio
              v-for="option in existingOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </c-radio-group>
        </c-layout>
      </v-window-item>

      <v-window-item :value="TYPE_TABS.createNew">
        <c-layout class="pa-4" gap="2" column>
          <c-label :label="t('metaAlarmRule.groupingLabels.groupBy')" />
          <c-radio-group
            :model-value="modelValue"
            name="meta_alarm_rule_type_create_new"
            hide-details
            @update:model-value="updateModel"
          >
            <v-radio
              v-for="option in createNewOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </c-radio-group>
        </c-layout>
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { computed, ref, watch } from 'vue';

import { META_ALARMS_RULE_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

const TYPE_TABS = {
  existing: 'existing',
  createNew: 'createNew',
};

const EXISTING_TYPES = [
  META_ALARMS_RULE_TYPES.relation,
  META_ALARMS_RULE_TYPES.corel,
];

const CREATE_NEW_TYPES = [
  META_ALARMS_RULE_TYPES.timebased,
  META_ALARMS_RULE_TYPES.attribute,
  META_ALARMS_RULE_TYPES.complex,
  META_ALARMS_RULE_TYPES.valuegroup,
];

export default {
  props: {
    modelValue: {
      type: String,
      default: META_ALARMS_RULE_TYPES.relation,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const activeTab = ref(TYPE_TABS.existing);

    const prepareOption = type => ({
      value: type,
      label: t(`metaAlarmRule.types.${type}.label`),
    });

    const existingOptions = computed(() => EXISTING_TYPES.map(prepareOption));
    const createNewOptions = computed(() => CREATE_NEW_TYPES.map(prepareOption));

    const updateActiveTab = (type) => {
      activeTab.value = CREATE_NEW_TYPES.includes(type)
        ? TYPE_TABS.createNew
        : TYPE_TABS.existing;
    };

    const updateModel = value => emit('update:modelValue', value);

    watch(
      () => props.modelValue,
      updateActiveTab,
      { immediate: true },
    );

    return {
      TYPE_TABS,
      t,
      activeTab,
      existingOptions,
      createNewOptions,
      updateModel,
    };
  },
};
</script>

<style lang="scss" scoped>
.meta-alarm-rule-type-field {
  --meta-alarm-rule-type-field-background: var(--v-application-background-darken1);

  overflow: hidden;
  border: thin solid var(--meta-alarm-rule-type-field-background);
  border-radius: 4px 4px 0 0;

  :deep(.meta-alarm-rule-type-field__tabs) {
    height: 44px;
  }

  :deep(.v-slide-group) {
    height: 44px;
    background: var(--meta-alarm-rule-type-field-background);
  }

  :deep(.v-slide-group__container),
  :deep(.v-slide-group__content) {
    width: 100%;
    height: 44px;
  }

  :deep(.v-tab) {
    flex: 1 1 0;
    align-items: center;
    justify-content: center;
    min-width: 0 !important;
    max-width: none !important;
    height: 44px;
    padding: 0 2px;
    text-transform: none;
  }

  :deep(.v-tab > .v-btn__content) {
    width: 100%;
  }

  :deep(.v-tab__slider) {
    display: none;
  }

  :deep(.v-tab--selected) {
    color: rgb(var(--v-theme-primary));

    .meta-alarm-rule-type-field__tab-btn {
      border-color: currentColor;
      background: rgb(var(--v-theme-surface));
      color: rgb(var(--v-theme-primary));
    }
  }

  &__tab-btn {
    width: 100%;
    height: 36px;
    min-width: 0;
    max-width: none;
    margin: 0;
    border: thin solid transparent;
    border-radius: 4px;
    background: var(--meta-alarm-rule-type-field-background);
    color: rgba(var(--v-theme-on-surface), 0.87);
  }

  :deep(.meta-alarm-rule-type-field__tab-btn),
  :deep(.meta-alarm-rule-type-field__tab-btn .v-btn__content) {
    letter-spacing: 0;
    text-transform: none !important;
  }
}
</style>
