<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsSelect.title') }}
      .font-italic.caption.ml-1 ({{ $t('settings.statsSelect.required') }})
    v-container
      v-alert(:value="errors.has('stats')", type="error") {{ $t('settings.statsSelect.required') }}
      v-btn(@click="openAddStatModal") {{ $t('modals.addStat.title.add') }}
      v-list.secondary(dark)
        v-list-group(v-for="(stat, key) in stats", :key="key")
          v-list-tile(slot="activator")
            v-list-tile-content
              v-list-tile-title {{ key }}
            v-list-tile-action
              v-layout
                v-btn.primary.mx-1(@click.stop="editStat(key, stat)", fab, small, depressed)
                  v-icon edit
                v-btn.error(@click.stop="deleteStat(key)", fab, small, depressed)
                  v-icon delete
          v-list-tile
            v-list-tile-title {{ $t('common.stat') }}: {{ stat.stat }}
          v-list-tile
            v-list-tile-title {{ $t('common.trend') }}: {{ stat.trend }}
          v-list-tile
            v-list-tile-title {{ $t('common.parameters') }}: {{ stat.parameters }}
</template>

<script>
import omit from 'lodash/omit';
import set from 'lodash/set';
import unset from 'lodash/unset';

import modalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';

export default {
  inject: ['$validator'],
  mixins: [modalMixin],
  model: {
    prop: 'stats',
    event: 'input',
  },
  props: {
    stats: {
      type: Object,
      default: () => ({}),
    },
  },
  watch: {
    stats(value) {
      this.$validator.validate('stats', value);
    },
  },
  created() {
    this.$validator.attach('stats', 'required', {
      getter: () => Object.values(this.stats),
      context: () => this,
    });
  },
  methods: {
    openAddStatModal() {
      this.showModal({
        name: MODALS.addStat,
        config: {
          title: 'modals.addStat.title.add',
          action: (stat) => {
            const newValue = { ...this.stats };
            const newStat = omit(stat, ['title', 'parameters', 'stat']);
            newStat.stat = stat.stat.value;
            newStat.parameters = {};
            stat.stat.options.forEach((option) => {
              newStat.parameters[option] = stat.parameters[option];
            });
            this.$emit('input', set(newValue, stat.title, newStat));
          },
        },
      });
    },

    deleteStat(stat) {
      const newValue = { ...this.stats };
      unset(newValue, stat);
      this.$emit('input', newValue);
    },

    editStat(statTitle, stat) {
      this.showModal({
        name: MODALS.addStat,
        config: {
          title: 'modals.addStat.title.edit',
          stat,
          statTitle,
          action: (newStat) => {
            // Delete the stat that we want to edit
            const newValue = { ...this.stats };
            unset(newValue, statTitle);
            // Set the edited stat in newValue object, and send it to parent with input event
            this.$emit('input', set(newValue, newStat.title, omit(newStat, ['title'])));
          },
        },
      });
    },
  },
};
</script>
