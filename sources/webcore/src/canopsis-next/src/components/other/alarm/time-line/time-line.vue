<template lang="pug">
  div.timeline
    ul(v-for="(steps, day) in groupedSteps", :key="day")
      li(v-for="(step, index) in steps", :key="`step-${index}`")
        .timeline-item(v-show="index === 0")
          .date {{ day | date('short', true) }}
        .timeline-item
          .time {{ step.t | date('H:mm:SS', true) }}
          time-line-flag.flag(:step="step")
          time-line-card(:step="step")
</template>

<script>
import moment from 'moment';
import { orderBy, groupBy } from 'lodash';

import TimeLineFlag from '@/components/other/alarm/time-line/time-line-flag.vue';
import TimeLineCard from '@/components/other/alarm/time-line/time-line-card.vue';

import entitiesAlarmMixin from '@/mixins/entities/alarm';

/**
   * Component for the timeline of an alarm, on the alarmslist
   *
   * @module alarm
   *
   * @prop {alarmProp} [alarmProps] - Properties of an alarm
   */
export default {
  components: { TimeLineFlag, TimeLineCard },
  mixins: [entitiesAlarmMixin],
  props: {
    alarmProps: {
      type: Object,
      required: true,
    },
  },
  computed: {
    groupedSteps() {
      const alarm = this.getAlarmItem(this.alarmProps._id);

      if (alarm && alarm.v.steps) {
        const orderedSteps = orderBy(alarm.v.steps.concat({
          a: 'system',
          _t: 'statecounter',
          m: '',
          t: 1541411224,
          val: {
            'state:2': 1,
            'state:0': 1,
            statedec: 1,
            stateinc: 1,
          },
        }), ['t'], 'desc');

        return groupBy(orderedSteps, step => moment.unix(step.t).startOf('day').format());
      }

      return {};
    },
  },
  mounted() {
    this.fetchAlarmItem({
      id: this.alarmProps.d,
      params: {
        sort_key: 't',
        sort_dir: 'DESC',
        limit: 1,
        with_steps: true,
      },
    });
  },
};
</script>

<style lang="scss" scoped>
  $border_line: #DDDDE0;
  $background: white;

  ul {
    list-style: none;
    color: #858585;

    &:last-child {
      li:last-child {
        .timeline-item:last-child {
          border-image: linear-gradient(
              to bottom,
              $border-line 60%,
              $background) 1 100%;
        }
      }
    }
  }

  .timeline {
    margin: 0 auto;
    width: 90%
  }

  .timeline-item {
    padding: 3em 2em 0em;
    position: relative;
    border-left: 2px solid $border_line;

    .time {
      position: absolute;
      left: 2em;
      top: 9px;
      display: block;
      font-size: 11px;
    }
  }

  .flag, .date {
    top: 0;
    position: absolute;
    background: $background;
  }

  .flag {
    top: 4px;
    left: -13px;
  }

  .date {
    top: 4px;
    left: -11px;
  }
</style>
