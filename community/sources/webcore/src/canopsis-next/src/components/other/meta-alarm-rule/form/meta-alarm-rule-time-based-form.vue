<template>
  <v-layout column>
    <v-layout>
      <v-flex xs6>
        <v-layout class="gap-2">
          <c-duration-field
            v-field="timebased.time_interval"
            :label="$t('metaAlarmRule.timeInterval')"
            required
          />
          <c-help-icon
            :text="$t('metaAlarmRule.timeIntervalHelpText')"
            icon="help"
            max-width="300"
            top
          />
        </v-layout>
      </v-flex>
    </v-layout>
    <v-layout>
      <v-flex v-if="withChildInactiveDelay" xs6>
        <v-layout class="gap-2">
          <c-enabled-field
            :value="!!timebased.child_inactive_delay"
            :label="$t('metaAlarmRule.childInactiveDelay')"
            @input="enableChildInactiveDelay"
          />
          <c-help-icon
            :text="$t('metaAlarmRule.childInactiveDelayHelpText')"
            icon="help"
            max-width="300"
            top
          />
        </v-layout>
        <c-duration-field
          v-if="timebased.child_inactive_delay"
          v-field="timebased.child_inactive_delay"
          clearable
        />
      </v-flex>
    </v-layout>
  </v-layout>
</template>

<script>
import { durationToForm } from '@/helpers/date/duration';

import { useModelField } from '@/hooks/form/model-field';

export default {
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
  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    const enableChildInactiveDelay = (value) => {
      updateField('child_inactive_delay', value ? durationToForm() : undefined);
    };

    return {
      enableChildInactiveDelay,
    };
  },
};
</script>
