<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsSelect') }}
    v-container
      v-card
        v-card-title
          v-select(
          v-model="form.stat",
          hide-details,
          :items="statsTypes",
          )
        v-divider
        v-layout.green.darken-4.white--text(align-center, justify-center)
          div.text-xs-center.my-2 Options
        v-container.pt-0(fluid)
          v-text-field(:placeholder="$t('common.title')", v-model="form.title", hide-details)
          v-switch(:label="$t('common.recursive')", v-model="form.parameters.recursive", hide-details)
          v-switch(:label="$t('common.trend')", v-model="form.trend", hide-details)
          v-select(
          :placeholder="$t('common.states')",
          :items="stateTypes",
          v-model="form.parameters.states",
          multiple,
          chips,
          hide-details
          )
          v-combobox(
          :placeholder="$t('common.authors')",
          v-model="form.parameters.authors",
          hide-details,
          chips,
          multiple
          )
          v-text-field(:placeholder="$t('common.sla')", v-model="form.sla", hide-details)
        v-btn.ma-0(@click="addStat") Add stat
</template>

<script>
import omit from 'lodash/omit';
import { STATS_TYPES, ENTITIES_STATES } from '@/constants';

export default {
  data() {
    return {
      stats: {},
      form: {
        stat: '',
        title: '',
        sla: '',
        trend: true,
        parameters: {
          states: [],
          recursive: false,
          authors: [],
        },
      },
      error: '',
    };
  },
  computed: {
    statsTypes() {
      return Object.values(STATS_TYPES).map(item => ({ value: item, text: this.$t(`stats.types.${item}`) }));
    },
    stateTypes() {
      return Object.keys(ENTITIES_STATES).map(item => ({ value: ENTITIES_STATES[item], text: item }));
    },
  },
  methods: {
    addStat() {
      if (this.stats[this.form.title]) {
        this.error = 'Stat with this title already exists';
      } else {
        this.error = '';
        this.stats[this.form.title] = omit(this.form, ['title']);
      }
    },
  },
};
</script>

