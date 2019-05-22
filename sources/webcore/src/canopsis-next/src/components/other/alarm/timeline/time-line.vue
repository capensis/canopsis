<template lang="pug">
  ul.timeline
    li(v-for="step in steps")
      .timeline-item(v-if="isNewDate(step.t)")
        .date {{ getFormattedDate(step.t) }}
      .timeline-item
        .time {{ getFormattedTime(step.t) }}
        div(v-if="step._t !== 'statecounter'")
          flag.flag(:type="stepType(step._t)", :step="step")
          .card
            .header
              alarm-chips.chips(v-if="stepType(step._t) !== $constants.ENTITY_INFOS_TYPE.action",
                                :value="step.val", :type="stepType(step._t)")
              p  &nbsp {{ step._t | stepTitle(step.a, step.role) }}
            .content
              p(v-if="isHTMLEnabled", v-html="step.m")
              p(v-else) {{ step.m }}
        div(v-else)
          alarm-flag.flag(isCroppedState)
          .header
            p Cropped State (since last change of status)
          .content
            table
              tr
                td State increased :
                td {{ step.val.stateinc }}
              tr
                td State decreases :
                td {{ step.val.statedec }}
              tr(v-for="(value, state) in stateSteps(step.val)")
                td State {{ stateName(state) }} :
                td {{ value }}
</template>

<script>
import moment from 'moment';
import { pickBy } from 'lodash';

import { stepTitle, stepType } from '@/helpers/timeline';

import Flag from '@/components/other/alarm/timeline/timeline-flag.vue';
import AlarmChips from '@/components/other/alarm/alarm-chips.vue';
import { numericSortObject } from '@/helpers/sorting';
import entitiesAlarmMixin from '@/mixins/entities/alarm';

/**
   * Component for the timeline of an alarm, on the alarmslist
   *
   * @module alarm
   *
   * @prop {alarmProp} [alarmProps] - Properties of an alarm
   */
export default {
  components: { AlarmChips, Flag },
  filters: {
    stepTitle,
  },
  mixins: [
    entitiesAlarmMixin,
  ],
  props: {
    alarmProps: {
      type: Object,
      required: true,
    },
    isHTMLEnabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    steps() {
      const alarm = this.getAlarmItem(this.alarmProps._id);
      if (alarm && alarm.v.steps) {
        const steps = [...alarm.v.steps];
        return numericSortObject(steps, 't', 'desc');
      }
      return [];
    },
    stateName(state) {
      const stateValue = parseInt(state.replace('state:', ''), 10);
      return this.$constants.ENTITIES_STATES_STYLES[stateValue].text;
    },
    stateSteps(steps) {
      return pickBy(steps, (value, key) => key.startsWith('state:'));
    },
  },
  mounted() {
    this.fetchAlarmItem({
      id: this.alarmProps.d,
      params: {
        opened: 'true',
        resolved: 'true',
        sort_key: 't',
        sort_dir: 'DESC',
        limit: '1',
        with_steps: 'true',
      },
    });
    this.lastDate = undefined;
  },
  updated() {
    // Useful like for example when the user change the translation
    this.lastDate = undefined;
  },
  methods: {
    stepType,
    isNewDate(timestamp) {
      const date = moment.unix(timestamp);
      if (!this.lastDate ||
            (date.diff(this.lastDate, 'days') < 0)) {
        this.lastDate = date;
        return true;
      }
      return false;
    },
    getFormattedDate(timestamp) {
      return moment.unix(timestamp).format('DD/MM/YYYY');
    },
    getFormattedTime(timestamp) {
      return moment.unix(timestamp).format('HH:mm:SS');
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
  }

  li:last-child {
    .timeline-item {
      border-image: linear-gradient(
          to bottom,
          $border-line 60%,
          $background) 1 100%;
    }
  }

  .timeline {
    margin: 0 auto;
    width: 90%;
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

  .content {
    padding-left: 20px;
    padding-top: 20px;
    overflow-wrap: break-word;
    word-break: break-all;
    width: 90%;
    max-height: 600px;
    overflow-y: auto;
  }

  .header {
    color: #686868;
    display: flex;
    align-items: baseline;
    font-weight: bold;
    border-bottom: solid 1px $border_line;
    padding-left: 5px;
    padding-top: 5px;


    .chips {
      font-size: 15px;
      height: 25px;
    }

    p {
      font-size: 15px;
    }

  }

  p {
    overflow-wrap: break-word;
    text-overflow: ellipsis;
    width: 90%;
  }
</style>
