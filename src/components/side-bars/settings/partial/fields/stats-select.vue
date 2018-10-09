<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsSelect') }}
    v-container
      v-btn(@click="openAddStatModal") {{ $t('modals.addStat.title.add') }}
      v-list(dark)
        v-list-group.my-1.grey(v-for="(stat, key) in value", :key="key")
          v-list-tile(slot="activator") {{ key }}
            v-layout(justify-end)
              v-btn.green.darken-4.white--text(@click.stop="editStat(key, stat)", fab, small, depressed)
                v-icon edit
              v-btn.red.darken-4.white--text(@click.stop="deleteStat(key)", fab, small, depressed)
                v-icon delete
          v-container(fluid)
            p {{ $t('common.stat') }}: {{ stat.stat.value }}
            p {{ $t('common.trend') }}: {{ stat.trend }}
            p {{ $t('common.parameters') }}: {{ stat.parameters }}
</template>

<script>
import omit from 'lodash/omit';
import set from 'lodash/set';
import unset from 'lodash/unset';
import modalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';

export default {
  mixins: [modalMixin],
  props: {
    value: {
      type: Object,
    },
  },
  data() {
    return {
      editing: false,
      editingStatTitle: '',
    };
  },
  methods: {
    openAddStatModal() {
      this.showModal({
        name: MODALS.addStat,
        config: {
          title: 'modals.addStat.title.add',
          action: (stat) => {
            const newValue = { ...this.value };
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
      const newValue = { ...this.value };
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
            const newValue = { ...this.value };
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
