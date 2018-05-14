<template lang="pug">
  ul.timeline(v-if="fetchComplete")
    li.timeline-item(v-for= "step in steps")
      time-item.time(:time="step.t")
      state-flag.flag(:val="step.val")
      .header
        state-chips.chips(:val="step.val")
        p {{ step._t }}
      .content
        p {{ step.m }}
</template>

<script>
import { normalize } from 'normalizr';

import { API_ROUTES } from '@/config';
import { alarmSchema } from '@/store/schemas';
import request from '@/services/request';
import StateFlag from '../BasicComponent/state-flag.vue';
import StateChips from '../BasicComponent/state-chips.vue';
import TimeItem from '../BasicComponent/time-item.vue';

export default {
  name: 'time-line',
  components: { TimeItem, StateChips, StateFlag },
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
    steps() {
      return Object.values(this.dataAlarms.entities.alarm)[0].v.steps;
    },
  },
  methods: {
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
     left: -15px;
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
