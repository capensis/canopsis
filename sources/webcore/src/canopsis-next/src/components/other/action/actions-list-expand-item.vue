<template lang="pug">
  v-tabs(color="secondary lighten-1", dark, centered)
    v-tab {{ $t('actions.table.expand.tabs.general') }}
    v-tab-item
      v-layout.py-3.secondary.lighten-2(row)
        v-flex(xs12, sm10, offset-sm1)
          v-card
            v-card-text
              v-list
                v-layout(wrap, justify-center, align-center)
                  v-flex(xs6)
                    v-list-tile
                      v-list-tile-avatar
                        v-icon person
                      v-list-tile-content
                        v-layout(align-center)
                          span.font-weight-bold Author
                          span : {{ action.parameters.author }}
                  template(v-if="action.type === $constants.ACTION_TYPES.pbehavior")
                    v-flex(xs6)
                      v-list-tile
                        v-list-tile-avatar
                          v-icon short_text
                        v-list-tile-content
                          v-layout(align-center)
                            span.font-weight-bold Name
                            span : {{ action.parameters.name }}
                    v-flex(xs6)
                      v-list-tile
                        v-list-tile-avatar
                          v-icon assignment
                        v-list-tile-content
                          v-layout(align-center)
                            span.font-weight-bold Type
                            span : {{ action.parameters.type_ }}
                    v-flex(xs6)
                      v-list-tile
                        v-list-tile-avatar
                          v-icon assignment
                        v-list-tile-content
                          v-layout(align-center)
                            span.font-weight-bold Reason
                            span : {{ action.parameters.reason }}
                    v-flex(xs6)
                      v-list-tile
                        v-list-tile-avatar
                          v-icon alarm_on
                        v-list-tile-content
                          v-layout(align-center)
                            span.font-weight-bold Start
                            span : {{ action.parameters.tstart | date('long') }}
                    v-flex(xs6)
                      v-list-tile
                        v-list-tile-avatar
                          v-icon alarm_off
                        v-list-tile-content
                          v-layout(align-center)
                            span.font-weight-bold End
                            span : {{ action.parameters.tstop | date('long') }}
                  template(v-if="action.type === $constants.ACTION_TYPES.snooze")
                    v-flex(xs6)
                      v-list-tile
                        v-list-tile-avatar
                          v-icon short_text
                        v-list-tile-content
                          v-layout(align-center)
                            span.font-weight-bold Duration
                            span : {{ action.parameters.duration }}
                    v-flex(xs12)
                      v-list-tile
                        v-list-tile-avatar
                          v-icon short_text
                        v-list-tile-content
                          v-layout(align-center)
                            span.font-weight-bold Message
                            span : {{ action.parameters.message }}
    v-tab {{ $t('actions.table.expand.tabs.hook') }}
    v-tab-item
      v-layout.py-3.secondary.lighten-2(row)
        v-flex(xs12, sm8, offset-sm2)
          v-card
            v-card-text
              div
                v-layout(row, wrap)
                  v-flex(xs12)
                    v-select(
                    :value="action.hook.triggers",
                    :items="availableTriggers",
                    disabled,
                    :label="$t('webhook.tabs.hook.fields.triggers')",
                    multiple,
                    chips
                    )
                  v-flex(xs12)
                    patterns-list(
                    :patterns="action.hook.event_patterns",
                    disabled,
                    :operators="$constants.WEBHOOK_EVENT_FILTER_RULE_OPERATORS",
                    )
</template>

<script>
import { WEBHOOK_TRIGGERS } from '@/constants';

import PatternsList from '@/components/other/shared/patterns-list/patterns-list.vue';

export default {
  components: {
    PatternsList,
  },
  props: {
    action: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      availableTriggers: Object.values(WEBHOOK_TRIGGERS),
    };
  },
};
</script>
