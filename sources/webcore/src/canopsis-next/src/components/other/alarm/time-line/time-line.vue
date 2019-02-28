<template lang="pug">
  div.timeline
    ul(v-for="(steps, day) in groupedSteps", :key="day")
      li(v-for="(step, index) in steps", :key="`step-${index}`")
        .timeline-item(v-show="index === 0")
          .date {{ day | date('short', true) }}
        .timeline-item
          .time {{ step.t | date('H:mm:SS', true) }}
          div(v-if="step._t !== 'statecounter'")
            flag.flag(:type="stepType(step._t)", :step="step")
            .card
              .header
                alarm-chips.chips(v-if="stepType(step._t) !== $constants.ENTITY_INFOS_TYPE.action",
                :value="step.val", :type="stepType(step._t)")
                p  &nbsp {{ step._t | stepTitle(step.a) }}
              .content
                p {{ step.m }}
          div(v-else)
            flag.flag(isCroppedState)
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
import { pickBy, orderBy, groupBy } from 'lodash';

import { ENTITIES_STATES_STYLES, ENTITY_INFOS_TYPE } from '@/constants';

import { stepTitle } from '@/helpers/timeline';

import Flag from '@/components/other/alarm/time-line/time-line-flag.vue';
import AlarmChips from '@/components/other/alarm/alarm-chips.vue';
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
        const orderedSteps = orderBy(alarm.v.steps, ['t'], 'desc');

        return groupBy(orderedSteps, step => moment.unix(step.t).startOf('day').format());
      }

      return {};
    },

    stepType() {
      return (type) => {
        if (type.startsWith('status')) {
          return ENTITY_INFOS_TYPE.status;
        } else if (type.startsWith('state')) {
          return ENTITY_INFOS_TYPE.state;
        }

        return ENTITY_INFOS_TYPE.action;
      };
    },

    stateName(state) {
      const stateValue = parseInt(state.replace('state:', ''), 10);
      return ENTITIES_STATES_STYLES[stateValue].text;
    },

    stateSteps(steps) {
      return pickBy(steps, (value, key) => key.startsWith('state:'));
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

  .content {
    padding-left: 20px;
    padding-top: 20px;
    overflow-wrap: break-word;
    width: 90%;
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
