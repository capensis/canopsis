<template>
  <v-layout class="gap-3">
    <v-flex :xs6="withChildInactiveDelay" xs12>
      <c-duration-field
        v-field="timebased.time_interval"
        :label="$t('metaAlarmRule.timeInterval')"
        required
      />
    </v-flex>
    <v-flex v-if="withChildInactiveDelay" xs6>
      <v-layout class="gap-3" align-center>
        <c-duration-field
          v-field="timebased.child_inactive_delay"
          :label="$t('metaAlarmRule.childInactiveDelay')"
          clearable
        />
        <div>
          <v-tooltip top>
            <template #activator="{ on }">
              <v-icon v-on="on">
                help
              </v-icon>
            </template>
            <span>{{ $t('metaAlarmRule.childInactiveDelayTooltip') }}</span>
          </v-tooltip>
        </div>
      </v-layout>
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'timebased',
    event: 'input',
  },
  props: {
    timebased: {
      type: Object,
      default: () => ({}),
    },
    withChildInactiveDelay: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
