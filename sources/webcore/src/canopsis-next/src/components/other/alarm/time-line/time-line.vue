<template lang="pug">
  div.timeline(:data-test="`alarmTimeLine-${alarm._id}`")
    ul(v-for="(steps, day) in groupedSteps", :key="day")
      li(v-for="(step, index) in steps", :key="`step-${index}`")
        .timeline-item(v-show="index === 0")
          .date {{ day | date('short', true) }}
        .timeline-item
          .time {{ step.t | date('H:mm:SS', true) }}
          time-line-flag.flag(:step="step")
          time-line-card(:step="step", :isHTMLEnabled="isHTMLEnabled")
</template>

<script>
import moment from 'moment';
import { orderBy, groupBy } from 'lodash';

import TimeLineFlag from '@/components/other/alarm/time-line/time-line-flag.vue';
import TimeLineCard from '@/components/other/alarm/time-line/time-line-card.vue';

import entitiesAlarmMixin from '@/mixins/entities/alarm';

export default {
  components: { TimeLineFlag, TimeLineCard },
  mixins: [entitiesAlarmMixin],
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    isHTMLEnabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      groupedSteps: {},
      steps: [],
    };
  },
  watch: {
    alarm(alarm, oldAlarm) {
      if (alarm._id !== oldAlarm._id) {
        this.groupedSteps = {};
        this.steps = [];
      }

      if (!alarm && alarm.v.steps) {
        this.fetchItem();
      } else {
        this.groupedSteps = this.groupSteps(alarm.v.steps);
        this.steps = alarm.v.steps;
      }
    },
  },
  mounted() {
    this.fetchItem();
  },
  methods: {
    fetchItem() {
      const params = {
        sort_key: 't',
        sort_dir: 'DESC',
        limit: 1,
        with_steps: true,
      };

      if (this.alarm.v.resolved) {
        params.resolved = true;
      }

      this.fetchAlarmItem({
        id: this.alarm._id,
        params,
      });
    },

    groupSteps(steps) {
      const orderedSteps = orderBy(steps, ['t'], 'desc');

      return groupBy(orderedSteps, step => moment.unix(step.t).startOf('day').format());
    },
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
    padding: 3em 2em 0;
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
