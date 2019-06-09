<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsSelect.title') }}
      div.font-italic.caption.ml-1(v-if="required") ({{ $t('settings.statsSelect.required') }})
      div.font-italic.caption.ml-1(v-else) ({{ $t('common.optionnal') }})
    v-container
      v-alert(:value="errors.has('stats')", type="error") {{ $t('settings.statsSelect.required') }}
      v-btn(@click="showAddStatModal") {{ $t('modals.addStat.title.add') }}
      v-list.secondary(dark)
        draggable(
        :value="orderedStats",
        :options="draggableOptions",
        @input="updateStatsPositions"
        )
          v-list-group(v-for="stat in orderedStats", :key="stat.title")
            v-list-tile(slot="activator")
              v-list-tile-content
                v-list-tile-title {{ stat.title }}
              v-list-tile-action
                v-layout
                  v-btn.primary.mx-1(@click.stop="showEditStatModal(stat.title, stat)", fab, small, depressed)
                    v-icon edit
                  v-btn.error(@click.stop="showDeleteStatModal(stat.title)", fab, small, depressed)
                    v-icon delete
            v-list-tile
              v-list-tile-title {{ $t('common.stat') }}: {{ stat.stat }}
            v-list-tile
              v-list-tile-title {{ $t('common.trend') }}: {{ stat.trend }}
            v-list-tile
              v-list-tile-title {{ $t('common.parameters') }}: {{ stat.parameters }}
</template>

<script>
import { omit } from 'lodash';
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { MODALS } from '@/constants';

import { setInSeveral } from '@/helpers/immutable';

import modalMixin from '@/mixins/modal';
import formMixin from '@/mixins/form';

export default {
  inject: ['$validator'],
  components: { Draggable },
  mixins: [modalMixin, formMixin],
  model: {
    prop: 'stats',
    event: 'input',
  },
  props: {
    stats: {
      type: Object,
      default: () => ({}),
    },
    required: {
      type: Boolean,
      default: false,
    },
    withTrend: {
      type: Boolean,
      default: false,
    },
    withSorting: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        disabled: !this.withSorting,
      };
    },

    orderedStats() {
      return Object.entries(this.stats)
        .map(([key, value]) => ({ ...value, title: key }))
        .sort((a, b) => a.position - b.position);
    },
  },
  watch: {
    stats(value) {
      if (this.required) {
        this.$validator.validate('stats', value);
      }
    },
  },
  created() {
    if (this.required) {
      this.$validator.attach({
        name: 'stats',
        rules: 'required',
        getter: () => Object.values(this.stats),
        context: () => this,
      });
    }
  },
  methods: {
    showAddStatModal() {
      this.showModal({
        name: MODALS.addStat,
        config: {
          title: this.$t('modals.addStat.title.add'),
          withTrend: this.withTrend,
          action: (stat) => {
            const newStat = {
              ...omit(stat, ['title', 'parameters']),

              parameters: stat.stat.options.reduce((acc, option) => {
                acc[option] = stat.parameters[option];

                return acc;
              }, {}),
            };

            this.updateField(stat.title, newStat);
          },
        },
      });
    },

    showEditStatModal(statTitle) {
      this.showModal({
        name: MODALS.addStat,
        config: {
          title: this.$t('modals.addStat.title.edit'),
          withTrend: this.withTrend,
          stat: this.stats[statTitle],
          statTitle,
          action: newStat => this.updateAndMoveField(statTitle, newStat.title, omit(newStat, ['title'])),
        },
      });
    },

    showDeleteStatModal(statTitle) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.removeField(statTitle),
        },
      });
    },

    updateStatsPositions(newOrderedStats) {
      const modifiers = newOrderedStats.reduce((acc, orderedStat, index) => {
        acc[orderedStat.title] = stat => ({ ...stat, position: index });

        return acc;
      }, {});

      const newStats = setInSeveral(this.stats, modifiers);

      this.updateModel(newStats);
    },
  },
};
</script>
