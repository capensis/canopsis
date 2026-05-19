<template>
  <v-tabs
    v-model="activeTab"
    class="meta-alarm-rule-type-field"
    fixed-tabs
    hide-slider
  >
    <v-tab :value="META_ALARMS_TYPE_TABS.existing">
      <v-btn depressed>
        {{ $t('metaAlarmRule.groupingTabs.existing') }}
      </v-btn>
    </v-tab>
    <v-tab-item>
      <c-label :label="$t('metaAlarmRule.groupingLabels.groupUnder')" />
      <v-radio-group v-field="value" hide-details>
        <v-radio
          v-for="option in firstTabOptions"
          :key="option.value"
          :label="option.label"
          :value="option.value"
          color="primary"
        />
      </v-radio-group>
    </v-tab-item>

    <v-tab :value="META_ALARMS_TYPE_TABS.createNew">
      <v-btn depressed>
        {{ $t('metaAlarmRule.groupingTabs.createNew') }}
      </v-btn>
    </v-tab>
    <v-tab-item>
      <c-label :label="$t('metaAlarmRule.groupingLabels.groupBy')" />
      <v-radio-group v-field="value" hide-details>
        <v-radio
          v-for="option in secondTabOptions"
          :key="option.value"
          :label="option.label"
          :value="option.value"
          color="primary"
        />
      </v-radio-group>
    </v-tab-item>
  </v-tabs>
</template>

<script>
import { computed, ref, onMounted } from 'vue';

import { META_ALARMS_RULE_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export const META_ALARMS_TYPE_TABS = {
  existing: 0,
  createNew: 1,
};

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: META_ALARMS_RULE_TYPES.relation,
    },
  },
  setup(props) {
    const { t, te } = useI18n();

    const activeTab = ref(META_ALARMS_TYPE_TABS.existing);

    const prepareOption = (type) => {
      const messageKey = `metaAlarmRule.types.${type}`;
      const { label } = te(messageKey) ? t(messageKey) : {};

      return {
        value: type,
        label,
      };
    };

    const firstTabOptions = computed(() => [
      META_ALARMS_RULE_TYPES.relation,
      META_ALARMS_RULE_TYPES.corel,
    ].map(prepareOption));

    const secondTabOptions = computed(() => [
      META_ALARMS_RULE_TYPES.timebased,
      META_ALARMS_RULE_TYPES.attribute,
      META_ALARMS_RULE_TYPES.complex,
      META_ALARMS_RULE_TYPES.valuegroup,
    ].map(prepareOption));

    onMounted(() => {
      if (secondTabOptions.value.some(option => option.value === props.value)) {
        activeTab.value = META_ALARMS_TYPE_TABS.createNew;
      }
    });

    return {
      META_ALARMS_TYPE_TABS,

      activeTab,

      firstTabOptions,
      secondTabOptions,
    };
  },
};
</script>

<style lang="scss" scoped>
.meta-alarm-rule-type-field {
  --tab-background-color-light: var(--v-application-background-darken1);
  --tab-background-color-dark: var(--v-application-background-lighten1);

  border-radius: 4px 4px 0 0;
  border: thin solid var(--tab-background-color-light);

  .theme--dark & {
    border-color: var(--tab-background-color-dark);
  }

  ::v-deep {
    > .v-tabs-items {
      > .v-window__container > .v-window-item {
        padding: 16px;
      }
    }

    > .v-tabs-bar {
      height: 44px;

      .theme--light & {
        background-color: var(--tab-background-color-light);
      }

      .theme--dark & {
        background-color: var(--tab-background-color-dark);
      }

      .v-tab {
        padding: 0 4px;

        .v-btn {
          width: 100%;
          border-radius: 4px;
          text-transform: none;

          .theme--light & {
            background-color: var(--tab-background-color-light);
          }

          .theme--dark & {
            background-color: var(--tab-background-color-dark);
          }
        }

        &.v-tab--active .v-btn {
          border:thin solid currentColor;
          color: var(--v-primary-base) !important;
          caret-color: var(--v-primary-base) !important;

          &.theme--light {
            background-color: var(--v-background-base, #FFFFFF);
          }

          &.theme--dark {
            background-color: var(--v-background-base, #1E1E1E);
          }
        }
      }
    }
  }
}
</style>
