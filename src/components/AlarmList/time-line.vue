<template lang="pug">
  ul.timeline(v-if="fetchComplete")
    li.timeline-item(v-for="step in steps")
      time-item.time(:time="step.t")
      div(v-if="step._t !== 'statecounter'")
        state-flag.flag(:val="step.val" :isStatus="isStatus(step._t)")
        .header
          alarm-chips.chips(:val="step.val" :isStatus="isStatus(step._t)")
          p {{ step._t }}
        .content
          p {{ step.m }}
      div(v-else)
        state-flag.flag(isCroppedState)
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
            tr(v-for="(value, state) in step.val" v-if="state.startsWith('state:')")
              td State {{ getStateAlarmConvention(parseInt(state.replace('state:',''),10))('text') }} :
              td {{ value }}
</template>

<script>
import { normalize } from 'normalizr';
import { createNamespacedHelpers } from 'vuex';


import { API_ROUTES } from '@/config';
import { alarmSchema } from '@/store/schemas';
import request from '@/services/request';
import StateFlag from '../BasicComponent/state-flag.vue';
import AlarmChips from '../BasicComponent/alarm-chips.vue';
import TimeItem from '../BasicComponent/time-item.vue';

const { mapGetters } = createNamespacedHelpers('entities/alarmConvention');

export default {
  name: 'time-line',
  components: { TimeItem, AlarmChips, StateFlag },
  mounted() {
    this.fetchAlarm();
  },
  props: {
    idAlarm: {
      type: String,
    },
  },
  data() {
    return {
      dataAlarms: undefined,
      fetchComplete: false,
    };
  },
  computed: {
    ...mapGetters(['getStateAlarmConvention']),
    steps() {
      return Object.values(this.dataAlarms.entities.alarm)[0].v.steps;
    },
  },
  methods: {
    isStatus(stepTitle) {
      return stepTitle.startsWith('status');
    },
    async fetchAlarm() {
      try {
        const [data] = await request.get(API_ROUTES.alarmList, {
          params: {
            filter: { d: this.idAlarm },
            opened: 'true',
            resolved: 'true',
            sort_key: 't',
            sort_dir: 'DESC',
            limit: '1',
            with_steps: 'true',
          },
        });
        this.dataAlarms = normalize(data.alarms, [alarmSchema]);
        this.fetchComplete = true;
      } catch (err) {
        console.error(err);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  ul {
    list-style: none;
  }
  .timeline {
    margin: 0 auto;
    width: 90%
  }

  $border_line: #DDDDE0;

  .timeline-item {
    padding: 3em 2em 2em;
    position: relative;
    border-left: 2px solid $border_line;

    .time{
      position: absolute;
      color: #686868;
      left: 2em;
      top: 9px;
      display: block;
      font-size: 11px;
    }

    &:last-child  {
     border-image: linear-gradient(
     to bottom,
     $border-line 60%,
     white) 1 100%;
    }

  }

  .flag {
     top: 0;
     position: absolute;
     left: -19px;
     background: white;
   }

  .content {
    color: #858585;
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

    .chips{
      font-size: 15px;
    }

    p{
      font-size: 20px;
    }

  }

  p{
    overflow-wrap: break-word;
    text-overflow: ellipsis;
    width: 90%;
  }
</style>
