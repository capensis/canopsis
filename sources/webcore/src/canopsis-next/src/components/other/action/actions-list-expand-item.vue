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
                          span.font-weight-bold {{ $t('actions.table.expand.tabs.author') }}
                          span : {{ action.parameters.author }}
                  template(v-if="action.type === $constants.ACTION_TYPES.pbehavior")
                    v-flex(xs6)
                      v-list-tile
                        v-list-tile-avatar
                          v-icon short_text
                        v-list-tile-content
                          v-layout(align-center)
                            span.font-weight-bold {{ $t('actions.table.expand.tabs.pbehavior.name') }}
                            span : {{ action.parameters.name }}
                    v-flex(xs6)
                      v-list-tile
                        v-list-tile-avatar
                          v-icon assignment
                        v-list-tile-content
                          v-layout(align-center)
                            span.font-weight-bold {{ $t('actions.table.expand.tabs.pbehavior.type') }}
                            span : {{ action.parameters.type_ }}
                    v-flex(xs6)
                      v-list-tile
                        v-list-tile-avatar
                          v-icon assignment
                        v-list-tile-content
                          v-layout(align-center)
                            span.font-weight-bold {{ $t('actions.table.expand.tabs.pbehavior.reason') }}
                            span : {{ action.parameters.reason }}
                    v-flex(xs6)
                      v-list-tile
                        v-list-tile-avatar
                          v-icon alarm_on
                        v-list-tile-content
                          v-layout(align-center)
                            span.font-weight-bold {{ $t('actions.table.expand.tabs.pbehavior.start') }}
                            span : {{ action.parameters.tstart | date('long') }}
                    v-flex(xs6)
                      v-list-tile
                        v-list-tile-avatar
                          v-icon alarm_off
                        v-list-tile-content
                          v-layout(align-center)
                            span.font-weight-bold {{ $t('actions.table.expand.tabs.pbehavior.end') }}
                            span : {{ action.parameters.tstop | date('long') }}
                  template(v-if="action.type === $constants.ACTION_TYPES.snooze")
                    v-flex(xs6)
                      v-list-tile
                        v-list-tile-avatar
                          v-icon short_text
                        v-list-tile-content
                          v-layout(align-center)
                            span.font-weight-bold {{ $t('actions.table.expand.tabs.snooze.duration') }}
                            span : {{ action.parameters.duration }}
                    v-flex(xs12)
                      v-list-tile
                        v-list-tile-avatar
                          v-icon short_text
                        v-list-tile-content
                          v-layout(align-center)
                            span.font-weight-bold {{ $t('actions.table.expand.tabs.snooze.message') }}
                            span : {{ action.parameters.message }}
    template(v-if="action.type === $constants.ACTION_TYPES.pbehavior")
      v-tab {{ $t('pbehaviors.tabs.comments') }}
      v-tab-item
        v-layout.py-3.secondary.lighten-2(row)
          v-flex(xs12, sm6, offset-sm3)
            v-card
              v-card-text
                v-list(two-line)
                  v-list-tile(v-if="!action.parameters.comments || !action.parameters.comments.length")
                    v-list-tile-content
                      v-list-tile-title {{ $t('tables.noData') }}
                  template(v-for="(comment, index) in action.parameters.comments")
                    v-list-tile(:key="comment._id")
                      v-list-tile-content
                        v-list-tile-title {{ comment.author }}
                        v-list-tile-sub-title {{ comment.message }}
                    v-divider(v-if="index < action.parameters.comments.length - 1", :key="`divider-${index}`")
      v-tab(v-if="rRule") {{ $t('pbehaviors.rrule') }}
      v-tab-item
        v-layout.py-3.secondary.lighten-2(row)
          v-flex(xs12, sm6, offset-sm3)
            v-card
              v-card-text()
                v-layout(row)
                  v-flex(xs2)
                    strong {{ $t('rRule.stringLabel') }}
                  v-flex(xs10)
                    p.rrule-paragraph {{ rRuleString }}
                v-layout(row)
                  v-flex(xs2)
                    strong {{ $t('rRule.textLabel') }}
                  v-flex(xs10)
                    p.rrule-paragraph {{ rRuleText }}

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
import { rrulestr } from 'rrule';

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
  computed: {
    rRule() {
      return this.action.parameters.rrule ? rrulestr(this.action.parameters.rrule) : null;
    },

    rRuleString() {
      return this.rRule ? this.rRule.toString() : '';
    },

    rRuleText() {
      return this.rRule ? this.rRule.toText() : '';
    },
  },
};
</script>
